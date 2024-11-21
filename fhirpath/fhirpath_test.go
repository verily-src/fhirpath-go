package fhirpath_test

import (
	"errors"
	"testing"
	"time"

	cpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/codes_go_proto"
	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	drpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/document_reference_go_proto"
	epb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/encounter_go_proto"
	lpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/list_go_proto"
	mrpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/medication_request_go_proto"
	opb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/observation_go_proto"
	ppb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	prpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/practitioner_go_proto"
	tpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/task_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/shopspring/decimal"
	"github.com/verily-src/fhirpath-go/fhirpath"
	"github.com/verily-src/fhirpath-go/fhirpath/compopts"
	"github.com/verily-src/fhirpath-go/fhirpath/evalopts"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
	"github.com/verily-src/fhirpath-go/internal/element/extension"
	"github.com/verily-src/fhirpath-go/internal/element/reference"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/fhirconv"
	"google.golang.org/protobuf/testing/protocmp"
)

type evaluateTestCase struct {
	name            string
	inputPath       string
	inputCollection []fhir.Resource
	wantCollection  system.Collection
	compileOptions  []fhirpath.CompileOption
	evaluateOptions []fhirpath.EvaluateOption
}

var patientChu = &ppb.Patient{
	Id:     fhir.ID("123"),
	Active: fhir.Boolean(true),
	Gender: &ppb.Patient_GenderCode{
		Value: cpb.AdministrativeGenderCode_FEMALE,
	},
	BirthDate: fhir.MustParseDate("2000-03-22"),
	Telecom: []*dtpb.ContactPoint{
		{
			System: &dtpb.ContactPoint_SystemCode{Value: cpb.ContactPointSystemCode_PHONE},
		},
	},
	Name: []*dtpb.HumanName{
		{
			Use: &dtpb.HumanName_UseCode{
				Value: cpb.NameUseCode_NICKNAME,
			},
			Given:  []*dtpb.String{fhir.String("Senpai")},
			Family: fhir.String("Chu"),
		},
		{
			Use: &dtpb.HumanName_UseCode{
				Value: cpb.NameUseCode_OFFICIAL,
			},
			Given:  []*dtpb.String{fhir.String("Kang")},
			Family: fhir.String("Chu"),
		},
	},
	Contact: []*ppb.Patient_Contact{
		{
			Name: &dtpb.HumanName{
				Given:  []*dtpb.String{fhir.String("Senpai")},
				Family: fhir.String("Rodusek"),
			},
		},
	},
}
var fooExtension, _ = extension.FromElement("foourl", fhir.String("foovalue"))
var barExtension, _ = extension.FromElement("barurl", fhir.String("barvalue"))
var nameVoldemort = &dtpb.HumanName{
	Given: []*dtpb.String{
		fhir.String("Lord"),
	},
	Family: fhir.String("Voldemort"),
}
var patientVoldemort = &ppb.Patient{
	Id:     fhir.ID("123"),
	Active: fhir.Boolean(true),
	Gender: &ppb.Patient_GenderCode{
		Value: cpb.AdministrativeGenderCode_FEMALE,
	},
	Deceased: &ppb.Patient_DeceasedX{
		Choice: &ppb.Patient_DeceasedX_Boolean{
			Boolean: fhir.Boolean(true),
		},
	},
	MultipleBirth: &ppb.Patient_MultipleBirthX{
		Choice: &ppb.Patient_MultipleBirthX_Integer{
			Integer: fhir.Integer(int32(2)),
		},
	},
	Meta: &dtpb.Meta{
		Tag: []*dtpb.Coding{
			{
				Code: fhir.Code("#blessed"),
			},
		},
	},
	Name: []*dtpb.HumanName{nameVoldemort},
	Extension: []*dtpb.Extension{
		fooExtension,
		barExtension,
	},
}
var docRef = &drpb.DocumentReference{
	Status: &drpb.DocumentReference_StatusCode{
		Value: cpb.DocumentReferenceStatusCode_CURRENT,
	},
	Content: []*drpb.DocumentReference_Content{
		{
			Attachment: &dtpb.Attachment{
				ContentType: &dtpb.Attachment_ContentTypeCode{
					Value: "image",
				},
				Url:   fhir.URL("http://image"),
				Title: fhir.String("title"),
			},
		},
	},
}
var questionnaireRef, _ = reference.Typed("Questionnaire", "1234")
var obsWithRef = &opb.Observation{
	Meta: &dtpb.Meta{
		Extension: []*dtpb.Extension{
			{
				Url: fhir.URI("https://extension"),
				Value: &dtpb.Extension_ValueX{
					Choice: &dtpb.Extension_ValueX_Reference{
						Reference: questionnaireRef,
					},
				},
			},
		},
	},
	DerivedFrom: []*dtpb.Reference{
		questionnaireRef,
	},
}
var listWithNilRef = &lpb.List{
	Entry: []*lpb.List_Entry{
		{Item: &dtpb.Reference{Type: fhir.URI("Location")}},
	},
}

