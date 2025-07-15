package impl

import (
	"errors"
	"fmt"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
)

var (
	ErrInvalidReference     = errors.New("invalid reference")
	ErrUnconfiguredResolver = errors.New("resolve() function requires a Resolver to be configured in the evaluation context")
)

func isValidResolveInput(item any) bool {
	switch item.(type) {
	case *dtpb.String, *dtpb.Uri, *dtpb.Url, *dtpb.Canonical, *dtpb.Reference, system.String, string:
		return true
	default:
		return false
	}
}

func stringifyReference(ref *dtpb.Reference) string {

	if ref == nil || ref.GetType() == nil || ref.GetReference() == nil {
		return ""
	}
	refType := ref.GetType().GetValue()

	switch ref.GetReference().(type) {
	case *dtpb.Reference_AccountId:
		return fmt.Sprintf("%s/%s", refType, ref.GetAccountId().GetValue())
	case *dtpb.Reference_ActivityDefinitionId:
		return fmt.Sprintf("%s/%s", refType, ref.GetActivityDefinitionId().GetValue())
	case *dtpb.Reference_AdverseEventId:
		return fmt.Sprintf("%s/%s", refType, ref.GetAdverseEventId().GetValue())
	case *dtpb.Reference_AllergyIntoleranceId:
		return fmt.Sprintf("%s/%s", refType, ref.GetAllergyIntoleranceId().GetValue())
	case *dtpb.Reference_AppointmentId:
		return fmt.Sprintf("%s/%s", refType, ref.GetAppointmentId().GetValue())
	case *dtpb.Reference_AppointmentResponseId:
		return fmt.Sprintf("%s/%s", refType, ref.GetAppointmentResponseId().GetValue())
	case *dtpb.Reference_AuditEventId:
		return fmt.Sprintf("%s/%s", refType, ref.GetAuditEventId().GetValue())
	case *dtpb.Reference_BasicId:
		return fmt.Sprintf("%s/%s", refType, ref.GetBasicId().GetValue())
	case *dtpb.Reference_BinaryId:
		return fmt.Sprintf("%s/%s", refType, ref.GetBinaryId().GetValue())
	case *dtpb.Reference_BiologicallyDerivedProductId:
		return fmt.Sprintf("%s/%s", refType, ref.GetBiologicallyDerivedProductId().GetValue())
	case *dtpb.Reference_BodyStructureId:
		return fmt.Sprintf("%s/%s", refType, ref.GetBodyStructureId().GetValue())
	case *dtpb.Reference_BundleId:
		return fmt.Sprintf("%s/%s", refType, ref.GetBundleId().GetValue())
	case *dtpb.Reference_CapabilityStatementId:
		return fmt.Sprintf("%s/%s", refType, ref.GetCapabilityStatementId().GetValue())
	case *dtpb.Reference_CarePlanId:
		return fmt.Sprintf("%s/%s", refType, ref.GetCarePlanId().GetValue())
	case *dtpb.Reference_CareTeamId:
		return fmt.Sprintf("%s/%s", refType, ref.GetCareTeamId().GetValue())
	case *dtpb.Reference_CatalogEntryId:
		return fmt.Sprintf("%s/%s", refType, ref.GetCatalogEntryId().GetValue())
	case *dtpb.Reference_ChargeItemId:
		return fmt.Sprintf("%s/%s", refType, ref.GetChargeItemId().GetValue())
	case *dtpb.Reference_ChargeItemDefinitionId:
		return fmt.Sprintf("%s/%s", refType, ref.GetChargeItemDefinitionId().GetValue())
	case *dtpb.Reference_ClaimId:
		return fmt.Sprintf("%s/%s", refType, ref.GetClaimId().GetValue())
	case *dtpb.Reference_ClaimResponseId:
		return fmt.Sprintf("%s/%s", refType, ref.GetClaimResponseId().GetValue())
	case *dtpb.Reference_ClinicalImpressionId:
		return fmt.Sprintf("%s/%s", refType, ref.GetClinicalImpressionId().GetValue())
	case *dtpb.Reference_CodeSystemId:
		return fmt.Sprintf("%s/%s", refType, ref.GetCodeSystemId().GetValue())
	case *dtpb.Reference_CommunicationId:
		return fmt.Sprintf("%s/%s", refType, ref.GetCommunicationId().GetValue())
	case *dtpb.Reference_CommunicationRequestId:
		return fmt.Sprintf("%s/%s", refType, ref.GetCommunicationRequestId().GetValue())
	case *dtpb.Reference_CompartmentDefinitionId:
		return fmt.Sprintf("%s/%s", refType, ref.GetCompartmentDefinitionId().GetValue())
	case *dtpb.Reference_CompositionId:
		return fmt.Sprintf("%s/%s", refType, ref.GetCompositionId().GetValue())
	case *dtpb.Reference_ConceptMapId:
		return fmt.Sprintf("%s/%s", refType, ref.GetConceptMapId().GetValue())
	case *dtpb.Reference_ConditionId:
		return fmt.Sprintf("%s/%s", refType, ref.GetConditionId().GetValue())
	case *dtpb.Reference_ConsentId:
		return fmt.Sprintf("%s/%s", refType, ref.GetConsentId().GetValue())
	case *dtpb.Reference_ContractId:
		return fmt.Sprintf("%s/%s", refType, ref.GetContractId().GetValue())
	case *dtpb.Reference_CoverageId:
		return fmt.Sprintf("%s/%s", refType, ref.GetCoverageId().GetValue())
	case *dtpb.Reference_CoverageEligibilityRequestId:
		return fmt.Sprintf("%s/%s", refType, ref.GetCoverageEligibilityRequestId().GetValue())
	case *dtpb.Reference_CoverageEligibilityResponseId:
		return fmt.Sprintf("%s/%s", refType, ref.GetCoverageEligibilityResponseId().GetValue())
	case *dtpb.Reference_DetectedIssueId:
		return fmt.Sprintf("%s/%s", refType, ref.GetDetectedIssueId().GetValue())
	case *dtpb.Reference_DeviceDefinitionId:
		return fmt.Sprintf("%s/%s", refType, ref.GetDeviceDefinitionId().GetValue())
	case *dtpb.Reference_DeviceId:
		return fmt.Sprintf("%s/%s", refType, ref.GetDeviceId().GetValue())
	case *dtpb.Reference_DeviceMetricId:
		return fmt.Sprintf("%s/%s", refType, ref.GetDeviceMetricId().GetValue())
	case *dtpb.Reference_DeviceRequestId:
		return fmt.Sprintf("%s/%s", refType, ref.GetDeviceRequestId().GetValue())
	case *dtpb.Reference_DeviceUseStatementId:
		return fmt.Sprintf("%s/%s", refType, ref.GetDeviceUseStatementId().GetValue())
	case *dtpb.Reference_DiagnosticReportId:
		return fmt.Sprintf("%s/%s", refType, ref.GetDiagnosticReportId().GetValue())
	case *dtpb.Reference_DocumentManifestId:
		return fmt.Sprintf("%s/%s", refType, ref.GetDocumentManifestId().GetValue())
	case *dtpb.Reference_DocumentReferenceId:
		return fmt.Sprintf("%s/%s", refType, ref.GetDocumentReferenceId().GetValue())
	case *dtpb.Reference_DomainResourceId:
		return fmt.Sprintf("%s/%s", refType, ref.GetDomainResourceId().GetValue())
	case *dtpb.Reference_EffectEvidenceSynthesisId:
		return fmt.Sprintf("%s/%s", refType, ref.GetEffectEvidenceSynthesisId().GetValue())
	case *dtpb.Reference_EncounterId:
		return fmt.Sprintf("%s/%s", refType, ref.GetEncounterId().GetValue())
	case *dtpb.Reference_EndpointId:
		return fmt.Sprintf("%s/%s", refType, ref.GetEndpointId().GetValue())
	case *dtpb.Reference_EnrollmentRequestId:
		return fmt.Sprintf("%s/%s", refType, ref.GetEnrollmentRequestId().GetValue())
	case *dtpb.Reference_EnrollmentResponseId:
		return fmt.Sprintf("%s/%s", refType, ref.GetEnrollmentResponseId().GetValue())
	case *dtpb.Reference_EpisodeOfCareId:
		return fmt.Sprintf("%s/%s", refType, ref.GetEpisodeOfCareId().GetValue())
	case *dtpb.Reference_EventDefinitionId:
		return fmt.Sprintf("%s/%s", refType, ref.GetEventDefinitionId().GetValue())
	case *dtpb.Reference_EvidenceId:
		return fmt.Sprintf("%s/%s", refType, ref.GetEvidenceId().GetValue())
	case *dtpb.Reference_EvidenceVariableId:
		return fmt.Sprintf("%s/%s", refType, ref.GetEvidenceVariableId().GetValue())
	case *dtpb.Reference_ExampleScenarioId:
		return fmt.Sprintf("%s/%s", refType, ref.GetExampleScenarioId().GetValue())
	case *dtpb.Reference_ExplanationOfBenefitId:
		return fmt.Sprintf("%s/%s", refType, ref.GetExplanationOfBenefitId().GetValue())
	case *dtpb.Reference_FamilyMemberHistoryId:
		return fmt.Sprintf("%s/%s", refType, ref.GetFamilyMemberHistoryId().GetValue())
	case *dtpb.Reference_FlagId:
		return fmt.Sprintf("%s/%s", refType, ref.GetFlagId().GetValue())
	case *dtpb.Reference_Fragment:
		return fmt.Sprintf("%s/%s", refType, ref.GetFragment())
	case *dtpb.Reference_GoalId:
		return fmt.Sprintf("%s/%s", refType, ref.GetGoalId().GetValue())
	case *dtpb.Reference_GraphDefinitionId:
		return fmt.Sprintf("%s/%s", refType, ref.GetGraphDefinitionId().GetValue())
	case *dtpb.Reference_GroupId:
		return fmt.Sprintf("%s/%s", refType, ref.GetGroupId().GetValue())
	case *dtpb.Reference_GuidanceResponseId:
		return fmt.Sprintf("%s/%s", refType, ref.GetGuidanceResponseId().GetValue())
	case *dtpb.Reference_HealthcareServiceId:
		return fmt.Sprintf("%s/%s", refType, ref.GetHealthcareServiceId().GetValue())
	case *dtpb.Reference_ImagingStudyId:
		return fmt.Sprintf("%s/%s", refType, ref.GetImagingStudyId().GetValue())
	case *dtpb.Reference_ImmunizationEvaluationId:
		return fmt.Sprintf("%s/%s", refType, ref.GetImmunizationEvaluationId().GetValue())
	case *dtpb.Reference_ImmunizationId:
		return fmt.Sprintf("%s/%s", refType, ref.GetImmunizationId().GetValue())
	case *dtpb.Reference_ImmunizationRecommendationId:
		return fmt.Sprintf("%s/%s", refType, ref.GetImmunizationRecommendationId().GetValue())
	case *dtpb.Reference_ImplementationGuideId:
		return fmt.Sprintf("%s/%s", refType, ref.GetImplementationGuideId().GetValue())
	case *dtpb.Reference_InsurancePlanId:
		return fmt.Sprintf("%s/%s", refType, ref.GetInsurancePlanId().GetValue())
	case *dtpb.Reference_InvoiceId:
		return fmt.Sprintf("%s/%s", refType, ref.GetInvoiceId().GetValue())
	case *dtpb.Reference_LibraryId:
		return fmt.Sprintf("%s/%s", refType, ref.GetLibraryId().GetValue())
	case *dtpb.Reference_LinkageId:
		return fmt.Sprintf("%s/%s", refType, ref.GetLinkageId().GetValue())
	case *dtpb.Reference_ListId:
		return fmt.Sprintf("%s/%s", refType, ref.GetListId().GetValue())
	case *dtpb.Reference_LocationId:
		return fmt.Sprintf("%s/%s", refType, ref.GetLocationId().GetValue())
	case *dtpb.Reference_MeasureId:
		return fmt.Sprintf("%s/%s", refType, ref.GetMeasureId().GetValue())
	case *dtpb.Reference_MeasureReportId:
		return fmt.Sprintf("%s/%s", refType, ref.GetMeasureReportId().GetValue())
	case *dtpb.Reference_MediaId:
		return fmt.Sprintf("%s/%s", refType, ref.GetMediaId().GetValue())
	case *dtpb.Reference_MedicationAdministrationId:
		return fmt.Sprintf("%s/%s", refType, ref.GetMedicationAdministrationId().GetValue())
	case *dtpb.Reference_MedicationDispenseId:
		return fmt.Sprintf("%s/%s", refType, ref.GetMedicationDispenseId().GetValue())
	case *dtpb.Reference_MedicationId:
		return fmt.Sprintf("%s/%s", refType, ref.GetMedicationId().GetValue())
	case *dtpb.Reference_MedicationKnowledgeId:
		return fmt.Sprintf("%s/%s", refType, ref.GetMedicationKnowledgeId().GetValue())
	case *dtpb.Reference_MedicationRequestId:
		return fmt.Sprintf("%s/%s", refType, ref.GetMedicationRequestId().GetValue())
	case *dtpb.Reference_MedicationStatementId:
		return fmt.Sprintf("%s/%s", refType, ref.GetMedicationStatementId().GetValue())
	case *dtpb.Reference_MedicinalProductAuthorizationId:
		return fmt.Sprintf("%s/%s", refType, ref.GetMedicinalProductAuthorizationId().GetValue())
	case *dtpb.Reference_MedicinalProductContraindicationId:
		return fmt.Sprintf("%s/%s", refType, ref.GetMedicinalProductContraindicationId().GetValue())
	case *dtpb.Reference_MedicinalProductId:
		return fmt.Sprintf("%s/%s", refType, ref.GetMedicinalProductId().GetValue())
	case *dtpb.Reference_MedicinalProductIndicationId:
		return fmt.Sprintf("%s/%s", refType, ref.GetMedicinalProductIndicationId().GetValue())
	case *dtpb.Reference_MedicinalProductIngredientId:
		return fmt.Sprintf("%s/%s", refType, ref.GetMedicinalProductIngredientId().GetValue())
	case *dtpb.Reference_MedicinalProductInteractionId:
		return fmt.Sprintf("%s/%s", refType, ref.GetMedicinalProductInteractionId().GetValue())
	case *dtpb.Reference_MedicinalProductManufacturedId:
		return fmt.Sprintf("%s/%s", refType, ref.GetMedicinalProductManufacturedId().GetValue())
	case *dtpb.Reference_MedicinalProductPackagedId:
		return fmt.Sprintf("%s/%s", refType, ref.GetMedicinalProductPackagedId().GetValue())
	case *dtpb.Reference_MedicinalProductPharmaceuticalId:
		return fmt.Sprintf("%s/%s", refType, ref.GetMedicinalProductPharmaceuticalId().GetValue())
	case *dtpb.Reference_MedicinalProductUndesirableEffectId:
		return fmt.Sprintf("%s/%s", refType, ref.GetMedicinalProductUndesirableEffectId().GetValue())
	case *dtpb.Reference_MessageDefinitionId:
		return fmt.Sprintf("%s/%s", refType, ref.GetMessageDefinitionId().GetValue())
	case *dtpb.Reference_MessageHeaderId:
		return fmt.Sprintf("%s/%s", refType, ref.GetMessageHeaderId().GetValue())
	case *dtpb.Reference_MetadataResourceId:
		return fmt.Sprintf("%s/%s", refType, ref.GetMetadataResourceId().GetValue())
	case *dtpb.Reference_MolecularSequenceId:
		return fmt.Sprintf("%s/%s", refType, ref.GetMolecularSequenceId().GetValue())
	case *dtpb.Reference_NamingSystemId:
		return fmt.Sprintf("%s/%s", refType, ref.GetNamingSystemId().GetValue())
	case *dtpb.Reference_NutritionOrderId:
		return fmt.Sprintf("%s/%s", refType, ref.GetNutritionOrderId().GetValue())
	case *dtpb.Reference_ObservationDefinitionId:
		return fmt.Sprintf("%s/%s", refType, ref.GetObservationDefinitionId().GetValue())
	case *dtpb.Reference_ObservationId:
		return fmt.Sprintf("%s/%s", refType, ref.GetObservationId().GetValue())
	case *dtpb.Reference_OperationDefinitionId:
		return fmt.Sprintf("%s/%s", refType, ref.GetOperationDefinitionId().GetValue())
	case *dtpb.Reference_OperationOutcomeId:
		return fmt.Sprintf("%s/%s", refType, ref.GetOperationOutcomeId().GetValue())
	case *dtpb.Reference_OrganizationAffiliationId:
		return fmt.Sprintf("%s/%s", refType, ref.GetOrganizationAffiliationId().GetValue())
	case *dtpb.Reference_OrganizationId:
		return fmt.Sprintf("%s/%s", refType, ref.GetOrganizationId().GetValue())
	case *dtpb.Reference_ParametersId:
		return fmt.Sprintf("%s/%s", refType, ref.GetParametersId().GetValue())
	case *dtpb.Reference_PatientId:
		return fmt.Sprintf("%s/%s", refType, ref.GetPatientId().GetValue())
	case *dtpb.Reference_PaymentNoticeId:
		return fmt.Sprintf("%s/%s", refType, ref.GetPaymentNoticeId().GetValue())
	case *dtpb.Reference_PaymentReconciliationId:
		return fmt.Sprintf("%s/%s", refType, ref.GetPaymentReconciliationId().GetValue())
	case *dtpb.Reference_PersonId:
		return fmt.Sprintf("%s/%s", refType, ref.GetPersonId().GetValue())
	case *dtpb.Reference_PlanDefinitionId:
		return fmt.Sprintf("%s/%s", refType, ref.GetPlanDefinitionId().GetValue())
	case *dtpb.Reference_PractitionerId:
		return fmt.Sprintf("%s/%s", refType, ref.GetPractitionerId().GetValue())
	case *dtpb.Reference_PractitionerRoleId:
		return fmt.Sprintf("%s/%s", refType, ref.GetPractitionerRoleId().GetValue())
	case *dtpb.Reference_ProcedureId:
		return fmt.Sprintf("%s/%s", refType, ref.GetProcedureId().GetValue())
	case *dtpb.Reference_ProvenanceId:
		return fmt.Sprintf("%s/%s", refType, ref.GetProvenanceId().GetValue())
	case *dtpb.Reference_QuestionnaireId:
		return fmt.Sprintf("%s/%s", refType, ref.GetQuestionnaireId().GetValue())
	case *dtpb.Reference_QuestionnaireResponseId:
		return fmt.Sprintf("%s/%s", refType, ref.GetQuestionnaireResponseId().GetValue())
	case *dtpb.Reference_RelatedPersonId:
		return fmt.Sprintf("%s/%s", refType, ref.GetRelatedPersonId().GetValue())
	case *dtpb.Reference_RequestGroupId:
		return fmt.Sprintf("%s/%s", refType, ref.GetRequestGroupId().GetValue())
	case *dtpb.Reference_ResearchDefinitionId:
		return fmt.Sprintf("%s/%s", refType, ref.GetResearchDefinitionId().GetValue())
	case *dtpb.Reference_ResearchElementDefinitionId:
		return fmt.Sprintf("%s/%s", refType, ref.GetResearchElementDefinitionId().GetValue())
	case *dtpb.Reference_ResearchStudyId:
		return fmt.Sprintf("%s/%s", refType, ref.GetResearchStudyId().GetValue())
	case *dtpb.Reference_ResearchSubjectId:
		return fmt.Sprintf("%s/%s", refType, ref.GetResearchSubjectId().GetValue())
	case *dtpb.Reference_ResourceId:
		return fmt.Sprintf("%s/%s", refType, ref.GetResourceId())
	case *dtpb.Reference_RiskAssessmentId:
		return fmt.Sprintf("%s/%s", refType, ref.GetRiskAssessmentId().GetValue())
	case *dtpb.Reference_RiskEvidenceSynthesisId:
		return fmt.Sprintf("%s/%s", refType, ref.GetRiskEvidenceSynthesisId().GetValue())
	case *dtpb.Reference_ScheduleId:
		return fmt.Sprintf("%s/%s", refType, ref.GetScheduleId().GetValue())
	case *dtpb.Reference_SearchParameterId:
		return fmt.Sprintf("%s/%s", refType, ref.GetSearchParameterId().GetValue())
	case *dtpb.Reference_ServiceRequestId:
		return fmt.Sprintf("%s/%s", refType, ref.GetServiceRequestId().GetValue())
	case *dtpb.Reference_SlotId:
		return fmt.Sprintf("%s/%s", refType, ref.GetSlotId().GetValue())
	case *dtpb.Reference_SpecimenDefinitionId:
		return fmt.Sprintf("%s/%s", refType, ref.GetSpecimenDefinitionId().GetValue())
	case *dtpb.Reference_SpecimenId:
		return fmt.Sprintf("%s/%s", refType, ref.GetSpecimenId().GetValue())
	case *dtpb.Reference_StructureDefinitionId:
		return fmt.Sprintf("%s/%s", refType, ref.GetStructureDefinitionId().GetValue())
	case *dtpb.Reference_StructureMapId:
		return fmt.Sprintf("%s/%s", refType, ref.GetStructureMapId().GetValue())
	case *dtpb.Reference_SubscriptionId:
		return fmt.Sprintf("%s/%s", refType, ref.GetSubscriptionId().GetValue())
	case *dtpb.Reference_SubstanceId:
		return fmt.Sprintf("%s/%s", refType, ref.GetSubstanceId().GetValue())
	case *dtpb.Reference_SubstanceNucleicAcidId:
		return fmt.Sprintf("%s/%s", refType, ref.GetSubstanceNucleicAcidId().GetValue())
	case *dtpb.Reference_SubstancePolymerId:
		return fmt.Sprintf("%s/%s", refType, ref.GetSubstancePolymerId().GetValue())
	case *dtpb.Reference_SubstanceProteinId:
		return fmt.Sprintf("%s/%s", refType, ref.GetSubstanceProteinId().GetValue())
	case *dtpb.Reference_SubstanceReferenceInformationId:
		return fmt.Sprintf("%s/%s", refType, ref.GetSubstanceReferenceInformationId().GetValue())
	case *dtpb.Reference_SubstanceSourceMaterialId:
		return fmt.Sprintf("%s/%s", refType, ref.GetSubstanceSourceMaterialId().GetValue())
	case *dtpb.Reference_SubstanceSpecificationId:
		return fmt.Sprintf("%s/%s", refType, ref.GetSubstanceSpecificationId().GetValue())
	case *dtpb.Reference_SupplyDeliveryId:
		return fmt.Sprintf("%s/%s", refType, ref.GetSupplyDeliveryId().GetValue())
	case *dtpb.Reference_SupplyRequestId:
		return fmt.Sprintf("%s/%s", refType, ref.GetSupplyRequestId().GetValue())
	case *dtpb.Reference_TaskId:
		return fmt.Sprintf("%s/%s", refType, ref.GetTaskId().GetValue())
	case *dtpb.Reference_TerminologyCapabilitiesId:
		return fmt.Sprintf("%s/%s", refType, ref.GetTerminologyCapabilitiesId().GetValue())
	case *dtpb.Reference_TestReportId:
		return fmt.Sprintf("%s/%s", refType, ref.GetTestReportId().GetValue())
	case *dtpb.Reference_TestScriptId:
		return fmt.Sprintf("%s/%s", refType, ref.GetTestScriptId().GetValue())
	case *dtpb.Reference_Uri:
		return fmt.Sprintf("%s/%s", refType, ref.GetUri())
	case *dtpb.Reference_ValueSetId:
		return fmt.Sprintf("%s/%s", refType, ref.GetValueSetId().GetValue())
	case *dtpb.Reference_VerificationResultId:
		return fmt.Sprintf("%s/%s", refType, ref.GetVerificationResultId().GetValue())
	case *dtpb.Reference_VisionPrescriptionId:
		return fmt.Sprintf("%s/%s", refType, ref.GetVisionPrescriptionId().GetValue())
	default:
		return ""
	}
}

