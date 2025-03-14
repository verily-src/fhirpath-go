package bundle_test

import (
	"fmt"
	"strconv"
	"testing"

	cpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/codes_go_proto"
	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	bpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/binary_go_proto"
	bcrpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"
	ppb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/verily-src/fhirpath-go/internal/bundle"
	"github.com/verily-src/fhirpath-go/internal/containedresource"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/resource"
	"google.golang.org/protobuf/testing/protocmp"
)

var (
	iuURI               = fhir.URI("Patient/IU")
	patient             = &ppb.Patient{Id: fhir.ID("IU")}
	ifNoneExistFullOpt  = bundle.WithIfNoneExist(&dtpb.Identifier{System: fhir.URI("system"), Value: fhir.String("value")})
	ifNoneExistSysOpt   = bundle.WithIfNoneExist(&dtpb.Identifier{System: fhir.URI("system")})
	ifNoneExistValOpt   = bundle.WithIfNoneExist(&dtpb.Identifier{Value: fhir.String("value")})
	ifNoneExistEmptyOpt = bundle.WithIfNoneExist(&dtpb.Identifier{})

	patientPostEntry           = makePatientEntry(cpb.HTTPVerbCode_POST, patient, nil, "")
	patientPostEntryFullHeader = makePatientEntry(cpb.HTTPVerbCode_POST, patient, nil, "identifier=system|value")
	patientPostEntrySysHeader  = makePatientEntry(cpb.HTTPVerbCode_POST, patient, nil, "identifier=system|")
	patientPostEntryValHeader  = makePatientEntry(cpb.HTTPVerbCode_POST, patient, nil, "identifier=value")
	uriPostEntry               = makePatientEntry(cpb.HTTPVerbCode_POST, patient, iuURI, "")
	patientPutEntry            = makePatientEntry(cpb.HTTPVerbCode_PUT, patient, nil, "")
	uriPutEntry                = makePatientEntry(cpb.HTTPVerbCode_PUT, patient, iuURI, "")
)

func makeEntryForGet(requestUrl string, fullUrl string) *bcrpb.Bundle_Entry {
	return &bcrpb.Bundle_Entry{
		FullUrl: fhir.URI(fullUrl),
		Request: &bcrpb.Bundle_Entry_Request{
			Method: &bcrpb.Bundle_Entry_Request_MethodCode{
				Value: cpb.HTTPVerbCode_GET,
			},
			Url: fhir.URI(requestUrl),
		},
	}
}

func TestPostEntry(t *testing.T) {
	testCases := []struct {
		name      string
		resource  fhir.Resource
		opts      []bundle.EntryOption
		wantEntry *bcrpb.Bundle_Entry
	}{
		{"no options", patient, nil, patientPostEntry},
		{"empty uri", patient, []bundle.EntryOption{bundle.WithFullURL("")}, patientPostEntry},
		{"uri provided", patient, []bundle.EntryOption{bundle.WithFullURL(iuURI.GetValue())}, uriPostEntry},
		{"apply if-none-exist", patient, []bundle.EntryOption{ifNoneExistFullOpt}, patientPostEntryFullHeader},
		{"apply if-none-exist with identifier.system", patient, []bundle.EntryOption{ifNoneExistSysOpt}, patientPostEntrySysHeader},
		{"apply if-none-exist with identifier.value", patient, []bundle.EntryOption{ifNoneExistValOpt}, patientPostEntryValHeader},
		{"ignore if-none-exist with empty identifier", patient, []bundle.EntryOption{ifNoneExistEmptyOpt}, patientPostEntry},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			entry := bundle.NewPostEntry(tc.resource, tc.opts...)

			got, want := entry, tc.wantEntry
			if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
				t.Errorf("PostEntry(%s) mismatch (-want, +got):\n%s", tc.name, diff)
			}
		})
	}
}

func TestPutEntry(t *testing.T) {
	testCases := []struct {
		name      string
		resource  fhir.Resource
		opts      []bundle.EntryOption
		wantEntry *bcrpb.Bundle_Entry
	}{
		{"no options", patient, nil, patientPutEntry},
		{"empty uri", patient, []bundle.EntryOption{bundle.WithFullURL("")}, patientPutEntry},
		{"uri provided", patient, []bundle.EntryOption{bundle.WithFullURL(iuURI.GetValue())}, uriPutEntry},
		{"ignore if-none-exist", patient, []bundle.EntryOption{ifNoneExistFullOpt}, patientPutEntry},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			entry := bundle.NewPutEntry(tc.resource, tc.opts...)

			got, want := entry, tc.wantEntry
			if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
				t.Errorf("PutEntry(%s) mismatch (-want, +got):\n%s", tc.name, diff)
			}
		})
	}
}