func testEvaluate(t *testing.T, testCases []evaluateTestCase) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			compiledExpression, err := fhirpath.Compile(tc.inputPath, tc.compileOptions...)
			if err != nil {
				t.Fatalf("Compiling \"%s\" returned unexpected error: %v", tc.inputPath, err)
			}

			got, err := compiledExpression.Evaluate(tc.inputCollection, tc.evaluateOptions...)

			if err != nil {
				t.Fatalf("Evaluating \"%s\" returned unexpected error: %v", tc.inputPath, err)
			}
			if diff := cmp.Diff(tc.wantCollection, got, protocmp.Transform()); diff != "" {
				t.Errorf("Evaluating \"%s\" returned unexpected diff (-want, +got)\n%s", tc.inputPath, diff)
			}
		})
	}
}

func TestEvaluate_PathSelection_ReturnsError(t *testing.T) {
	end := system.MustParseDateTime("@2016-01-01T12:22:33Z")
	task := makeTaskWithEndTime(end)

	testCases := []struct {
		name    string
		path    string
		input   fhir.Resource
		wantErr error
	}{
		{
			name:    "Invalid value_us field on DateTime",
			path:    "(Task.input.value as DataRequirement).dateFilter[0].value.end.value_us",
			input:   task,
			wantErr: fhirpath.ErrInvalidField,
		}, {
			name:    "Invalid timezone field on DateTime",
			path:    "(Task.input.value as DataRequirement).dateFilter[0].value.end.timezone",
			input:   task,
			wantErr: fhirpath.ErrInvalidField,
		}, {
			name:    "Field is not in correct casing, but exists",
			path:    "Patient.multiple_birth",
			input:   patientVoldemort,
			wantErr: fhirpath.ErrInvalidField,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sut, err := fhirpath.Compile(tc.path)
			if err != nil {
				t.Fatalf("fhirpath.Compile(%v): unexpected err: %v", tc.name, err)
			}

			_, err = sut.Evaluate([]fhir.Resource{tc.input})

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Errorf("fhirpath.Compile(%v): want err '%v', got '%v'", tc.name, got, want)
			}
		})
	}
}

func TestEvaluate_PathSelection_ReturnsResult(t *testing.T) {
	practitioner := &prpb.Practitioner{
		Name: []*dtpb.HumanName{nameVoldemort},
	}
	end := system.MustParseDateTime("@2016-01-01T12:22:33Z")
	task := makeTaskWithEndTime(end)

	testCases := []evaluateTestCase{
		{
			name:            "Patient.name.given returns given name",
			inputPath:       "Patient.name.given",
			inputCollection: []fhir.Resource{patientVoldemort},
			wantCollection:  system.Collection{fhir.String("Lord")},
		},
		{
			name:            "Patient.name returns name",
			inputPath:       "Patient.name",
			inputCollection: []fhir.Resource{patientVoldemort},
			wantCollection:  system.Collection{nameVoldemort},
		},
		{
			name:            "Extension with resource type returns extensions",
			inputPath:       "Patient.extension",
			inputCollection: []fhir.Resource{patientVoldemort},
			wantCollection:  system.Collection{fooExtension, barExtension},
		},
		{
			name:            "Extension without resource type returns extensions",
			inputPath:       "extension",
			inputCollection: []fhir.Resource{patientVoldemort},
			wantCollection:  system.Collection{fooExtension, barExtension},
		},
		{
			name:            "Patient.name returns empty on non-patient resource",
			inputPath:       "Patient.name",
			inputCollection: []fhir.Resource{practitioner},
			wantCollection:  system.Collection{},
		},
		{
			name:            "Accessing code field returns code",
			inputPath:       "Patient.gender",
			inputCollection: []fhir.Resource{patientVoldemort},
			wantCollection:  system.Collection{patientVoldemort.Gender},
		},
		{
			name:            "converts value field of primitive to System primitive",
			inputPath:       "Patient.name.given.value",
			inputCollection: []fhir.Resource{patientVoldemort},
			wantCollection:  system.Collection{system.String("Lord")},
		},
		{
			name:            "returns empty on non-existent field",
			inputPath:       "Patient.language",
			inputCollection: []fhir.Resource{patientVoldemort},
			wantCollection:  system.Collection{},
		},
		{
			name:            "returns value from a field with the _value suffix",
			inputPath:       "Encounter.class",
			inputCollection: []fhir.Resource{&epb.Encounter{ClassValue: fhir.Coding("class-system", "class-code")}},
			wantCollection:  system.Collection{fhir.Coding("class-system", "class-code")},
		},
		{
			name:      "value as Quantity returns fhir Quantity datatype",
			inputPath: "Observation.value as Quantity",
			inputCollection: []fhir.Resource{
				&opb.Observation{
					Value: &opb.Observation_ValueX{
						Choice: &opb.Observation_ValueX_Quantity{
							Quantity: &dtpb.Quantity{
								Value: fhir.Decimal(float64(22.2)),
							},
						},
					},
				},
			},
			wantCollection: []any{
				&dtpb.Quantity{
					Value: fhir.Decimal(float64(22.2)),
				},
			},
		},
		{
			name:      "Quantity with addition returns system.Quantity",
			inputPath: "Observation.value as Quantity + 2",
			inputCollection: []fhir.Resource{
				&opb.Observation{
					Value: &opb.Observation_ValueX{
						Choice: &opb.Observation_ValueX_Quantity{
							Quantity: &dtpb.Quantity{
								Value: fhir.Decimal(float64(22.2)),
							},
						},
					},
				},
			},
			wantCollection: []any{system.MustParseQuantity("24.2", "")},
		},
		{
			name:            "reference field returns Type/ID",
			inputPath:       "Observation.derivedFrom[0].reference",
			inputCollection: []fhir.Resource{obsWithRef},
			wantCollection:  system.Collection{fhir.String("Questionnaire/1234")},
		},
		{
			name:            "reference extension field returns Type/ID",
			inputPath:       "Observation.meta.extension('https://extension').value.reference",
			inputCollection: []fhir.Resource{obsWithRef},
			wantCollection:  system.Collection{fhir.String("Questionnaire/1234")},
		},
		{
			name:            "nil reference does not panic",
			inputPath:       "List.entry.item.where(type = 'Location').reference",
			inputCollection: []fhir.Resource{listWithNilRef},
			wantCollection:  system.Collection{},
		},
		{
			name:            "Valid access of time field",
			inputPath:       "(Task.input.value as DataRequirement).dateFilter[0].value.end.value",
			inputCollection: []fhir.Resource{task},
			wantCollection:  system.Collection{system.String(fhirconv.DateTimeToString(end.ToProtoDateTime()))},
		},
	}
	testEvaluate(t, testCases)
}

