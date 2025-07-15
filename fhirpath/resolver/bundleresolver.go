package resolver

import (
	"fmt"
	"regexp"

	cpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/codes_go_proto"
	r4pb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"
	"github.com/verily-src/fhirpath-go/internal/containedresource"
	"github.com/verily-src/fhirpath-go/internal/element/reference"
	"github.com/verily-src/fhirpath-go/internal/fhir"
)

var (
	ErrMultipleResourcesWithSameIDAndVersion = fmt.Errorf("multiple bundle entries with same id and version found for the same resource type")
	ErrUnsupportedBundleType                 = fmt.Errorf("only bundles of type SEARCHSET or COLLECTION are supported by BundleResolver, another found")
	ErrNilBundleInit                         = fmt.Errorf("cannot initialize BundleResolver with a nil bundle")
	ErrMissingMetaOrLastUpdated              = fmt.Errorf("undefined behavior, multiple versionless absolute reference matches found, with missing meta or meta.lastUpdated fields")
)

// BundleResolver implements the `Resolver` interface and
// resolves FHIR references (URIs or plain URLs) to resources within a FHIR Bundle.
type BundleResolver struct {
	bundle *r4pb.Bundle
}

// NewBundleResolver initializes a BundleResolver for SEARCHSET and COLLECTION type bundles
func NewBundleResolver(bundle *r4pb.Bundle) (*BundleResolver, error) {
	if bundle == nil {
		return nil, ErrNilBundleInit
	}
	bundleType := bundle.GetType().GetValue()
	if bundleType != cpb.BundleTypeCode_SEARCHSET && bundleType != cpb.BundleTypeCode_COLLECTION {
		return nil, ErrUnsupportedBundleType
	}
	return &BundleResolver{
		bundle: bundle,
	}, nil
}

// resolvableResources holds maps and slices of resources extracted from a FHIR Bundle
// to facilitate efficient lookup and resolution of FHIR references.
type resolvableResources struct {
	// urnResourceMap stores resources accessible by their URN (e.g., "urn:uuid:..." or "urn:oid:...").
	// The key is the URN string, and the value is the corresponding fhir.Resource.
	urnResourceMap map[string]fhir.Resource
	// versionlessUrlResourcesMap stores resources accessible by their absolute URL without a version ID.
	// The key is a concatenation of the base URL and the resource's identity string (e.g., "http://example.com/Patient/123").
	// The value is a slice of fhir.Resource.
	versionlessUrlResourcesMap map[string][]fhir.Resource
	// versionedURLResourcesMap stores resources accessible by their absolute URL including a version ID.
	// The key is a concatenation of the base URL, resource's identity string, and version ID (e.g., "http://example.com/Patient/123/_history/4").
	// The value is a slice of fhir.Resource.
	versionedURLResourcesMap map[string][]fhir.Resource
	// rootURLs stores a list of base URLs found in the bundle's entries that are absolute RESTful URIs.
	// These are used to resolve relative FHIR references by constructing full URLs.
	rootURLs []string
}

// fromBundle populates a resolvableResources struct from a FHIR Bundle.
// It populates maps for resolving resources by their full URL and versioned URL,
// and extracts base URLs from absolute RESTful URIs within the bundle entries.
// Invalid entries or entries without a full URL are skipped.
func (rr *resolvableResources) fromBundle(bundle *r4pb.Bundle) {
	for _, entry := range bundle.Entry {
		if entry == nil || entry.Resource == nil {
			continue
		}
		entryResource := containedresource.Unwrap(entry.Resource)

		// if the FullUrl is missing, this resource cannot be resolved
		fullURL := entry.GetFullUrl().GetValue()
		if fullURL == "" {
			continue
		}

		if isURN(fullURL) {
			rr.urnResourceMap[fullURL] = entryResource
			continue
		}

		_, _, isVersioned, isAbsoluteURL := absoluteURLInfo(fullURL)
		if isAbsoluteURL {
			key := fullURL
			if isVersioned {
				rr.versionedURLResourcesMap[key] = append(rr.versionedURLResourcesMap[key], entryResource)
			} else {
				rr.versionlessUrlResourcesMap[key] = append(rr.versionlessUrlResourcesMap[key], entryResource)
			}
		}

		isRestful := restfulURLRegex.MatchString(fullURL)
		if isRestful {
			litInfo, err := reference.LiteralInfoFromURI(fullURL)
			if err != nil {
				continue
			}

			rootURL := litInfo.ServiceBaseURL()
			if rootURL == "" {
				continue
			}
			rr.rootURLs = append(rr.rootURLs, rootURL)
		}
	}
}

func (rr *resolvableResources) isEmpty() bool {
	return len(rr.versionlessUrlResourcesMap) == 0 && len(rr.versionedURLResourcesMap) == 0 && len(rr.rootURLs) == 0 && len(rr.urnResourceMap) == 0
}

