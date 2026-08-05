package store

import (
	"context"
	"crypto/sha256"
	"math/rand"
	"strings"
	"testing"
)

func TestCompressDatabasePayloadRejectsBoundedCompressedOverflow(t *testing.T) {
	body := make([]byte, 4096)
	if _, err := rand.New(rand.NewSource(7)).Read(body); err != nil {
		t.Fatal(err)
	}
	if _, err := compressDatabasePayload(body, 128); payloadStorageErrorCode(err) != "database_payload_too_large" {
		t.Fatalf("err=%v", err)
	}
	compressed, err := compressDatabasePayload([]byte(strings.Repeat("a", 4096)), 128)
	if err != nil {
		t.Fatal(err)
	}
	if len(compressed) > 128 {
		t.Fatalf("compressed bytes=%d", len(compressed))
	}
}

func TestValidateExternalPayload(t *testing.T) {
	body := []byte("vless://external")
	digest := sha256.Sum256(body)
	valid := &ExternalPayload{
		SHA256: digest, OriginalBytes: int64(len(body)), CompressedBytes: 42,
		Compression: "gzip", ArchiveKey: "sha256/ab/cd/body.gz",
	}
	if err := validateExternalPayload(valid); err != nil {
		t.Fatal(err)
	}

	for name, mutate := range map[string]func(*ExternalPayload){
		"missing_sha256":    func(payload *ExternalPayload) { payload.SHA256 = [32]byte{} },
		"missing_key":       func(payload *ExternalPayload) { payload.ArchiveKey = "" },
		"zero_original":     func(payload *ExternalPayload) { payload.OriginalBytes = 0 },
		"zero_compressed":   func(payload *ExternalPayload) { payload.CompressedBytes = 0 },
		"wrong_compression": func(payload *ExternalPayload) { payload.Compression = "zstd" },
	} {
		t.Run(name, func(t *testing.T) {
			payload := *valid
			mutate(&payload)
			if err := validateExternalPayload(&payload); err == nil || !strings.Contains(err.Error(), "external payload") {
				t.Fatalf("err=%v", err)
			}
		})
	}
}

func TestClaimUnparsedFetchesByStorageRejectsInvalidKindBeforeDatabase(t *testing.T) {
	if _, err := ClaimUnparsedFetchesByStorage(context.Background(), nil, "other", 1); err == nil || !strings.Contains(err.Error(), "storage kind") {
		t.Fatalf("err=%v", err)
	}
}