func makeTaskWithEndTime(end system.DateTime) *tpb.Task {
	start := system.MustParseDateTime("@2016-01-01T12:00:00Z")
	task := &tpb.Task{
		Input: []*tpb.Task_Parameter{
			{
				Value: &tpb.Task_Parameter_ValueX{
					Choice: &tpb.Task_Parameter_ValueX_DataRequirement{
						DataRequirement: &dtpb.DataRequirement{
							DateFilter: []*dtpb.DataRequirement_DateFilter{
								{
									Value: &dtpb.DataRequirement_DateFilter_ValueX{
										Choice: &dtpb.DataRequirement_DateFilter_ValueX_Period{
											Period: &dtpb.Period{
												Start: start.ToProtoDateTime(),
												End:   end.ToProtoDateTime(),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return task
}

func TestEvaluate_LegacyPathSelection_ReturnsResult(t *testing.T) {
	compileOptions := []fhirpath.CompileOption{compopts.Permissive()}
	end := system.MustParseDateTime("@2016-01-01T12:22:33Z")
	task := makeTaskWithEndTime(end)

	testCases := []evaluateTestCase{
		{
			name:            "Legacy: Evaluates ValueX fields and value_us fields",
			inputPath:       "(Task.input.value as DataRequirement).dateFilter[0].value.period.end.value_us",
			inputCollection: []fhir.Resource{task},
			wantCollection:  system.Collection{end},
			compileOptions:  compileOptions,
		},
	}
	testEvaluate(t, testCases)
}

func TestEvaluate_Literal_ReturnsLiteral(t *testing.T) {
	decimal := system.Decimal(decimal.NewFromFloat(1.450))
	date, _ := system.ParseDate("2023-05-30")
	time, _ := system.ParseTime("08:30:55.999")
	dateTime, _ := system.ParseDateTime("2023-06-14T13:48:55.555Z")
	quantity, _ := system.ParseQuantity("20", "years")

	testCases := []evaluateTestCase{
		{
			name:            "null literal returns empty collection",
			inputPath:       "{}",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{},
		},
		{
			name:            "boolean literal returns Boolean",
			inputPath:       "true",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "string literal returns escaped string",
			inputPath:       "'string test\\ 1\\''",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.String("string test 1'")},
		},
		{
			name:            "integer literal returns Integer",
			inputPath:       "23",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Integer(23)},
		},
		{
			name:            "decimal literal returns Decimal",
			inputPath:       "1.450",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{decimal},
		},
		{
			name:            "date literal returns Date",
			inputPath:       "@2023-05-30",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{date},
		},
		{
			name:            "time literal returns Time",
			inputPath:       "@T08:30:55.999",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{time},
		},
		{
			name:            "dateTime literal returns DateTime",
			inputPath:       "@2023-06-14T13:48:55.555Z",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{dateTime},
		},
		{
			name:            "quantity literal returns Quantity",
			inputPath:       "20 years",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{quantity},
		},
	}

	testEvaluate(t, testCases)
}

func TestEvaluate_ThisInvocation_Evaluates(t *testing.T) {
	testCases := []evaluateTestCase{
		{
			name:            "returns nickname with where()",
			inputPath:       "Patient.name.given.where($this = 'Senpai')",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{dtpb.String{Value: "Senpai"}},
		},
	}

	testEvaluate(t, testCases)
}

func TestEvaluate_Index_ReturnsIndex(t *testing.T) {
	nameOne := &dtpb.HumanName{
		Given: []*dtpb.String{
			fhir.String("Kobe"),
			fhir.String("Bean"),
		},
		Family: fhir.String("Bryant"),
	}
	nameTwo := &dtpb.HumanName{
		Given: []*dtpb.String{
			fhir.String("The"),
		},
		Family: fhir.String("Goat"),
	}
	patient := &ppb.Patient{
		Name: []*dtpb.HumanName{
			nameOne,
			nameTwo,
		},
	}

	testCases := []evaluateTestCase{
		{
			name:            "first index returns result",
			inputPath:       "Patient.name[0]",
			inputCollection: []fhir.Resource{patient},
			wantCollection:  system.Collection{nameOne},
		},
		{
			name:            "second index returns result",
			inputPath:       "Patient.name[1]",
			inputCollection: []fhir.Resource{patient},
			wantCollection:  system.Collection{nameTwo},
		},
		{
			name:            "indexing name.given",
			inputPath:       "Patient.name.given[2]",
			inputCollection: []fhir.Resource{patient},
			wantCollection:  system.Collection{fhir.String("The")},
		},
		{
			name:            "indexing multiple times",
			inputPath:       "Patient.name[0].given[1]",
			inputCollection: []fhir.Resource{patient},
			wantCollection:  system.Collection{fhir.String("Bean")},
		},
		{
			name:            "out of bounds index returns empty",
			inputPath:       "Patient.name.given[5]",
			inputCollection: []fhir.Resource{patient},
			wantCollection:  system.Collection{},
		},
		{
			name:            "empty collection index returns empty",
			inputPath:       "Patient.name.given[{}]",
			inputCollection: []fhir.Resource{patient},
			wantCollection:  system.Collection{},
		},
	}

	testEvaluate(t, testCases)
}

func TestEvaluateEquality_ReturnsBoolean(t *testing.T) {
	request := &mrpb.MedicationRequest{
		Intent: &mrpb.MedicationRequest_IntentCode{Value: cpb.MedicationRequestIntentCode_FILLER_ORDER},
	}

	testCases := []evaluateTestCase{
		{
			name:            "querying active field",
			inputPath:       "Patient.active = true",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "two contrary conditions on 2 resources with an OR, first one true",
			inputPath:       "Patient.active = true or Observation.status = 'final'",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "two contrary conditions on 2 resources with an OR, second one true",
			inputPath:       "Observation.status = 'final' or Patient.active = true",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "inverse of active field",
			inputPath:       "Patient.active != true",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(false)},
		},
		{
			name:            "querying given name",
			inputPath:       "Patient.name[0].given = 'Senpai'",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "equality of complex types",
			inputPath:       "Patient.name[0].given = Patient.contact.name.given",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "complex types not equal",
			inputPath:       "Patient.name.family != Patient.contact.name.family",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "comparing non-equal fields",
			inputPath:       "Patient.name.family = Patient.contact.name.family",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(false)},
		},
		{
			name:            "comparing non-existent field",
			inputPath:       "Patient.maritalStatus = false",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{},
		},
		{
			name:            "comparing dates",
			inputPath:       "Patient.birthDate = @2000-03-22",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "compare date with dateTime",
			inputPath:       "@2012-12-31 = @2012-12-31T",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "compare dateTime with date",
			inputPath:       "@2012-12-31T = @2012-12-31",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "comparing non-equal dates",
			inputPath:       "@2000-01-02 != @2000-01-01",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "comparing mismatched date precision",
			inputPath:       "@2000-01 = @2000-01-03",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{},
		},
		{
			name:            "comparing mismatched date precision that isn't equal",
			inputPath:       "@2000-02 = @2000-01-03",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{system.Boolean(false)},
		},
		{
			name:            "respects timezones for DateTime comparison",
			inputPath:       "@2000-02-01T12:30:00Z = @2000-02-01T13:30:00+01:00",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "comparing gender code",
			inputPath:       "Patient.gender = 'female'",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "comparing code with non-enum value",
			inputPath:       "DocumentReference.content[0].attachment.contentType = 'image'",
			inputCollection: []fhir.Resource{docRef},
			wantCollection:  []any{system.Boolean(true)},
		},
		{
			name:            "comparing name use code",
			inputPath:       "Patient.name[0].use = 'nickname'",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "comparing telecom system code",
			inputPath:       "telecom.system = 'phone'",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "comparing incorrect telecom code",
			inputPath:       "telecom.system = 'carrier pigeon'",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(false)},
		},
		{
			name:            "mismatched case on code",
			inputPath:       "Patient.name.use = 'NICKNAME'",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(false)},
		},
		{
			name:            "comparing multi-word code",
			inputPath:       "MedicationRequest.intent = 'filler-order'",
			inputCollection: []fhir.Resource{request},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "mismatched case for multi-word code",
			inputPath:       "MedicationRequest.intent = 'fillerOrder'",
			inputCollection: []fhir.Resource{request},
			wantCollection:  system.Collection{system.Boolean(false)},
		},
		{
			name:            "comparing decimal to integer",
			inputPath:       "1 = 1.000",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "comparing decimal to quantity",
			inputPath:       "24.3 = 24.3 'kg'",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "comparing integer to quantity",
			inputPath:       "2 = 2.0 'lbs'",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
	}

	testEvaluate(t, testCases)
}

func TestParenthesizedExpression_MaintainsPrecedence(t *testing.T) {
	patient := &ppb.Patient{
		Name: []*dtpb.HumanName{
			{
				Given: []*dtpb.String{
					fhir.String("Alex"),
					fhir.String("Jon"),
					fhir.String("Matt"),
					fhir.String("Heming"),
				},
			},
		},
	}
	testCases := []evaluateTestCase{
		{
			name:            "evaluates parenthesized equality first",
			inputPath:       "true = (false = false)",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "evaluates parenthesized expressions in order",
			inputPath:       "true = ('Alex' = (name.given[0]))",
			inputCollection: []fhir.Resource{patient},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
	}

	testEvaluate(t, testCases)
}

func TestFunctionInvocation_Evaluates(t *testing.T) {
	testTime := time.Now()
	testDateTime, _ := system.DateTimeFromProto(fhir.DateTime(testTime))
	testCases := []evaluateTestCase{
		{
			name:            "returns nickname with where()",
			inputPath:       "Patient.name.where(use = 'nickname')",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{patientChu.Name[0]},
		},
		{
			name:            "returns official name with where()",
			inputPath:       "Patient.name.where(use = 'official')",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{patientChu.Name[1]},
		},
		{
			name:            "returns true with exists()",
			inputPath:       "Patient.name.exists(use = 'official')",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "returns false with exists()",
			inputPath:       "Patient.name.exists(use = 'random-use')",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(false)},
		},
		{
			name:            "returns true with exists() with BooleanExpression",
			inputPath:       "Patient.name.exists(use = 'official' and given = 'Kang')",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "returns true with exists() with BooleanExpression",
			inputPath:       "Patient.name.exists(use = 'random-use' or given = 'Kang')",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "returns true with where() and exists()",
			inputPath:       "Patient.name.where(use = 'official').exists()",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "chaining exists() is fine when the first exists() evaluates to true",
			inputPath:       "Patient.name.where(use = 'official').exists().exists().exists()",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "chaining exists() when the first exists() evaluates to false gives correct but ambiguous result",
			inputPath:       "Patient.name.where(use = 'random').exists().exists()",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "chaining empty() gives correct but ambiguous result",
			inputPath:       "Patient.name.where(use = 'random').empty().empty()",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(false)},
		},
		{
			name:            "returns false with where() and empty()",
			inputPath:       "Patient.name.where(use = 'official').empty()",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(false)},
		},
		{
			name:            "returns true with where() and empty()",
			inputPath:       "Patient.name.where(use = 'random').empty()",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "returns empty if no elements match where condition",
			inputPath:       "Patient.name.where(family = 'Suresh')",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{},
		},
		{
			name:            "evaluates timeOfDay() based on context, not dependent on latent factors",
			inputPath:       "timeOfDay().delay() = timeOfDay()",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(true)},
			compileOptions: []fhirpath.CompileOption{compopts.AddFunction("delay", func(in system.Collection) (system.Collection, error) {
				time.Sleep(time.Second * 2)
				return in, nil
			})},
		},
		{
			name:            "evaluates now() based on context, not dependent on latent factors",
			inputPath:       "now().delay() = now()",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(true)},
			compileOptions: []fhirpath.CompileOption{compopts.AddFunction("delay", func(in system.Collection) (system.Collection, error) {
				time.Sleep(time.Second * 2)
				return in, nil
			})},
		},
		{
			name:            "evaluates today() based on context, not dependent on latent factors",
			inputPath:       "today().delay() = today()",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(true)},
			compileOptions: []fhirpath.CompileOption{compopts.AddFunction("delay", func(in system.Collection) (system.Collection, error) {
				time.Sleep(time.Second * 2)
				return in, nil
			})},
		},
		{
			name:            "evaluates now() using overridden time",
			inputPath:       "now()",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{testDateTime},
			evaluateOptions: []fhirpath.EvaluateOption{evalopts.OverrideTime(testTime)},
		},
		{
			name:            "evaluate with custom function 'patient()'",
			inputPath:       "patient() = Patient",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(true)},
			compileOptions: []fhirpath.CompileOption{compopts.AddFunction("patient", func(system.Collection) (system.Collection, error) {
				return system.Collection{patientChu}, nil
			})},
		},
		{
			name:            "evaluate with custom function startsWith()",
			inputPath:       "Patient.name[0].family.startsWith('Ch')",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "evaluate with custom function endsWith()",
			inputPath:       "Patient.name[0].family.endsWith('hu')",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "evaluate with custom function length()",
			inputPath:       "Patient.name[0].family.length()",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Integer(3)},
		},
		{
			name:            "evaluate with custom function upper()",
			inputPath:       "Patient.name[0].given.upper()",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.String("SENPAI")},
		},
		{
			name:            "evaluate with custom function lower()",
			inputPath:       "Patient.name[0].family.lower()",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.String("chu")},
		},
		{
			name:            "evaluate with custom function contains()",
			inputPath:       "Patient.name[0].given.contains('pai')",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "evaluate with custom function toChars()",
			inputPath:       "Patient.name[0].family.toChars()",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection: system.Collection{
				system.String('C'),
				system.String('h'),
				system.String('u'),
			},
		},
		{
			name:            "evaluate with custom function substring()",
			inputPath:       "Patient.name[0].given.substring(1, 4)",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.String("enpa")},
		},
		{
			name:            "evaluate with custom function indexOf()",
			inputPath:       "Patient.name[0].given.indexOf('pa')",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Integer(3)},
		},
		{
			name:            "evaluate with custom function matches()",
			inputPath:       "Patient.name[0].family.matches('^[A-Za-z]*$')",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "evaluate with custom function replace()",
			inputPath:       "Patient.name[0].given.replace('Senpai', 'Oppa')",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.String("Oppa")},
		},
		{
			name:            "evaluate with custom function replaceMatches()",
			inputPath:       "Patient.name[0].family.replaceMatches('[A-Z]', 'zzz')",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.String("zzzhu")},
		},
		{
			name:            "returns full name with select()",
			inputPath:       "Patient.name.where(use = 'official').select(given.first() + ' ' + family)",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.String("Kang Chu")},
		},
		{
			name:            "projection on given name with select()",
			inputPath:       "name.given.select($this = 'Kang')",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(false), system.Boolean(true)},
		},
		{
			name:            "returns concatenated family name value with join()",
			inputPath:       "name.family.value.join('-')",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.String("Chu-Chu")},
			compileOptions:  []fhirpath.CompileOption{compopts.WithExperimentalFuncs()},
		},
		{
			name:            "returns concatenated family name with join()",
			inputPath:       "name.family.join('-')",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.String("Chu-Chu")},
			compileOptions:  []fhirpath.CompileOption{compopts.WithExperimentalFuncs()},
		},
	}

	testEvaluate(t, testCases)
}