func TestDeleteEntry(t *testing.T) {
	testCases := []struct {
		name      string
		inType    resource.Type
		inID      string
		opts      []bundle.EntryOption
		wantEntry *bcrpb.Bundle_Entry
	}{
		{
			name:   "no options",
			inType: resource.Patient,
			inID:   "1234",
			opts:   nil,
			wantEntry: &bcrpb.Bundle_Entry{
				Request: &bcrpb.Bundle_Entry_Request{
					Method: &bcrpb.Bundle_Entry_Request_MethodCode{
						Value: cpb.HTTPVerbCode_DELETE,
					},
					Url: &dtpb.Uri{
						Value: "Patient/1234",
					},
				},
			},
		},
		{
			name:   "with option id",
			inType: resource.Patient,
			inID:   "1234",
			opts:   []bundle.EntryOption{bundle.WithFullURL("test-full-url")},
			wantEntry: &bcrpb.Bundle_Entry{
				Request: &bcrpb.Bundle_Entry_Request{
					Method: &bcrpb.Bundle_Entry_Request_MethodCode{
						Value: cpb.HTTPVerbCode_DELETE,
					},
					Url: &dtpb.Uri{
						Value: "Patient/1234",
					},
				},
				FullUrl: &dtpb.Uri{
					Value: "test-full-url",
				},
			},
		},
		{
			name:   "ignore if-none-exist",
			inType: resource.Patient,
			inID:   "1234",
			opts:   []bundle.EntryOption{ifNoneExistFullOpt},
			wantEntry: &bcrpb.Bundle_Entry{
				Request: &bcrpb.Bundle_Entry_Request{
					Method: &bcrpb.Bundle_Entry_Request_MethodCode{
						Value: cpb.HTTPVerbCode_DELETE,
					},
					Url: &dtpb.Uri{
						Value: "Patient/1234",
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			entry := bundle.NewDeleteEntry(tc.inType, tc.inID, tc.opts...)

			got, want := entry, tc.wantEntry
			if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
				t.Errorf("NewDeleteEntry(%s) mismatch (-want, +got):\n%s", tc.name, diff)
			}
		})
	}
}

func TestPatchEntryFromBytes(t *testing.T) {

	patch := []byte(`[{"op":"add","path":"/active","value":true}]`)
	validIdentity, _ := resource.NewIdentity("Patient", "123", "0")

	testCases := []struct {
		name      string
		resID     *resource.Identity
		payload   []byte
		wantEntry *bcrpb.Bundle_Entry
		wantError error
	}{
		{
			name:    "Valid patch",
			resID:   validIdentity,
			payload: patch,
			wantEntry: &bcrpb.Bundle_Entry{
				Request: &bcrpb.Bundle_Entry_Request{
					Method: &bcrpb.Bundle_Entry_Request_MethodCode{
						Value: cpb.HTTPVerbCode_PATCH,
					},
					Url: &dtpb.Uri{
						Value: "Patient/123",
					},
				},
				Resource: &bcrpb.ContainedResource{
					OneofResource: &bcrpb.ContainedResource_Binary{
						Binary: &bpb.Binary{
							ContentType: &bpb.Binary_ContentTypeCode{Value: "application/json-patch+json"},
							Data:        &dtpb.Base64Binary{Value: patch},
						},
					},
				},
			},
		},
		{
			name:      "Invalid resource identity",
			resID:     &resource.Identity{},
			payload:   patch,
			wantError: bundle.ErrInvalidIdentity,
		},
		{
			name:      "Nil resource identity",
			resID:     nil,
			payload:   patch,
			wantError: bundle.ErrInvalidIdentity,
		},
		{
			name:      "Nil payload",
			resID:     validIdentity,
			payload:   nil,
			wantError: bundle.ErrMissingPayload,
		},
		{
			name:      "Empty payload",
			resID:     validIdentity,
			payload:   []byte{},
			wantError: bundle.ErrMissingPayload,
		},
		{
			name:      "Invalid payload - single op",
			resID:     validIdentity,
			payload:   []byte(`{"op":"add","path":"/active","value":true}`),
			wantError: bundle.ErrInvalidPayload,
		},
		{
			name:      "Invalid payload - wrong field",
			resID:     validIdentity,
			payload:   []byte(`[{"op":"add","url":"/active","value":true}]`),
			wantError: bundle.ErrInvalidPayload,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			entry, err := bundle.PatchEntryFromBytes(tc.resID, tc.payload)

			if tc.wantError == nil {
				got, want := entry, tc.wantEntry
				if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
					t.Errorf("PatchEntryFromBytes(%s) mismatch (-want, +got):\n%s", tc.name, diff)
				}
			} else {
				got, want := err, tc.wantError
				if !cmp.Equal(got, want, cmpopts.EquateErrors()) {
					t.Errorf("PatchEntryFromBytes(%s): unexpected error got %s, want %s", tc.name, got, want)
				}
			}
		})
	}
}

