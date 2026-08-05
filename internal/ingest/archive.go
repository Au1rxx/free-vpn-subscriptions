package ingest

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/store"
)

// ArchiveReference identifies one gzip-compressed body by the SHA-256 of its
// uncompressed bytes.
type ArchiveReference struct {
	Key             string
	SHA256          [32]byte
	OriginalBytes   int64
	CompressedBytes int64
}

type ArchiveStatus struct {
	Files, Referenced, Missing, Corrupt, Unreferenced int64
	Bytes                                             int64
	QuarantineFiles, QuarantineBytes                  int64
}

// PayloadArchive stores large response bodies below one owner-only root.
type PayloadArchive struct {
	root string
}

type ArchiveError struct {
	Code string
	Err  error
}

func (e *ArchiveError) Error() string { return e.Code + ": " + e.Err.Error() }
func (e *ArchiveError) Unwrap() error { return e.Err }

func archiveErrorCode(err error) string {
	if typed, ok := err.(*ArchiveError); ok {
		return typed.Code
	}
	return ""
}

func NewPayloadArchive(root string) (*PayloadArchive, error) {
	if strings.TrimSpace(root) == "" {
		return nil, fmt.Errorf("archive root is required")
	}
	absolute, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(absolute, 0o700); err != nil {
		return nil, err
	}
	if err := os.Chmod(absolute, 0o700); err != nil {
		return nil, err
	}
	resolved, err := filepath.EvalSymlinks(absolute)
	if err != nil {
		return nil, err
	}
	return &PayloadArchive{root: resolved}, nil
}

func (a *PayloadArchive) Root() string { return a.root }

func InspectArchive(ctx context.Context, archive *PayloadArchive, references []ArchiveReference, quarantineDirectory string) (ArchiveStatus, error) {
	if archive == nil {
		return ArchiveStatus{}, fmt.Errorf("payload archive is required")
	}
	files := make(map[string]int64)
	status := ArchiveStatus{}
	err := filepath.WalkDir(archive.root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".gz") {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(archive.root, path)
		if err != nil {
			return err
		}
		key := filepath.ToSlash(relative)
		files[key] = info.Size()
		status.Files++
		status.Bytes += info.Size()
		return nil
	})
	if err != nil {
		return ArchiveStatus{}, err
	}
	expected := make(map[string]bool, len(references))
	for _, reference := range references {
		expected[reference.Key] = true
		size, ok := files[reference.Key]
		if !ok {
			status.Missing++
			continue
		}
		path, err := archive.resolve(reference)
		if err != nil || archive.rejectSymlinkComponents(path) != nil || (reference.CompressedBytes > 0 && size != reference.CompressedBytes) {
			status.Corrupt++
			continue
		}
		status.Referenced++
	}
	for key := range files {
		if !expected[key] {
			status.Unreferenced++
		}
	}
	if strings.TrimSpace(quarantineDirectory) == "" {
		return status, nil
	}
	err = filepath.WalkDir(quarantineDirectory, func(_ string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			if os.IsNotExist(walkErr) {
				return nil
			}
			return walkErr
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		status.QuarantineFiles++
		status.QuarantineBytes += info.Size()
		return nil
	})
	if os.IsNotExist(err) {
		err = nil
	}
	return status, err
}

func (a *PayloadArchive) Put(ctx context.Context, body []byte) (ArchiveReference, error) {
	if err := ctx.Err(); err != nil {
		return ArchiveReference{}, err
	}
	digest := sha256.Sum256(body)
	reference := ArchiveReference{
		Key:           archiveKey(digest),
		SHA256:        digest,
		OriginalBytes: int64(len(body)),
	}
	destination, err := a.resolve(reference)
	if err != nil {
		return ArchiveReference{}, err
	}
	if err := a.makeSafeDirectory(filepath.Dir(destination)); err != nil {
		return ArchiveReference{}, err
	}
	if info, err := os.Stat(destination); err == nil {
		reference.CompressedBytes = info.Size()
		if _, err := a.Read(ctx, reference); err != nil {
			return ArchiveReference{}, err
		}
		return reference, nil
	} else if !os.IsNotExist(err) {
		return ArchiveReference{}, err
	}

	staging := filepath.Join(a.root, ".staging")
	if err := a.makeSafeDirectory(staging); err != nil {
		return ArchiveReference{}, err
	}
	temporary, err := os.CreateTemp(staging, ".payload-*.tmp")
	if err != nil {
		return ArchiveReference{}, err
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)
	if err := temporary.Chmod(0o600); err != nil {
		temporary.Close()
		return ArchiveReference{}, err
	}
	writer := gzip.NewWriter(temporary)
	if _, err := io.Copy(writer, &contextReader{ctx: ctx, reader: bytes.NewReader(body)}); err != nil {
		writer.Close()
		temporary.Close()
		return ArchiveReference{}, err
	}
	if err := writer.Close(); err != nil {
		temporary.Close()
		return ArchiveReference{}, err
	}
	if err := temporary.Sync(); err != nil {
		temporary.Close()
		return ArchiveReference{}, err
	}
	info, err := temporary.Stat()
	if err != nil {
		temporary.Close()
		return ArchiveReference{}, err
	}
	if err := temporary.Close(); err != nil {
		return ArchiveReference{}, err
	}
	if err := ctx.Err(); err != nil {
		return ArchiveReference{}, err
	}
	if err := os.Rename(temporaryPath, destination); err != nil {
		return ArchiveReference{}, err
	}
	if err := syncDirectory(filepath.Dir(destination)); err != nil {
		return ArchiveReference{}, err
	}
	reference.CompressedBytes = info.Size()
	return reference, nil
}