func TestTypeExpression_Evaluates(t *testing.T) {
	testCases := []evaluateTestCase{
		{
			name:            "returns true for resource type check",
			inputPath:       "Patient is Patient",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "returns true for resource subtype relationship",
			inputPath:       "Patient is FHIR.Resource",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "returns true for primitive type check",
			inputPath:       "Patient.deceased is boolean",
			inputCollection: []fhir.Resource{patientVoldemort},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "returns false for primitive type case mismatch",
			inputPath:       "Patient.deceased is Boolean",
			inputCollection: []fhir.Resource{patientVoldemort},
			wantCollection:  system.Collection{system.Boolean(false)},
		},
		{
			name:            "returns true for system type check",
			inputPath:       "1 is Integer",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "propagates empty collection",
			inputPath:       "{} is Boolean",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{},
		},
		{
			name:            "passes through as expression",
			inputPath:       "Patient as Patient",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{patientChu},
		},
		{
			name:            "passes through as expression for subtype relationship",
			inputPath:       "Patient.name.use[0] as FHIR.Element",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{patientChu.Name[0].Use},
		},
		{
			name:            "returns empty if as expression is not correct type",
			inputPath:       "Patient.name.family[0] as HumanName",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{},
		},
		{
			name:            "unwraps polymorphic type with as expression",
			inputPath:       "Patient.deceased as boolean",
			inputCollection: []fhir.Resource{patientVoldemort},
			wantCollection:  system.Collection{fhir.Boolean(true)},
		},
		{
			name:            "passes through system type with as expression",
			inputPath:       "@2000-12-05 as Date",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{system.MustParseDate("2000-12-05")},
		},
	}

	testEvaluate(t, testCases)
}

