package fhirtest

import (
	"fmt"
	"time"

	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/stablerand"
)

// stableRandomID generates a random ID value that will be stable across
// multiple test executions.
func stableRandomID() *datatypes_go_proto.Id {
	return &datatypes_go_proto.Id{
		Value: randomID(),
	}
}

// stableRandomVersionID generates a random version-ID value that will be
// stable across multiple test executions.
func stableRandomVersionID() *datatypes_go_proto.Id {
	return &datatypes_go_proto.Id{
		Value: randomVersionID(),
	}
}

func stableRandomInstant() *datatypes_go_proto.Instant {
	return fhir.Instant(stableRandomTime())
}

func StableRandomMeta() *datatypes_go_proto.Meta {
	return &datatypes_go_proto.Meta{
		LastUpdated: stableRandomInstant(),
		VersionId:   stableRandomVersionID(),
	}
}

func randomVersionID() string {
	const versionIDLength = 26
	return stablerand.AlnumString(versionIDLength)
}

// This is a different implementation than fhir.RandomID() since this test
// library manually sets a random seed.
func randomID() string {
	uuidBase := stablerand.HexString(32)
	return fmt.Sprintf(
		"%v-%v-%v-%v-%v",
		uuidBase[0:8],
		uuidBase[8:12],
		uuidBase[12:16],
		uuidBase[16:20],
		uuidBase[20:],
	)
}

func stableRandomTime() time.Time {
	const (
		// Timestamp for 2020-01-01 T12:00:00
		baseTime int64 = 1577898000

		// Variation of up to 1 year
		timeVariation = time.Hour * 24 * 365
	)

	base := time.Unix(baseTime, 0)
	return stablerand.Time(base, timeVariation)
}
