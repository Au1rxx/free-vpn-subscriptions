package ingest

import (
	"bytes"
	"context"
	"crypto/sha256"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPayloadArchivePutReadAndReuse(t *testing.T) {
	archive, err := NewPayloadArchive(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	body := bytes.Repeat([]byte("vless://fixture\n"), 4096)
	wantDigest := sha256.Sum256(body)

	first, err := archive.Put(context.Background(), body)
	if err != nil {
		t.Fatal(err)
	}
	second, err := archive.Put(context.Background(), body)
	if err != nil {
		t.Fatal(err)
	}
	wantKey := "sha256/" + strings.ToLower(strings.Join([]string{
		toHex(wantDigest[0]), toHex(wantDigest[1]), strings.ToLower(hexDigest(wantDigest)) + ".gz",
	}, "/"))
	if first.Key != wantKey || second.Key != wantKey {
		t.Fatalf("keys=%q/%q want=%q", first.Key, second.Key, wantKey)
	}
	if first.SHA256 != wantDigest || first.OriginalBytes != int64(len(body)) || first.CompressedBytes < 1 {
		t.Fatalf("reference=%+v", first)
	}
	got, err := archive.Read(context.Background(), first)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, body) {
		t.Fatal("archive body changed")
	}
	matches, err := filepath.Glob(filepath.Join(archive.Root(), "sha256", "*", "*", "*.gz"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 1 {
		t.Fatalf("archive files=%d", len(matches))
	}
}

func TestPayloadArchiveReadRejectsUnsafeOrCorruptFile(t *testing.T) {
	archive, err := NewPayloadArchive(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	body := []byte("trojan://fixture")
	ref, err := archive.Put(context.Background(), body)
	if err != nil {
		t.Fatal(err)
	}

	unsafe := ref
	unsafe.Key = "../../outside.gz"
	if _, err := archive.Read(context.Background(), unsafe); archiveErrorCode(err) != "archive_key_invalid" {
		t.Fatalf("unsafe key err=%v", err)
	}

	path := filepath.Join(archive.Root(), filepath.FromSlash(ref.Key))
	if err := os.WriteFile(path, []byte("not-gzip"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := archive.Read(context.Background(), ref); archiveErrorCode(err) != "archive_integrity" {
		t.Fatalf("corrupt file err=%v", err)
	}
}

func TestPayloadArchivePutRejectsSymlinkedHashDirectory(t *testing.T) {
	root := t.TempDir()
	archive, err := NewPayloadArchive(root)
	if err != nil {
		t.Fatal(err)
	}
	body := []byte("ss://symlink-fixture")
	digest := sha256.Sum256(body)
	outside := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "sha256"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "sha256", toHex(digest[0]))); err != nil {
		t.Fatal(err)
	}
	if _, err := archive.Put(context.Background(), body); archiveErrorCode(err) != "archive_path_unsafe" {
		t.Fatalf("symlink err=%v", err)
	}
}

func TestPayloadArchiveReadRejectsSymlinkedHashDirectory(t *testing.T) {
	root := t.TempDir()
	archive, err := NewPayloadArchive(root)
	if err != nil {
		t.Fatal(err)
	}
	ref, err := archive.Put(context.Background(), []byte("vmess://read-symlink-fixture"))
	if err != nil {
		t.Fatal(err)
	}
	firstHashDir := filepath.Join(root, "sha256", toHex(ref.SHA256[0]))
	outside := filepath.Join(t.TempDir(), "hash-dir")
	if err := os.Rename(firstHashDir, outside); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, firstHashDir); err != nil {
		t.Fatal(err)
	}
	if _, err := archive.Read(context.Background(), ref); archiveErrorCode(err) != "archive_path_unsafe" {
		t.Fatalf("symlink read err=%v", err)
	}
}

func TestInspectArchiveReportsReferencesAndQuarantineWithoutDeleting(t *testing.T) {
	archive, err := NewPayloadArchive(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	first, err := archive.Put(context.Background(), []byte("first"))
	if err != nil {
		t.Fatal(err)
	}
	second, err := archive.Put(context.Background(), []byte("second"))
	if err != nil {
		t.Fatal(err)
	}
	missing := first
	missing.SHA256 = sha256.Sum256([]byte("missing"))
	missing.Key = archiveKey(missing.SHA256)
	quarantine := t.TempDir()
	if err := os.WriteFile(filepath.Join(quarantine, "one.corrupt"), []byte("bad"), 0o600); err != nil {
		t.Fatal(err)
	}

	status, err := InspectArchive(context.Background(), archive, []ArchiveReference{first, missing}, quarantine)
	if err != nil {
		t.Fatal(err)
	}
	if status.Files != 2 || status.Referenced != 1 || status.Missing != 1 || status.Unreferenced != 1 || status.Corrupt != 0 {
		t.Fatalf("status=%+v second=%+v", status, second)
	}
	if status.QuarantineFiles != 1 || status.QuarantineBytes != 3 || status.Bytes < first.CompressedBytes+second.CompressedBytes {
		t.Fatalf("status=%+v", status)
	}
	if _, err := os.Stat(filepath.Join(quarantine, "one.corrupt")); err != nil {
		t.Fatalf("inspection changed quarantine: %v", err)
	}
}

func toHex(value byte) string {
	const digits = "0123456789abcdef"
	return string([]byte{digits[value>>4], digits[value&0x0f]})
}

func hexDigest(digest [32]byte) string {
	const digits = "0123456789abcdef"
	encoded := make([]byte, len(digest)*2)
	for i, value := range digest {
		encoded[i*2] = digits[value>>4]
		encoded[i*2+1] = digits[value&0x0f]
	}
	return string(encoded)
}