func TestBooleanExpression_Evaluates(t *testing.T) {
	testCases := []evaluateTestCase{
		{
			name:            "evaluates and correctly",
			inputPath:       "true and false",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{system.Boolean(false)},
		},
		{
			name:            "evaluates boolean correctly with protos",
			inputPath:       "Patient.active and Patient.deceased",
			inputCollection: []fhir.Resource{patientVoldemort},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "evaluates or correctly",
			inputPath:       "true or false",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "propogates empty collections correctly",
			inputPath:       "false or {}",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{},
		},
		{
			name:            "evaluates xor correctly",
			inputPath:       "true xor true",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{system.Boolean(false)},
		},
		{
			name:            "evaluates implies correctly",
			inputPath:       "false implies false",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "not function inverts input",
			inputPath:       "Patient.active.not()",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(false)},
		},
	}

	testEvaluate(t, testCases)
}

func TestComparisonExpression_ReturnsBool(t *testing.T) {
	testCases := []evaluateTestCase{
		{
			name:            "compares strings",
			inputPath:       "'abc' > 'ABC'",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "compares integer with decimal",
			inputPath:       "4 <= 4.0",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "compares quantities of the same precision",
			inputPath:       "3.2 'kg' > 9.7 'kg'",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{system.Boolean(false)},
		},
		{
			name:            "returns empty for quantities of different precision",
			inputPath:       "99.9 'cm' < 1 'm'",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{},
		},
		{
			name:            "compares dates correctly",
			inputPath:       "@2018-03-01 >= @2018-03-01",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "returns empty for mismatched time precision",
			inputPath:       "@T08:30 > @T08:30:00",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{},
		},
		{
			name:            "correctly compares times",
			inputPath:       "@T10:29:59.999 < @T10:30:00",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "validate the age of an individual",
			inputPath:       "Patient.birthDate + 23 'years' <= today()",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
	}

	testEvaluate(t, testCases)
}

func TestArithmetic_ReturnsResult(t *testing.T) {
	testCases := []evaluateTestCase{
		{
			name:            "adds dates with quantity",
			inputPath:       "@2012-12-12 + 12 days",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{system.MustParseDate("2012-12-24")},
		},
		{
			name:            "concatenates strings",
			inputPath:       "'hello ' & 'world'",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{system.String("hello world")},
		},
		{
			name:            "subtracts integer from quantity",
			inputPath:       "8 'kg' - 4",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{system.MustParseQuantity("4", "kg")},
		},
		{
			name:            "multiplies values together",
			inputPath:       "8 * 4.2",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{system.Decimal(decimal.NewFromFloat(33.6))},
		},
		{
			name:            "divides values",
			inputPath:       "8 / 2.5",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{system.Decimal(decimal.NewFromFloat(3.2))},
		},
		{
			name:            "performs floor division",
			inputPath:       "29 div 10",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{system.Integer(2)},
		},
		{
			name:            "performs modulo operation",
			inputPath:       "100 mod 11",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{system.Integer(1)},
		},
	}

	testEvaluate(t, testCases)
}

func TestCompile_ReturnsError(t *testing.T) {
	testCases := []struct {
		name           string
		inputPath      string
		compileOptions []fhirpath.CompileOption
	}{
		{
			name:      "mismatched parentheses",
			inputPath: "Patient.name.where(use = official",
		},
		{
			name:      "invalid character",
			inputPath: "Patient.*name",
		},
		{
			name:      "invalid expression (misspelling)",
			inputPath: "Patient.name aand Patient.name",
		},
		{
			name:      "invalid expression (non-existent operator)",
			inputPath: "Patient.name nor Patient.name",
		},
		{
			name:      "invalid character (lexer error)",
			inputPath: "Patient^",
		},
		{
			name:      "non-existent function",
			inputPath: "Patient.notAFunc()",
		},
		{
			name:           "expanding function table with bad function",
			inputPath:      "Patient.badFn()",
			compileOptions: []fhirpath.CompileOption{compopts.AddFunction("badFn", func() {})},
		},
		{
			name:           "attempting to override existing function",
			inputPath:      "Patient.where()",
			compileOptions: []fhirpath.CompileOption{compopts.AddFunction("where", func(system.Collection) (system.Collection, error) { return nil, nil })},
		},
		{
			name:      "evaluating function with mismatched arity",
			inputPath: "Patient.name.where(use = 'official', use = 'nickname')",
		},
		{
			name:      "evaluating function with invalid arguments",
			inputPath: "Patient.name.where(invalid $ expr)",
		},
		{
			name:      "resolving invalid type specifier",
			inputPath: "1 is System.Patient",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := fhirpath.Compile(tc.inputPath, tc.compileOptions...); err == nil {
				t.Errorf("Compiling \"%s\" doesn't raise error when expected to", tc.inputPath)
			}
		})
	}
}

