package terminology

import (
	"context"

	pgp "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/parameters_go_proto"
)

type ValueSetValidateCodeOptions struct {
	// The value set OID or UUID.
	ID string
	// The code system ID, OID, or URI.
	System string
	// The code to be checked for validity.
	Code string
	// The effective date for determining validity, format should be
	// YYYY-MM-DD. If empty, the service will return a result based
	// on the latest dated system.
	Date string
	// The value set revision to validate against.
	ValueSetVersion string
}

type Service interface {
	ValueSetValidateCode(ctx context.Context, opts *ValueSetValidateCodeOptions) (*pgp.Parameters, error)
}