func toString(item any) string {
	switch item := item.(type) {
	case *dtpb.String:
		return item.GetValue()
	case *dtpb.Uri:
		return item.GetValue()
	case *dtpb.Url:
		return item.GetValue()
	case *dtpb.Canonical:
		return item.GetValue()
	case *dtpb.Reference:
		return stringifyReference(item)
	case system.String:
		return string(item)
	default:
		return ""
	}
}

// Resolve locates the target of each reference in the input collection. For each item that
// is a string representing a URI (or canonical or URL), or a Reference, Resolve locates
// the target resource and adds it to the resulting collection. If an item does not resolve
// to a resource, it is ignored. If the input is empty, the output will be empty.
func Resolve(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	if argLen := len(args); argLen != 0 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 0", ErrWrongArity, argLen)
	}

	toResolve := []string{}
	for _, item := range input {
		if isValidResolveInput(item) {
			toResolve = append(toResolve, toString(item))
		}
	}

	if len(toResolve) == 0 {
		return system.Collection{}, nil
	}

	resolverImpl := ctx.Resolver
	if resolverImpl == nil {
		return nil, ErrUnconfiguredResolver
	}
	resources, err := resolverImpl.Resolve(toResolve)
	if err != nil {
		return nil, err
	}

	resolved := system.Collection{}
	for _, res := range resources {
		resolved = append(resolved, res)
	}
	return resolved, nil
}