func TestCollectionEntry(t *testing.T) {
	wantEntry := &bcrpb.Bundle_Entry{
		Resource: containedresource.Wrap(patient),
		FullUrl:  &dtpb.Uri{Value: "Patient/IU"},
	}
	entry := bundle.NewCollectionEntry(patient)

	got, want := entry, wantEntry
	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Errorf("CollectionEntry mismatch (-want, +got):\n%s", diff)
	}
}

// Test helper to create a BundleEntry of a Patient
func makePatientEntry(method cpb.HTTPVerbCode_Value, res fhir.Resource, uri *dtpb.Uri, header string) *bcrpb.Bundle_Entry {
	requestURL := "Patient"
	if method == cpb.HTTPVerbCode_PUT {
		requestURL = resource.URIString(res)
	}
	entry := &bcrpb.Bundle_Entry{
		Resource: containedresource.Wrap(res),
		Request: &bcrpb.Bundle_Entry_Request{
			Method: &bcrpb.Bundle_Entry_Request_MethodCode{
				Value: method,
			},
			Url: &dtpb.Uri{
				Value: requestURL,
			},
		},
	}
	if uriString := uri.GetValue(); uriString != "" {
		entry.FullUrl = uri
	}
	if header != "" {
		entry.Request.IfNoneExist = fhir.String(header)
	}
	return entry
}

func TestEntryReference(t *testing.T) {
	entry := &bcrpb.Bundle_Entry{
		FullUrl: &dtpb.Uri{
			Value: "patient-full-url",
		},
		Resource: containedresource.Wrap(resource.New("Patient")),
	}
	wantRef := &dtpb.Reference{
		Type: &dtpb.Uri{
			Value: "Patient",
		},
		Reference: &dtpb.Reference_Uri{
			Uri: &dtpb.String{
				Value: "patient-full-url",
			},
		},
	}

	ref := bundle.EntryReference(entry)

	got, want := ref, wantRef
	if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
		t.Errorf("EntryReference: (-got +want):\n%v", diff)
	}
}

