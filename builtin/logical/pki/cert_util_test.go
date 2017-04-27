package pki

import (
	"fmt"
	"testing"

	"strings"

	"github.com/hashicorp/vault/logical"
)

func TestPki_FetchCertBySerial(t *testing.T) {
	storage := &logical.InmemStorage{}

	cases := map[string]struct {
		Req    *logical.Request
		Prefix string
		Serial string
	}{
		"valid cert": {
			&logical.Request{
				Operation: logical.ReadOperation,
				Storage:   storage,
			},
			"certs/",
			"00:00:00:00:00:00:00:00",
		},
		"revoked cert": {
			&logical.Request{
				Operation: logical.ReadOperation,
				Storage:   storage,
			},
			"revoked/",
			"11:11:11:11:11:11:11:11",
		},
	}

	// Test for colon-based paths in storage
	for name, tc := range cases {
		storageKey := fmt.Sprintf("%s%s", tc.Prefix, tc.Serial)
		err := storage.Put(&logical.StorageEntry{
			Key:   storageKey,
			Value: []byte("some data"),
		})
		if err != nil {
			t.Fatalf("error writing to storage on %s colon-based storage path: %s", name, err)
		}

		certEntry, err := fetchCertBySerial(tc.Req, tc.Prefix, tc.Serial)
		if err != nil {
			t.Fatalf("error on %s for colon-based storage path: %s", name, err)
		}

		// Check for non-nil on valid/revoked certs
		if certEntry == nil {
			t.Fatalf("nil on %s for colon-based storage path", name)
		}

		// Ensure that cert serials are converted/updated after fetch
		expectedKey := fmt.Sprintf("%s%s", tc.Prefix, strings.Replace(strings.ToLower(tc.Serial), ":", "-", -1))
		se, err := storage.Get(expectedKey)
		if err != nil {
			t.Fatalf("error on %s for colon-based storage path:%s", name, err)
		}
		if strings.Compare(expectedKey, se.Key) != 0 {
			t.Fatalf("expected: %s, got: %s", expectedKey, certEntry.Key)
		}
	}

	// Reset storage
	storage = &logical.InmemStorage{}

	// Test for hyphen-base paths in storage
	for name, tc := range cases {
		storageKey := fmt.Sprintf("%s%s", tc.Prefix, strings.Replace(strings.ToLower(tc.Serial), ":", "-", -1))
		err := storage.Put(&logical.StorageEntry{
			Key:   storageKey,
			Value: []byte("some data"),
		})
		if err != nil {
			t.Fatalf("error writing to storage on %s hyphen-based storage path: %s", name, err)
		}

		certEntry, err := fetchCertBySerial(tc.Req, tc.Prefix, tc.Serial)
		if err != nil {
			t.Fatalf("error on %s for hyphen-based storage path: %s", name, err)
		}

		// Check for non-nil on valid/revoked certs
		if certEntry == nil {
			t.Fatalf("nil on %s for hyphen-based storage path", name)
		}
	}
}