func TestEvaluate_ReturnsError(t *testing.T) {
	alwaysFails := func(system.Collection) (system.Collection, error) {
		return nil, errors.New("some error")
	}
	testCases := []struct {
		name            string
		inputPath       string
		inputCollection []fhir.Resource
		compileOptions  []fhirpath.CompileOption
		evaluateOptions []fhirpath.EvaluateOption
	}{
		{
			name:            "non-integer index returns error",
			inputPath:       "Patient.name['not a number']",
			inputCollection: []fhir.Resource{patientChu},
		},
		{
			name:            "evaluating failing function propagates error",
			inputPath:       "alwaysFails()",
			inputCollection: []fhir.Resource{},
			compileOptions:  []fhirpath.CompileOption{compopts.AddFunction("alwaysFails", alwaysFails)},
		},
		{
			name:            "evaluating is expression on non-singleton collection",
			inputPath:       "Patient.name is string",
			inputCollection: []fhir.Resource{patientChu},
		},
		{
			name:            "comparing unsupported types",
			inputPath:       "true > 0",
			inputCollection: []fhir.Resource{},
		},
		{
			name:            "arithmetic on unsupported types",
			inputPath:       "1 + true",
			inputCollection: []fhir.Resource{},
		},
		{
			name:            "misspelled identifier raises error",
			inputPath:       "Patient.nam.given",
			inputCollection: []fhir.Resource{patientVoldemort},
		},
		{
			name:            "overriding existing constant",
			inputPath:       "'valid fhirpath'",
			inputCollection: []fhir.Resource{},
			evaluateOptions: []fhirpath.EvaluateOption{
				evalopts.EnvVariable("context", system.String("context")),
			},
		},
		{
			name:            "adding unsupported type as constant",
			inputPath:       "%var",
			inputCollection: []fhir.Resource{},
			evaluateOptions: []fhirpath.EvaluateOption{
				evalopts.EnvVariable("var", 1),
			},
		},
		{
			name:            "adding unsupported type within collection as constant",
			inputPath:       "%collection",
			inputCollection: []fhir.Resource{},
			evaluateOptions: []fhirpath.EvaluateOption{
				evalopts.EnvVariable("collection", system.Collection{system.Integer(1), 1}),
			},
		},
		{
			name:            "negating unsupported type",
			inputPath:       "-'string'",
			inputCollection: []fhir.Resource{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expression, err := fhirpath.Compile(tc.inputPath, tc.compileOptions...)
			if err != nil {
				t.Fatalf("compiling \"%s\" raised unexpected error: %v", tc.inputPath, err)
			}
			if _, err = expression.Evaluate(tc.inputCollection, tc.evaluateOptions...); err == nil {
				t.Errorf("Evaluating expression \"%s\" doesn't raise error when expected to", tc.inputPath)
			}
		})
	}
}

