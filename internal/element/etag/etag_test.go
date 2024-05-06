package etag_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/verily-src/fhirpath-go/internal/element/etag"
)

func TestVersionIDFromEtag(t *testing.T) {
	testCases := []struct {
		name          string
		etag          string
		wantVersionID string
		wantErr       error
	}{
		{
			name:          "valid etag",
			etag:          `W/"foo"`,
			wantVersionID: "foo",
			wantErr:       nil,
		},
		{
			name:    "empty etag",
			etag:    "",
			wantErr: etag.ErrInvalidEtagVersionID,
		},
		{
			name:    "etag without version ID",
			etag:    `W/""`,
			wantErr: etag.ErrInvalidEtagVersionID,
		},
		{
			name:    "etag with version ID with non-ASCII letters",
			etag:    `W/"?()"`,
			wantErr: etag.ErrInvalidEtagVersionID,
		},
		{
			name:    "etag with long version ID",
			etag:    `W/"12345678901234567890123456789012345678901234567890123456789012345"`,
			wantErr: etag.ErrInvalidEtagVersionID,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotVersionID, gotErr := etag.VersionIDFromEtag(tc.etag)
			if !cmp.Equal(gotErr, tc.wantErr, cmpopts.EquateErrors()) {
				t.Fatalf("VersionIDFromEtag(%s) error mismatch: got [%v], want [%v]", tc.name, gotErr, tc.wantErr)
			}
			if gotErr == nil && gotVersionID != tc.wantVersionID {
				t.Errorf("VersionIDFromEtag(%s) versionID mismatch: got [%s], want [%s]",
					tc.name, gotVersionID, tc.wantVersionID)
			}
		})
	}
}