func (a *PayloadArchive) Read(ctx context.Context, reference ArchiveReference) ([]byte, error) {
	path, err := a.resolve(reference)
	if err != nil {
		return nil, err
	}
	if err := a.rejectSymlinkComponents(path); err != nil {
		if os.IsNotExist(err) {
			return nil, &ArchiveError{Code: "archive_missing", Err: err}
		}
		return nil, err
	}
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, &ArchiveError{Code: "archive_missing", Err: err}
		}
		return nil, err
	}
	defer file.Close()
	if reference.CompressedBytes > 0 {
		info, statErr := file.Stat()
		if statErr != nil {
			return nil, statErr
		}
		if info.Size() != reference.CompressedBytes {
			return nil, &ArchiveError{Code: "archive_integrity", Err: fmt.Errorf("compressed bytes=%d want=%d", info.Size(), reference.CompressedBytes)}
		}
	}
	reader, err := gzip.NewReader(file)
	if err != nil {
		return nil, &ArchiveError{Code: "archive_integrity", Err: err}
	}
	defer reader.Close()
	if reference.OriginalBytes < 0 || reference.OriginalBytes == int64(^uint64(0)>>1) {
		return nil, &ArchiveError{Code: "archive_integrity", Err: fmt.Errorf("invalid original bytes %d", reference.OriginalBytes)}
	}
	body, err := io.ReadAll(io.LimitReader(&contextReader{ctx: ctx, reader: reader}, reference.OriginalBytes+1))
	if err != nil {
		return nil, &ArchiveError{Code: "archive_integrity", Err: err}
	}
	if int64(len(body)) != reference.OriginalBytes {
		return nil, &ArchiveError{Code: "archive_integrity", Err: fmt.Errorf("original bytes=%d want=%d", len(body), reference.OriginalBytes)}
	}
	if digest := sha256.Sum256(body); digest != reference.SHA256 {
		return nil, &ArchiveError{Code: "archive_integrity", Err: fmt.Errorf("sha256 mismatch")}
	}
	return body, nil
}

func (a *PayloadArchive) resolve(reference ArchiveReference) (string, error) {
	want := archiveKey(reference.SHA256)
	if reference.Key != want || strings.Contains(reference.Key, "\\") {
		return "", &ArchiveError{Code: "archive_key_invalid", Err: fmt.Errorf("key %q does not match sha256", reference.Key)}
	}
	path := filepath.Join(a.root, filepath.FromSlash(reference.Key))
	relative, err := filepath.Rel(a.root, path)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", &ArchiveError{Code: "archive_key_invalid", Err: fmt.Errorf("key %q escapes archive root", reference.Key)}
	}
	return path, nil
}

func (a *PayloadArchive) makeSafeDirectory(path string) error {
	relative, err := filepath.Rel(a.root, path)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return &ArchiveError{Code: "archive_path_unsafe", Err: fmt.Errorf("directory escapes archive root")}
	}
	current := a.root
	for _, part := range strings.Split(relative, string(filepath.Separator)) {
		if part == "" || part == "." {
			continue
		}
		current = filepath.Join(current, part)
		info, statErr := os.Lstat(current)
		switch {
		case statErr == nil && info.Mode()&os.ModeSymlink != 0:
			return &ArchiveError{Code: "archive_path_unsafe", Err: fmt.Errorf("symlink component %q", current)}
		case statErr == nil && !info.IsDir():
			return &ArchiveError{Code: "archive_path_unsafe", Err: fmt.Errorf("non-directory component %q", current)}
		case os.IsNotExist(statErr):
			if err := os.Mkdir(current, 0o700); err != nil {
				return err
			}
		case statErr != nil:
			return statErr
		}
	}
	return nil
}

func (a *PayloadArchive) rejectSymlinkComponents(path string) error {
	relative, err := filepath.Rel(a.root, path)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return &ArchiveError{Code: "archive_path_unsafe", Err: fmt.Errorf("path escapes archive root")}
	}
	current := a.root
	for _, part := range strings.Split(relative, string(filepath.Separator)) {
		if part == "" || part == "." {
			continue
		}
		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return &ArchiveError{Code: "archive_path_unsafe", Err: fmt.Errorf("symlink component %q", current)}
		}
	}
	return nil
}

func archiveKey(digest [32]byte) string {
	encoded := hex.EncodeToString(digest[:])
	return "sha256/" + encoded[:2] + "/" + encoded[2:4] + "/" + encoded + ".gz"
}

type contextReader struct {
	ctx    context.Context
	reader io.Reader
}

func shouldArchiveBody(bodyBytes, thresholdBytes int64) bool {
	return thresholdBytes > 0 && bodyBytes > thresholdBytes
}

func archiveFetchWrite(ctx context.Context, archive *PayloadArchive, write store.FetchWrite) (store.FetchWrite, error) {
	if archive == nil {
		return store.FetchWrite{}, fmt.Errorf("payload archive is required")
	}
	if len(write.Body) == 0 {
		return store.FetchWrite{}, fmt.Errorf("payload archive body is required")
	}
	reference, err := archive.Put(ctx, write.Body)
	if err != nil {
		return store.FetchWrite{}, err
	}
	write.ResponseBytes = int64(len(write.Body))
	write.Body = nil
	write.ExternalPayload = &store.ExternalPayload{
		SHA256: reference.SHA256, OriginalBytes: reference.OriginalBytes,
		CompressedBytes: reference.CompressedBytes, Compression: "gzip", ArchiveKey: reference.Key,
	}
	return write, nil
}

func (r *contextReader) Read(buffer []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.reader.Read(buffer)
}