func TestExternalConstantExpression_ReturnsConstant(t *testing.T) {
	testCases := []evaluateTestCase{
		{
			name:            "system type constant",
			inputPath:       "%var",
			inputCollection: []fhir.Resource{},
			evaluateOptions: []fhirpath.EvaluateOption{
				evalopts.EnvVariable("var", system.String("hello")),
			},
			wantCollection: system.Collection{system.String("hello")},
		},
		{
			name:            "proto type constant",
			inputPath:       "%patient",
			inputCollection: []fhir.Resource{},
			evaluateOptions: []fhirpath.EvaluateOption{
				evalopts.EnvVariable("patient", patientChu),
			},
			wantCollection: system.Collection{patientChu},
		},
		{
			name:            "collection constant containing system and proto types",
			inputPath:       "%collection",
			inputCollection: []fhir.Resource{},
			evaluateOptions: []fhirpath.EvaluateOption{
				evalopts.EnvVariable("collection", system.Collection{system.String("hello"), patientChu}),
			},
			wantCollection: system.Collection{system.String("hello"), patientChu},
		},
		{
			name:            "returns input as %context variable",
			inputPath:       "%context",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{patientChu},
		},
		{
			name:            "returns ucum url as %ucum",
			inputPath:       "%ucum",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{system.String("http://unitsofmeasure.org")},
		},
	}

	testEvaluate(t, testCases)
}