func TestSetEntryIfMatch(t *testing.T) {
	testCases := []struct {
		name        string
		entry       *bcrpb.Bundle_Entry
		wantVersion string
	}{
		{
			"does not override ifMatch with resource-derived value if it's already set",
			&bcrpb.Bundle_Entry{
				Request: &bcrpb.Bundle_Entry_Request{
					Method:  &bcrpb.Bundle_Entry_Request_MethodCode{Value: cpb.HTTPVerbCode_POST},
					IfMatch: fhir.String("already set"),
				},
				Resource: &bcrpb.ContainedResource{
					OneofResource: &bcrpb.ContainedResource_Patient{Patient: &ppb.Patient{Meta: &dtpb.Meta{VersionId: fhir.ID("derived")}}},
				},
			},
			"already set",
		},
		{
			"does not set ifMatch if there's no resource version",
			&bcrpb.Bundle_Entry{
				Request: &bcrpb.Bundle_Entry_Request{Method: &bcrpb.Bundle_Entry_Request_MethodCode{Value: cpb.HTTPVerbCode_POST}},
				Resource: &bcrpb.ContainedResource{
					OneofResource: &bcrpb.ContainedResource_Patient{Patient: &ppb.Patient{Meta: &dtpb.Meta{VersionId: fhir.ID("")}}},
				},
			},
			"",
		},
		{
			"does not set ifMatch with resource-derived value if the method is unspecified",
			&bcrpb.Bundle_Entry{
				Request: &bcrpb.Bundle_Entry_Request{},
				Resource: &bcrpb.ContainedResource{
					OneofResource: &bcrpb.ContainedResource_Patient{Patient: &ppb.Patient{Meta: &dtpb.Meta{VersionId: fhir.ID("derived")}}},
				},
			},
			"",
		},
		{
			"does not set ifMatch with resource-derived value if it's a GET",
			&bcrpb.Bundle_Entry{
				Request: &bcrpb.Bundle_Entry_Request{Method: &bcrpb.Bundle_Entry_Request_MethodCode{Value: cpb.HTTPVerbCode_GET}},
				Resource: &bcrpb.ContainedResource{
					OneofResource: &bcrpb.ContainedResource_Patient{Patient: &ppb.Patient{Meta: &dtpb.Meta{VersionId: fhir.ID("derived")}}},
				},
			},
			"",
		},
		{
			"does not set ifMatch with resource-derived value if it's a HEAD",
			&bcrpb.Bundle_Entry{
				Request: &bcrpb.Bundle_Entry_Request{Method: &bcrpb.Bundle_Entry_Request_MethodCode{Value: cpb.HTTPVerbCode_HEAD}},
				Resource: &bcrpb.ContainedResource{
					OneofResource: &bcrpb.ContainedResource_Patient{Patient: &ppb.Patient{Meta: &dtpb.Meta{VersionId: fhir.ID("derived")}}},
				},
			},
			"",
		},
		{
			"does not set ifMatch with resource-derived value if it's a POST",
			&bcrpb.Bundle_Entry{
				Request: &bcrpb.Bundle_Entry_Request{Method: &bcrpb.Bundle_Entry_Request_MethodCode{Value: cpb.HTTPVerbCode_POST}},
				Resource: &bcrpb.ContainedResource{
					OneofResource: &bcrpb.ContainedResource_Patient{Patient: &ppb.Patient{Meta: &dtpb.Meta{VersionId: fhir.ID("derived")}}},
				},
			},
			"",
		},
		{
			"sets ifMatch with resource-derived value if it's a PUT",
			&bcrpb.Bundle_Entry{
				Request: &bcrpb.Bundle_Entry_Request{Method: &bcrpb.Bundle_Entry_Request_MethodCode{Value: cpb.HTTPVerbCode_PUT}},
				Resource: &bcrpb.ContainedResource{
					OneofResource: &bcrpb.ContainedResource_Patient{Patient: &ppb.Patient{Meta: &dtpb.Meta{VersionId: fhir.ID("derived")}}},
				},
			},
			`W/"derived"`,
		},
		{
			"does not set ifMatch with resource-derived value if it's a PATCH",
			&bcrpb.Bundle_Entry{
				Request: &bcrpb.Bundle_Entry_Request{Method: &bcrpb.Bundle_Entry_Request_MethodCode{Value: cpb.HTTPVerbCode_PATCH}},
				Resource: &bcrpb.ContainedResource{
					OneofResource: &bcrpb.ContainedResource_Patient{Patient: &ppb.Patient{Meta: &dtpb.Meta{VersionId: fhir.ID("derived")}}},
				},
			},
			"",
		},
		{
			"does not set ifMatch with resource-derived value if it's a DELETE",
			&bcrpb.Bundle_Entry{
				Request: &bcrpb.Bundle_Entry_Request{Method: &bcrpb.Bundle_Entry_Request_MethodCode{Value: cpb.HTTPVerbCode_DELETE}},
				Resource: &bcrpb.ContainedResource{
					OneofResource: &bcrpb.ContainedResource_Patient{Patient: &ppb.Patient{Meta: &dtpb.Meta{VersionId: fhir.ID("derived")}}},
				},
			},
			"",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bundle.SetEntryIfMatch(tc.entry)

			gotVersion := tc.entry.GetRequest().GetIfMatch().GetValue()
			if diff := cmp.Diff(gotVersion, tc.wantVersion); diff != "" {
				t.Errorf("SetEntryIfMatch(%s) version got = %v, want = %v", tc.name, gotVersion, tc.wantVersion)
			}
		})
	}
}

func TestStatusCodeFromEntry(t *testing.T) {
	_, invalidStatusCodeErr := strconv.Atoi("inv")

	testCases := []struct {
		name      string
		entry     *bcrpb.Bundle_Entry
		wantCode  int
		wantError error
	}{
		{
			name: "response entry with valid code",
			entry: &bcrpb.Bundle_Entry{
				Response: &bcrpb.Bundle_Entry_Response{
					Status: fhir.String("200 OK"),
				},
			},
			wantCode:  200,
			wantError: nil,
		},
		{
			name: "response entry with missing code",
			entry: &bcrpb.Bundle_Entry{
				Response: &bcrpb.Bundle_Entry_Response{},
			},
			wantCode:  0,
			wantError: bundle.ErrRspBundleEntryMissingStatus,
		},
		{
			name: "response entry with short code",
			entry: &bcrpb.Bundle_Entry{
				Response: &bcrpb.Bundle_Entry_Response{
					Status: fhir.String("11"),
				},
			},
			wantCode:  0,
			wantError: bundle.ErrRspBundleEntryShortStatus,
		},
		{
			name: "response entry with invalid code",
			entry: &bcrpb.Bundle_Entry{
				Response: &bcrpb.Bundle_Entry_Response{
					Status: fhir.String("invalid"),
				},
			},
			wantError: fmt.Errorf("invalid (atoi) Bundle_Entry.response.status: %v", invalidStatusCodeErr),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotCode, gotError := bundle.StatusCodeFromEntry(tc.entry)

			if tc.wantError != nil {
				if gotError.Error() != tc.wantError.Error() {
					t.Errorf("StatusCodeFromEntry(%s) error got = %v, want = %v", tc.name, gotError, tc.wantError)
				}
			}

			if diff := cmp.Diff(gotCode, tc.wantCode); diff != "" {
				t.Errorf("StatusCodeFromEntry(%s) code got = %v, want = %v", tc.name, gotCode, tc.wantCode)
			}
		})
	}
}
