package etag

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	// FHIR version ID follows the [A-Za-z0-9\-\.]{1,64} pattern
	// https://build.fhir.org/resource-definitions.html#Meta.versionId
	versionIdInEtagRegexp   = regexp.MustCompile(`W\/"([A-Za-z0-9\-\.]{1,64})"`)
	ErrInvalidEtagVersionID = errors.New("invalid version ID in etag")
)

// Returns the version ID from an ETag string.
func VersionIDFromEtag(etag string) (string, error) {
	matches := versionIdInEtagRegexp.FindStringSubmatch(etag)
	// leftmost match of the regular expression in string itself and then the matches, if any.
	// We expect the second match to be the version ID.
	if len(matches) != 2 {
		return "", fmt.Errorf("%w: got version ID: %s", ErrInvalidEtagVersionID, etag)
	}
	return matches[1], nil
}