func TestPolarityExpression(t *testing.T) {
	testCases := []evaluateTestCase{
		{
			name:            "negates integer",
			inputPath:       "-1",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{system.Integer(-1)},
		},
		{
			name:            "does nothing when using '+'",
			inputPath:       "+2.45",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{system.Decimal(decimal.NewFromFloat(2.45))},
		},
		{
			name:            "negates field from proto",
			inputPath:       "-(Patient.multipleBirth as integer)",
			inputCollection: []fhir.Resource{patientVoldemort},
			wantCollection:  system.Collection{system.Integer(-2)},
		},
		{
			name:            "performs arithmetic correctly with negatives",
			inputPath:       "-1 - (-2)",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{system.Integer(1)},
		},
	}

	testEvaluate(t, testCases)
}

func TestAll_Evaluates(t *testing.T) {
	testCases := []evaluateTestCase{
		{
			name:            "returns false if not all elements are integers",
			inputPath:       "Patient.name.given.all($this is Integer)",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(false)},
		},
		{
			name:            "returns true if input is empty",
			inputPath:       "{}.all($this is Integer)",
			inputCollection: []fhir.Resource{},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "returns true if born during the 21st century",
			inputPath:       "Patient.birthDate.all($this >= @2000-01-01 and $this < @2100-01-01)",
			inputCollection: []fhir.Resource{patientChu},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
	}

	testEvaluate(t, testCases)
}

func TestMustCompile_CompileError_Panics(t *testing.T) {
	defer func() { _ = recover() }()

	fhirpath.MustCompile("Patient.name.where(use = official")

	t.Errorf("MustCompile: Expected panic")
}

func TestMustCompile_ValidExpression_ReturnsExpression(t *testing.T) {
	result := fhirpath.MustCompile("Patient.name")

	if result == nil {
		t.Errorf("MustCompile: Expected result")
	}
}