// Regex Sources -
// - https://hl7.org/fhir/datatypes.html#primitive
// - https://hl7.org/fhir/r4/references.html#literal
var (
	oidRegex        = regexp.MustCompile(`^urn:oid:[0-2](\.(0|[1-9][0-9]*))+$`)
	uuidRegex       = regexp.MustCompile(`^urn:uuid:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	restfulURLRegex = regexp.MustCompile(`^((http|https)://([A-Za-z0-9\-\.:%\\$]*/)+)?(Account|ActivityDefinition|ActorDefinition|AdministrableProductDefinition|AdverseEvent|AllergyIntolerance|Appointment|AppointmentResponse|ArtifactAssessment|AuditEvent|Basic|Binary|BiologicallyDerivedProduct|BiologicallyDerivedProductDispense|BodyStructure|Bundle|CapabilityStatement|CarePlan|CareTeam|ChargeItem|ChargeItemDefinition|Citation|Claim|ClaimResponse|ClinicalAssessment|ClinicalUseDefinition|CodeSystem|Communication|CommunicationRequest|CompartmentDefinition|Composition|ConceptMap|Condition|ConditionDefinition|Consent|Contract|Coverage|CoverageEligibilityRequest|CoverageEligibilityResponse|DetectedIssue|Device|DeviceAlert|DeviceAssociation|DeviceDefinition|DeviceDispense|DeviceMetric|DeviceRequest|DeviceUsage|DiagnosticReport|DocumentReference|Encounter|EncounterHistory|Endpoint|EnrollmentRequest|EnrollmentResponse|EpisodeOfCare|EventDefinition|Evidence|EvidenceVariable|ExampleScenario|ExplanationOfBenefit|FamilyMemberHistory|Flag|FormularyItem|GenomicStudy|Goal|GraphDefinition|Group|GuidanceResponse|HealthcareService|ImagingSelection|ImagingStudy|Immunization|ImmunizationEvaluation|ImmunizationRecommendation|ImplementationGuide|Ingredient|InsurancePlan|InsuranceProduct|InventoryItem|InventoryReport|Invoice|Library|Linkage|List|Location|ManufacturedItemDefinition|Measure|MeasureReport|Medication|MedicationAdministration|MedicationDispense|MedicationKnowledge|MedicationRequest|MedicationStatement|MedicinalProductDefinition|MessageDefinition|MessageHeader|MolecularDefinition|MolecularSequence|NamingSystem|NutritionIntake|NutritionOrder|NutritionProduct|Observation|ObservationDefinition|OperationDefinition|OperationOutcome|Organization|OrganizationAffiliation|PackagedProductDefinition|Patient|PaymentNotice|PaymentReconciliation|Permission|Person|PersonalRelationship|PlanDefinition|Practitioner|PractitionerRole|Procedure|Provenance|Questionnaire|QuestionnaireResponse|RegulatedAuthorization|RelatedPerson|RequestOrchestration|Requirements|ResearchStudy|ResearchSubject|RiskAssessment|Schedule|SearchParameter|ServiceRequest|Slot|Specimen|SpecimenDefinition|StructureDefinition|StructureMap|Subscription|SubscriptionStatus|SubscriptionTopic|Substance|SubstanceDefinition|SubstanceNucleicAcid|SubstancePolymer|SubstanceProtein|SubstanceReferenceInformation|SubstanceSourceMaterial|SupplyDelivery|SupplyRequest|Task|TerminologyCapabilities|TestPlan|TestReport|TestScript|Transport|ValueSet|VerificationResult|VisionPrescription)/[A-Za-z0-9\-\.]{1,64}(/_history/[A-Za-z0-9\-\.]{1,64})?(#[A-Za-z0-9\-\.]{1,64})?$`)
)

// isURN checks if a reference value is a URN (e.g., "urn:uuid:..." or "urn:oid:...").
func isURN(refValue string) bool {
	return oidRegex.MatchString(refValue) || uuidRegex.MatchString(refValue)
}

// absoluteURLInfo is a helper function that returns base URL, identity string, isVersioned boolean, and a boolean indicating if the given reference is an absolute URL, for a given reference string
func absoluteURLInfo(refValue string) (baseURL string, identityStr string, isVersioned bool, isAbsoluteURL bool) {
	identity, err := reference.IdentityFromAbsoluteURL(refValue)
	if err != nil {
		return "", "", false, false
	}
	litInfo, err := reference.LiteralInfoFromURI(refValue)
	if err != nil {
		return "", "", false, false
	}
	baseURL = litInfo.ServiceBaseURL()
	_, isVersioned = identity.VersionID()
	return baseURL, identity.String(), isVersioned, true
}

func isRelativeURL(refValue string) bool {
	_, err := reference.IdentityFromRelativeURI(refValue)
	return err == nil
}

// toAbsolute is a helper function that returns an absolute URL, given a base URL, and a relative identity
func toAbsolute(baseURL string, identity string) string {
	return fmt.Sprintf("%s/%s", baseURL, identity)
}

func (rr *resolvableResources) resolveURN(ref string) (fhir.Resource, error) {
	resource, ok := rr.urnResourceMap[ref]
	if !ok {
		return nil, nil
	}
	return resource, nil
}

// getLatestResource returns the resource with the latest meta.LastUpdated field, from a slice of Resources.
// If any of the given resources has a missing meta.LastUpdated field, an Undefined behavior error is thrown.
func getLatestResource(resources []fhir.Resource) (fhir.Resource, error) {
	if resources[0].GetMeta().GetLastUpdated() == nil {
		return nil, ErrMissingMetaOrLastUpdated
	}
	resolvedResource := resources[0]
	resolvedLastUpdated := resolvedResource.GetMeta().GetLastUpdated().GetValueUs()

	for _, res := range resources[1:] {
		if res.GetMeta().GetLastUpdated() == nil {
			return nil, ErrMissingMetaOrLastUpdated
		}
		resLastUpdated := res.GetMeta().GetLastUpdated().GetValueUs()
		if resolvedLastUpdated > resLastUpdated {
			resolvedResource = res
			resolvedLastUpdated = resLastUpdated
		}
	}
	return resolvedResource, nil
}

func (rr *resolvableResources) resolveAbsolute(baseURL string, identityStr string, isVersioned bool) (fhir.Resource, error) {
	key := toAbsolute(baseURL, identityStr)
	if isVersioned {
		// Versioned absolute URL
		resources, found := rr.versionedURLResourcesMap[key]
		if !found {
			return nil, nil
		}
		if len(resources) > 1 {
			return nil, ErrMultipleResourcesWithSameIDAndVersion
		}
		return resources[0], nil
	}
	// Versionless absolute URL
	resources, found := rr.versionlessUrlResourcesMap[key]
	if !found {
		return nil, nil
	}

	if len(resources) == 1 {
		return resources[0], nil
	}
	// If more than one versionless matches are found, return the one with the latest meta.LastUpdated field.
	// If any of the matched resources has a missing meta.LastUpdated field, an Undefined behavior error is thrown.
	return getLatestResource(resources)
}

func (rr *resolvableResources) resolveRelative(ref string) (fhir.Resource, error) {
	if len(rr.rootURLs) == 0 {
		return nil, nil
	}
	for _, rootURL := range rr.rootURLs {
		constructedRef := toAbsolute(rootURL, ref)
		baseURL, identityStr, isVersioned, isAbsoluteURL := absoluteURLInfo(constructedRef)
		if !isAbsoluteURL {
			continue
		}
		resolvedRes, err := rr.resolveAbsolute(baseURL, identityStr, isVersioned)
		if resolvedRes != nil && err == nil {
			return resolvedRes, nil
		}
	}
	return nil, nil
}

func (rr *resolvableResources) resolveReference(ref string) (fhir.Resource, error) {
	if isURN(ref) {
		return rr.resolveURN(ref)
	}
	baseURL, identityStr, isVersioned, isAbsoluteURL := absoluteURLInfo(ref)
	if isAbsoluteURL {
		return rr.resolveAbsolute(baseURL, identityStr, isVersioned)
	}
	if isRelativeURL(ref) {
		return rr.resolveRelative(ref)
	}
	return nil, nil
}

// Resolve implements the FHIRPath resolve() function in the context of a bundle of type SEARCHSET or COLLECTION.
// It identifies and returns FHIR resources from the resolver's bundle that match the provided list of reference strings.
// The BundleResolver strictly supports references of type URN and URL (both absolute and relative), but not canonical references.
// It returns an empty slice if the bundle is empty, or if no matching resources are found for the given references.
// Invalid references and invalid bundle resource entries are ignored.
// Source for the algorithm implementation -
// https://build.fhir.org/bundle.html#references
func (br *BundleResolver) Resolve(toResolveRefs []string) ([]fhir.Resource, error) {
	if br.bundle == nil || len(br.bundle.Entry) == 0 || len(toResolveRefs) == 0 {
		return []fhir.Resource{}, nil
	}

	rr := &resolvableResources{
		urnResourceMap:             map[string]fhir.Resource{},
		versionlessUrlResourcesMap: map[string][]fhir.Resource{},
		versionedURLResourcesMap:   map[string][]fhir.Resource{},
		rootURLs:                   []string{},
	}
	rr.fromBundle(br.bundle)

	if rr.isEmpty() {
		return []fhir.Resource{}, nil
	}

	resolvedResources := []fhir.Resource{}
	for _, ref := range toResolveRefs {
		resolvedResource, err := rr.resolveReference(ref)
		if err != nil {
			return nil, err
		}
		if resolvedResource != nil {
			resolvedResources = append(resolvedResources, resolvedResource)
		}
	}
	return resolvedResources, nil
}
