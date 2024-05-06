package fhir_test

import (
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/account_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/activity_definition_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/adverse_event_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/allergy_intolerance_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/appointment_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/appointment_response_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/audit_event_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/basic_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/binary_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/biologically_derived_product_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/body_structure_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/capability_statement_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/care_plan_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/care_team_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/catalog_entry_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/charge_item_definition_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/charge_item_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/claim_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/claim_response_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/clinical_impression_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/code_system_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/communication_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/communication_request_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/compartment_definition_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/composition_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/concept_map_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/condition_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/consent_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/contract_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/coverage_eligibility_request_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/coverage_eligibility_response_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/coverage_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/detected_issue_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/device_definition_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/device_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/device_metric_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/device_request_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/device_use_statement_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/diagnostic_report_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/document_manifest_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/document_reference_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/effect_evidence_synthesis_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/encounter_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/endpoint_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/enrollment_request_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/enrollment_response_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/episode_of_care_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/event_definition_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/evidence_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/evidence_variable_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/example_scenario_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/explanation_of_benefit_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/family_member_history_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/flag_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/goal_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/graph_definition_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/group_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/guidance_response_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/healthcare_service_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/imaging_study_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/immunization_evaluation_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/immunization_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/immunization_recommendation_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/implementation_guide_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/insurance_plan_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/invoice_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/library_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/linkage_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/list_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/location_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/measure_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/measure_report_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/media_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/medication_administration_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/medication_dispense_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/medication_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/medication_knowledge_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/medication_request_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/medication_statement_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/medicinal_product_authorization_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/medicinal_product_contraindication_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/medicinal_product_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/medicinal_product_indication_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/medicinal_product_ingredient_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/medicinal_product_interaction_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/medicinal_product_manufactured_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/medicinal_product_packaged_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/medicinal_product_pharmaceutical_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/medicinal_product_undesirable_effect_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/message_definition_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/message_header_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/molecular_sequence_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/naming_system_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/nutrition_order_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/observation_definition_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/observation_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/operation_definition_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/operation_outcome_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/organization_affiliation_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/organization_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/parameters_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/payment_notice_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/payment_reconciliation_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/person_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/plan_definition_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/practitioner_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/practitioner_role_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/procedure_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/provenance_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/questionnaire_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/questionnaire_response_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/related_person_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/request_group_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/research_definition_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/research_element_definition_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/research_study_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/research_subject_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/risk_assessment_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/risk_evidence_synthesis_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/schedule_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/search_parameter_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/service_request_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/slot_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/specimen_definition_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/specimen_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/structure_definition_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/structure_map_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/subscription_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/substance_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/substance_nucleic_acid_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/substance_polymer_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/substance_protein_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/substance_reference_information_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/substance_source_material_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/substance_specification_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/supply_delivery_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/supply_request_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/task_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/terminology_capabilities_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/test_report_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/test_script_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/value_set_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/verification_result_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/vision_prescription_go_proto"
	"github.com/verily-src/fhirpath-go/internal/fhir"
)

// Note: The tests in this file are all statically validated by the interface
// type assigned to the anonymous variables.

// Resource types (https://www.hl7.org/fhir/resource.html#Resource).
// This may exclude certain Trial-Use types that google/fhir doesn't implement.

var _ fhir.Resource = (*account_go_proto.Account)(nil)
var _ fhir.Resource = (*activity_definition_go_proto.ActivityDefinition)(nil)
var _ fhir.Resource = (*adverse_event_go_proto.AdverseEvent)(nil)
var _ fhir.Resource = (*allergy_intolerance_go_proto.AllergyIntolerance)(nil)
var _ fhir.Resource = (*appointment_go_proto.Appointment)(nil)
var _ fhir.Resource = (*appointment_response_go_proto.AppointmentResponse)(nil)
var _ fhir.Resource = (*audit_event_go_proto.AuditEvent)(nil)
var _ fhir.Resource = (*basic_go_proto.Basic)(nil)
var _ fhir.Resource = (*biologically_derived_product_go_proto.BiologicallyDerivedProduct)(nil)
var _ fhir.Resource = (*body_structure_go_proto.BodyStructure)(nil)
var _ fhir.Resource = (*capability_statement_go_proto.CapabilityStatement)(nil)
var _ fhir.Resource = (*care_plan_go_proto.CarePlan)(nil)
var _ fhir.Resource = (*care_team_go_proto.CareTeam)(nil)
var _ fhir.Resource = (*catalog_entry_go_proto.CatalogEntry)(nil)
var _ fhir.Resource = (*charge_item_go_proto.ChargeItem)(nil)
var _ fhir.Resource = (*charge_item_definition_go_proto.ChargeItemDefinition)(nil)
var _ fhir.Resource = (*claim_go_proto.Claim)(nil)
var _ fhir.Resource = (*claim_response_go_proto.ClaimResponse)(nil)
var _ fhir.Resource = (*clinical_impression_go_proto.ClinicalImpression)(nil)
var _ fhir.Resource = (*code_system_go_proto.CodeSystem)(nil)
var _ fhir.Resource = (*communication_go_proto.Communication)(nil)
var _ fhir.Resource = (*communication_request_go_proto.CommunicationRequest)(nil)
var _ fhir.Resource = (*compartment_definition_go_proto.CompartmentDefinition)(nil)
var _ fhir.Resource = (*composition_go_proto.Composition)(nil)
var _ fhir.Resource = (*concept_map_go_proto.ConceptMap)(nil)
var _ fhir.Resource = (*condition_go_proto.Condition)(nil)
var _ fhir.Resource = (*consent_go_proto.Consent)(nil)
var _ fhir.Resource = (*contract_go_proto.Contract)(nil)
var _ fhir.Resource = (*coverage_go_proto.Coverage)(nil)
var _ fhir.Resource = (*coverage_eligibility_request_go_proto.CoverageEligibilityRequest)(nil)
var _ fhir.Resource = (*coverage_eligibility_response_go_proto.CoverageEligibilityResponse)(nil)
var _ fhir.Resource = (*detected_issue_go_proto.DetectedIssue)(nil)
var _ fhir.Resource = (*device_go_proto.Device)(nil)
var _ fhir.Resource = (*device_definition_go_proto.DeviceDefinition)(nil)
var _ fhir.Resource = (*device_metric_go_proto.DeviceMetric)(nil)
var _ fhir.Resource = (*device_request_go_proto.DeviceRequest)(nil)
var _ fhir.Resource = (*device_use_statement_go_proto.DeviceUseStatement)(nil)
var _ fhir.Resource = (*diagnostic_report_go_proto.DiagnosticReport)(nil)
var _ fhir.Resource = (*document_manifest_go_proto.DocumentManifest)(nil)
var _ fhir.Resource = (*document_reference_go_proto.DocumentReference)(nil)
var _ fhir.Resource = (*effect_evidence_synthesis_go_proto.EffectEvidenceSynthesis)(nil)
var _ fhir.Resource = (*encounter_go_proto.Encounter)(nil)
var _ fhir.Resource = (*endpoint_go_proto.Endpoint)(nil)
var _ fhir.Resource = (*enrollment_request_go_proto.EnrollmentRequest)(nil)
var _ fhir.Resource = (*enrollment_response_go_proto.EnrollmentResponse)(nil)
var _ fhir.Resource = (*episode_of_care_go_proto.EpisodeOfCare)(nil)
var _ fhir.Resource = (*event_definition_go_proto.EventDefinition)(nil)
var _ fhir.Resource = (*evidence_go_proto.Evidence)(nil)
var _ fhir.Resource = (*evidence_variable_go_proto.EvidenceVariable)(nil)
var _ fhir.Resource = (*example_scenario_go_proto.ExampleScenario)(nil)
var _ fhir.Resource = (*explanation_of_benefit_go_proto.ExplanationOfBenefit)(nil)
var _ fhir.Resource = (*family_member_history_go_proto.FamilyMemberHistory)(nil)
var _ fhir.Resource = (*flag_go_proto.Flag)(nil)
var _ fhir.Resource = (*goal_go_proto.Goal)(nil)
var _ fhir.Resource = (*graph_definition_go_proto.GraphDefinition)(nil)
var _ fhir.Resource = (*group_go_proto.Group)(nil)
var _ fhir.Resource = (*guidance_response_go_proto.GuidanceResponse)(nil)
var _ fhir.Resource = (*healthcare_service_go_proto.HealthcareService)(nil)
var _ fhir.Resource = (*imaging_study_go_proto.ImagingStudy)(nil)
var _ fhir.Resource = (*immunization_go_proto.Immunization)(nil)
var _ fhir.Resource = (*immunization_evaluation_go_proto.ImmunizationEvaluation)(nil)
var _ fhir.Resource = (*immunization_recommendation_go_proto.ImmunizationRecommendation)(nil)
var _ fhir.Resource = (*implementation_guide_go_proto.ImplementationGuide)(nil)
var _ fhir.Resource = (*insurance_plan_go_proto.InsurancePlan)(nil)
var _ fhir.Resource = (*invoice_go_proto.Invoice)(nil)
var _ fhir.Resource = (*library_go_proto.Library)(nil)
var _ fhir.Resource = (*linkage_go_proto.Linkage)(nil)
var _ fhir.Resource = (*list_go_proto.List)(nil)
var _ fhir.Resource = (*location_go_proto.Location)(nil)
var _ fhir.Resource = (*measure_go_proto.Measure)(nil)
var _ fhir.Resource = (*measure_report_go_proto.MeasureReport)(nil)
var _ fhir.Resource = (*media_go_proto.Media)(nil)
var _ fhir.Resource = (*medication_go_proto.Medication)(nil)
var _ fhir.Resource = (*medication_administration_go_proto.MedicationAdministration)(nil)
var _ fhir.Resource = (*medication_dispense_go_proto.MedicationDispense)(nil)
var _ fhir.Resource = (*medication_knowledge_go_proto.MedicationKnowledge)(nil)
var _ fhir.Resource = (*medication_request_go_proto.MedicationRequest)(nil)
var _ fhir.Resource = (*medication_statement_go_proto.MedicationStatement)(nil)
var _ fhir.Resource = (*medicinal_product_go_proto.MedicinalProduct)(nil)
var _ fhir.Resource = (*medicinal_product_authorization_go_proto.MedicinalProductAuthorization)(nil)
var _ fhir.Resource = (*medicinal_product_contraindication_go_proto.MedicinalProductContraindication)(nil)
var _ fhir.Resource = (*medicinal_product_indication_go_proto.MedicinalProductIndication)(nil)
var _ fhir.Resource = (*medicinal_product_ingredient_go_proto.MedicinalProductIngredient)(nil)
var _ fhir.Resource = (*medicinal_product_interaction_go_proto.MedicinalProductInteraction)(nil)
var _ fhir.Resource = (*medicinal_product_manufactured_go_proto.MedicinalProductManufactured)(nil)
var _ fhir.Resource = (*medicinal_product_packaged_go_proto.MedicinalProductPackaged)(nil)
var _ fhir.Resource = (*medicinal_product_pharmaceutical_go_proto.MedicinalProductPharmaceutical)(nil)
var _ fhir.Resource = (*medicinal_product_undesirable_effect_go_proto.MedicinalProductUndesirableEffect)(nil)
var _ fhir.Resource = (*message_definition_go_proto.MessageDefinition)(nil)
var _ fhir.Resource = (*message_header_go_proto.MessageHeader)(nil)
var _ fhir.Resource = (*molecular_sequence_go_proto.MolecularSequence)(nil)
var _ fhir.Resource = (*naming_system_go_proto.NamingSystem)(nil)
var _ fhir.Resource = (*nutrition_order_go_proto.NutritionOrder)(nil)
var _ fhir.Resource = (*observation_go_proto.Observation)(nil)
var _ fhir.Resource = (*observation_definition_go_proto.ObservationDefinition)(nil)
var _ fhir.Resource = (*operation_definition_go_proto.OperationDefinition)(nil)
var _ fhir.Resource = (*operation_outcome_go_proto.OperationOutcome)(nil)
var _ fhir.Resource = (*organization_go_proto.Organization)(nil)
var _ fhir.Resource = (*organization_affiliation_go_proto.OrganizationAffiliation)(nil)
var _ fhir.Resource = (*patient_go_proto.Patient)(nil)
var _ fhir.Resource = (*payment_notice_go_proto.PaymentNotice)(nil)
var _ fhir.Resource = (*payment_reconciliation_go_proto.PaymentReconciliation)(nil)
var _ fhir.Resource = (*person_go_proto.Person)(nil)
var _ fhir.Resource = (*plan_definition_go_proto.PlanDefinition)(nil)
var _ fhir.Resource = (*practitioner_go_proto.Practitioner)(nil)
var _ fhir.Resource = (*practitioner_role_go_proto.PractitionerRole)(nil)
var _ fhir.Resource = (*procedure_go_proto.Procedure)(nil)
var _ fhir.Resource = (*provenance_go_proto.Provenance)(nil)
var _ fhir.Resource = (*questionnaire_go_proto.Questionnaire)(nil)
var _ fhir.Resource = (*questionnaire_response_go_proto.QuestionnaireResponse)(nil)
var _ fhir.Resource = (*related_person_go_proto.RelatedPerson)(nil)
var _ fhir.Resource = (*request_group_go_proto.RequestGroup)(nil)
var _ fhir.Resource = (*research_definition_go_proto.ResearchDefinition)(nil)
var _ fhir.Resource = (*research_element_definition_go_proto.ResearchElementDefinition)(nil)
var _ fhir.Resource = (*research_study_go_proto.ResearchStudy)(nil)
var _ fhir.Resource = (*research_subject_go_proto.ResearchSubject)(nil)
var _ fhir.Resource = (*risk_assessment_go_proto.RiskAssessment)(nil)
var _ fhir.Resource = (*risk_evidence_synthesis_go_proto.RiskEvidenceSynthesis)(nil)
var _ fhir.Resource = (*schedule_go_proto.Schedule)(nil)
var _ fhir.Resource = (*search_parameter_go_proto.SearchParameter)(nil)
var _ fhir.Resource = (*service_request_go_proto.ServiceRequest)(nil)
var _ fhir.Resource = (*slot_go_proto.Slot)(nil)
var _ fhir.Resource = (*specimen_go_proto.Specimen)(nil)
var _ fhir.Resource = (*specimen_definition_go_proto.SpecimenDefinition)(nil)
var _ fhir.Resource = (*structure_definition_go_proto.StructureDefinition)(nil)
var _ fhir.Resource = (*structure_map_go_proto.StructureMap)(nil)
var _ fhir.Resource = (*subscription_go_proto.Subscription)(nil)
var _ fhir.Resource = (*substance_go_proto.Substance)(nil)
var _ fhir.Resource = (*substance_nucleic_acid_go_proto.SubstanceNucleicAcid)(nil)
var _ fhir.Resource = (*substance_polymer_go_proto.SubstancePolymer)(nil)
var _ fhir.Resource = (*substance_protein_go_proto.SubstanceProtein)(nil)
var _ fhir.Resource = (*substance_reference_information_go_proto.SubstanceReferenceInformation)(nil)
var _ fhir.Resource = (*substance_source_material_go_proto.SubstanceSourceMaterial)(nil)
var _ fhir.Resource = (*substance_specification_go_proto.SubstanceSpecification)(nil)
var _ fhir.Resource = (*supply_delivery_go_proto.SupplyDelivery)(nil)
var _ fhir.Resource = (*supply_request_go_proto.SupplyRequest)(nil)
var _ fhir.Resource = (*task_go_proto.Task)(nil)
var _ fhir.Resource = (*terminology_capabilities_go_proto.TerminologyCapabilities)(nil)
var _ fhir.Resource = (*test_report_go_proto.TestReport)(nil)
var _ fhir.Resource = (*test_script_go_proto.TestScript)(nil)
var _ fhir.Resource = (*value_set_go_proto.ValueSet)(nil)
var _ fhir.Resource = (*verification_result_go_proto.VerificationResult)(nil)
var _ fhir.Resource = (*vision_prescription_go_proto.VisionPrescription)(nil)
var _ fhir.Resource = (*parameters_go_proto.Parameters)(nil)
var _ fhir.Resource = (*binary_go_proto.Binary)(nil)
var _ fhir.Resource = (*bundle_and_contained_resource_go_proto.Bundle)(nil)

// DomainResource types (https://www.hl7.org/fhir/resource.html#DomainResource).
// This may exclude certain Trial-Use types that google/fhir doesn't implement.

var _ fhir.DomainResource = (*account_go_proto.Account)(nil)
var _ fhir.DomainResource = (*activity_definition_go_proto.ActivityDefinition)(nil)
var _ fhir.DomainResource = (*adverse_event_go_proto.AdverseEvent)(nil)
var _ fhir.DomainResource = (*allergy_intolerance_go_proto.AllergyIntolerance)(nil)
var _ fhir.DomainResource = (*appointment_go_proto.Appointment)(nil)
var _ fhir.DomainResource = (*appointment_response_go_proto.AppointmentResponse)(nil)
var _ fhir.DomainResource = (*audit_event_go_proto.AuditEvent)(nil)
var _ fhir.DomainResource = (*basic_go_proto.Basic)(nil)
var _ fhir.DomainResource = (*biologically_derived_product_go_proto.BiologicallyDerivedProduct)(nil)
var _ fhir.DomainResource = (*body_structure_go_proto.BodyStructure)(nil)
var _ fhir.DomainResource = (*capability_statement_go_proto.CapabilityStatement)(nil)
var _ fhir.DomainResource = (*care_plan_go_proto.CarePlan)(nil)
var _ fhir.DomainResource = (*care_team_go_proto.CareTeam)(nil)
var _ fhir.DomainResource = (*catalog_entry_go_proto.CatalogEntry)(nil)
var _ fhir.DomainResource = (*charge_item_go_proto.ChargeItem)(nil)
var _ fhir.DomainResource = (*charge_item_definition_go_proto.ChargeItemDefinition)(nil)
var _ fhir.DomainResource = (*claim_go_proto.Claim)(nil)
var _ fhir.DomainResource = (*claim_response_go_proto.ClaimResponse)(nil)
var _ fhir.DomainResource = (*clinical_impression_go_proto.ClinicalImpression)(nil)
var _ fhir.DomainResource = (*code_system_go_proto.CodeSystem)(nil)
var _ fhir.DomainResource = (*communication_go_proto.Communication)(nil)
var _ fhir.DomainResource = (*communication_request_go_proto.CommunicationRequest)(nil)
var _ fhir.DomainResource = (*compartment_definition_go_proto.CompartmentDefinition)(nil)
var _ fhir.DomainResource = (*composition_go_proto.Composition)(nil)
var _ fhir.DomainResource = (*concept_map_go_proto.ConceptMap)(nil)
var _ fhir.DomainResource = (*condition_go_proto.Condition)(nil)
var _ fhir.DomainResource = (*consent_go_proto.Consent)(nil)
var _ fhir.DomainResource = (*contract_go_proto.Contract)(nil)
var _ fhir.DomainResource = (*coverage_go_proto.Coverage)(nil)
var _ fhir.DomainResource = (*coverage_eligibility_request_go_proto.CoverageEligibilityRequest)(nil)
var _ fhir.DomainResource = (*coverage_eligibility_response_go_proto.CoverageEligibilityResponse)(nil)
var _ fhir.DomainResource = (*detected_issue_go_proto.DetectedIssue)(nil)
var _ fhir.DomainResource = (*device_go_proto.Device)(nil)
var _ fhir.DomainResource = (*device_definition_go_proto.DeviceDefinition)(nil)
var _ fhir.DomainResource = (*device_metric_go_proto.DeviceMetric)(nil)
var _ fhir.DomainResource = (*device_request_go_proto.DeviceRequest)(nil)
var _ fhir.DomainResource = (*device_use_statement_go_proto.DeviceUseStatement)(nil)
var _ fhir.DomainResource = (*diagnostic_report_go_proto.DiagnosticReport)(nil)
var _ fhir.DomainResource = (*document_manifest_go_proto.DocumentManifest)(nil)
var _ fhir.DomainResource = (*document_reference_go_proto.DocumentReference)(nil)
var _ fhir.DomainResource = (*effect_evidence_synthesis_go_proto.EffectEvidenceSynthesis)(nil)
var _ fhir.DomainResource = (*encounter_go_proto.Encounter)(nil)
var _ fhir.DomainResource = (*endpoint_go_proto.Endpoint)(nil)
var _ fhir.DomainResource = (*enrollment_request_go_proto.EnrollmentRequest)(nil)
var _ fhir.DomainResource = (*enrollment_response_go_proto.EnrollmentResponse)(nil)
var _ fhir.DomainResource = (*episode_of_care_go_proto.EpisodeOfCare)(nil)
var _ fhir.DomainResource = (*event_definition_go_proto.EventDefinition)(nil)
var _ fhir.DomainResource = (*evidence_go_proto.Evidence)(nil)
var _ fhir.DomainResource = (*evidence_variable_go_proto.EvidenceVariable)(nil)
var _ fhir.DomainResource = (*example_scenario_go_proto.ExampleScenario)(nil)
var _ fhir.DomainResource = (*explanation_of_benefit_go_proto.ExplanationOfBenefit)(nil)
var _ fhir.DomainResource = (*family_member_history_go_proto.FamilyMemberHistory)(nil)
var _ fhir.DomainResource = (*flag_go_proto.Flag)(nil)
var _ fhir.DomainResource = (*goal_go_proto.Goal)(nil)
var _ fhir.DomainResource = (*graph_definition_go_proto.GraphDefinition)(nil)
var _ fhir.DomainResource = (*group_go_proto.Group)(nil)
var _ fhir.DomainResource = (*guidance_response_go_proto.GuidanceResponse)(nil)
var _ fhir.DomainResource = (*healthcare_service_go_proto.HealthcareService)(nil)
var _ fhir.DomainResource = (*imaging_study_go_proto.ImagingStudy)(nil)
var _ fhir.DomainResource = (*immunization_go_proto.Immunization)(nil)
var _ fhir.DomainResource = (*immunization_evaluation_go_proto.ImmunizationEvaluation)(nil)
var _ fhir.DomainResource = (*immunization_recommendation_go_proto.ImmunizationRecommendation)(nil)
var _ fhir.DomainResource = (*implementation_guide_go_proto.ImplementationGuide)(nil)
var _ fhir.DomainResource = (*insurance_plan_go_proto.InsurancePlan)(nil)
var _ fhir.DomainResource = (*invoice_go_proto.Invoice)(nil)
var _ fhir.DomainResource = (*library_go_proto.Library)(nil)
var _ fhir.DomainResource = (*linkage_go_proto.Linkage)(nil)
var _ fhir.DomainResource = (*list_go_proto.List)(nil)
var _ fhir.DomainResource = (*location_go_proto.Location)(nil)
var _ fhir.DomainResource = (*measure_go_proto.Measure)(nil)
var _ fhir.DomainResource = (*measure_report_go_proto.MeasureReport)(nil)
var _ fhir.DomainResource = (*media_go_proto.Media)(nil)
var _ fhir.DomainResource = (*medication_go_proto.Medication)(nil)
var _ fhir.DomainResource = (*medication_administration_go_proto.MedicationAdministration)(nil)
var _ fhir.DomainResource = (*medication_dispense_go_proto.MedicationDispense)(nil)
var _ fhir.DomainResource = (*medication_knowledge_go_proto.MedicationKnowledge)(nil)
var _ fhir.DomainResource = (*medication_request_go_proto.MedicationRequest)(nil)
var _ fhir.DomainResource = (*medication_statement_go_proto.MedicationStatement)(nil)
var _ fhir.DomainResource = (*medicinal_product_go_proto.MedicinalProduct)(nil)
var _ fhir.DomainResource = (*medicinal_product_authorization_go_proto.MedicinalProductAuthorization)(nil)
var _ fhir.DomainResource = (*medicinal_product_contraindication_go_proto.MedicinalProductContraindication)(nil)
var _ fhir.DomainResource = (*medicinal_product_indication_go_proto.MedicinalProductIndication)(nil)
var _ fhir.DomainResource = (*medicinal_product_ingredient_go_proto.MedicinalProductIngredient)(nil)
var _ fhir.DomainResource = (*medicinal_product_interaction_go_proto.MedicinalProductInteraction)(nil)
var _ fhir.DomainResource = (*medicinal_product_manufactured_go_proto.MedicinalProductManufactured)(nil)
var _ fhir.DomainResource = (*medicinal_product_packaged_go_proto.MedicinalProductPackaged)(nil)
var _ fhir.DomainResource = (*medicinal_product_pharmaceutical_go_proto.MedicinalProductPharmaceutical)(nil)
var _ fhir.DomainResource = (*medicinal_product_undesirable_effect_go_proto.MedicinalProductUndesirableEffect)(nil)
var _ fhir.DomainResource = (*message_definition_go_proto.MessageDefinition)(nil)
var _ fhir.DomainResource = (*message_header_go_proto.MessageHeader)(nil)
var _ fhir.DomainResource = (*molecular_sequence_go_proto.MolecularSequence)(nil)
var _ fhir.DomainResource = (*naming_system_go_proto.NamingSystem)(nil)
var _ fhir.DomainResource = (*nutrition_order_go_proto.NutritionOrder)(nil)
var _ fhir.DomainResource = (*observation_go_proto.Observation)(nil)
var _ fhir.DomainResource = (*observation_definition_go_proto.ObservationDefinition)(nil)
var _ fhir.DomainResource = (*operation_definition_go_proto.OperationDefinition)(nil)
var _ fhir.DomainResource = (*operation_outcome_go_proto.OperationOutcome)(nil)
var _ fhir.DomainResource = (*organization_go_proto.Organization)(nil)
var _ fhir.DomainResource = (*organization_affiliation_go_proto.OrganizationAffiliation)(nil)
var _ fhir.DomainResource = (*patient_go_proto.Patient)(nil)
var _ fhir.DomainResource = (*payment_notice_go_proto.PaymentNotice)(nil)
var _ fhir.DomainResource = (*payment_reconciliation_go_proto.PaymentReconciliation)(nil)
var _ fhir.DomainResource = (*person_go_proto.Person)(nil)
var _ fhir.DomainResource = (*plan_definition_go_proto.PlanDefinition)(nil)
var _ fhir.DomainResource = (*practitioner_go_proto.Practitioner)(nil)
var _ fhir.DomainResource = (*practitioner_role_go_proto.PractitionerRole)(nil)
var _ fhir.DomainResource = (*procedure_go_proto.Procedure)(nil)
var _ fhir.DomainResource = (*provenance_go_proto.Provenance)(nil)
var _ fhir.DomainResource = (*questionnaire_go_proto.Questionnaire)(nil)
var _ fhir.DomainResource = (*questionnaire_response_go_proto.QuestionnaireResponse)(nil)
var _ fhir.DomainResource = (*related_person_go_proto.RelatedPerson)(nil)
var _ fhir.DomainResource = (*request_group_go_proto.RequestGroup)(nil)
var _ fhir.DomainResource = (*research_definition_go_proto.ResearchDefinition)(nil)
var _ fhir.DomainResource = (*research_element_definition_go_proto.ResearchElementDefinition)(nil)
var _ fhir.DomainResource = (*research_study_go_proto.ResearchStudy)(nil)
var _ fhir.DomainResource = (*research_subject_go_proto.ResearchSubject)(nil)
var _ fhir.DomainResource = (*risk_assessment_go_proto.RiskAssessment)(nil)
var _ fhir.DomainResource = (*risk_evidence_synthesis_go_proto.RiskEvidenceSynthesis)(nil)
var _ fhir.DomainResource = (*schedule_go_proto.Schedule)(nil)
var _ fhir.DomainResource = (*search_parameter_go_proto.SearchParameter)(nil)
var _ fhir.DomainResource = (*service_request_go_proto.ServiceRequest)(nil)
var _ fhir.DomainResource = (*slot_go_proto.Slot)(nil)
var _ fhir.DomainResource = (*specimen_go_proto.Specimen)(nil)
var _ fhir.DomainResource = (*specimen_definition_go_proto.SpecimenDefinition)(nil)
var _ fhir.DomainResource = (*structure_definition_go_proto.StructureDefinition)(nil)
var _ fhir.DomainResource = (*structure_map_go_proto.StructureMap)(nil)
var _ fhir.DomainResource = (*subscription_go_proto.Subscription)(nil)
var _ fhir.DomainResource = (*substance_go_proto.Substance)(nil)
var _ fhir.DomainResource = (*substance_nucleic_acid_go_proto.SubstanceNucleicAcid)(nil)
var _ fhir.DomainResource = (*substance_polymer_go_proto.SubstancePolymer)(nil)
var _ fhir.DomainResource = (*substance_protein_go_proto.SubstanceProtein)(nil)
var _ fhir.DomainResource = (*substance_reference_information_go_proto.SubstanceReferenceInformation)(nil)
var _ fhir.DomainResource = (*substance_source_material_go_proto.SubstanceSourceMaterial)(nil)
var _ fhir.DomainResource = (*substance_specification_go_proto.SubstanceSpecification)(nil)
var _ fhir.DomainResource = (*supply_delivery_go_proto.SupplyDelivery)(nil)
var _ fhir.DomainResource = (*supply_request_go_proto.SupplyRequest)(nil)
var _ fhir.DomainResource = (*task_go_proto.Task)(nil)
var _ fhir.DomainResource = (*terminology_capabilities_go_proto.TerminologyCapabilities)(nil)
var _ fhir.DomainResource = (*test_report_go_proto.TestReport)(nil)
var _ fhir.DomainResource = (*test_script_go_proto.TestScript)(nil)
var _ fhir.DomainResource = (*value_set_go_proto.ValueSet)(nil)
var _ fhir.DomainResource = (*verification_result_go_proto.VerificationResult)(nil)
var _ fhir.DomainResource = (*vision_prescription_go_proto.VisionPrescription)(nil)

// CanonicalResource types from https://www.hl7.org/fhir/canonicalresource.html#bnr.
// The list of CanonicalResources here is slightly smaller than on HL7's
// site since it includes trial-use types, and objects with trial-use fields
// which are not modeled in the google/fhir protos.

var _ fhir.CanonicalResource = (*activity_definition_go_proto.ActivityDefinition)(nil)
var _ fhir.CanonicalResource = (*code_system_go_proto.CodeSystem)(nil)
var _ fhir.CanonicalResource = (*event_definition_go_proto.EventDefinition)(nil)
var _ fhir.CanonicalResource = (*library_go_proto.Library)(nil)
var _ fhir.CanonicalResource = (*measure_go_proto.Measure)(nil)
var _ fhir.CanonicalResource = (*message_definition_go_proto.MessageDefinition)(nil)
var _ fhir.CanonicalResource = (*plan_definition_go_proto.PlanDefinition)(nil)
var _ fhir.CanonicalResource = (*questionnaire_go_proto.Questionnaire)(nil)
var _ fhir.CanonicalResource = (*research_definition_go_proto.ResearchDefinition)(nil)
var _ fhir.CanonicalResource = (*research_element_definition_go_proto.ResearchElementDefinition)(nil)
var _ fhir.CanonicalResource = (*structure_definition_go_proto.StructureDefinition)(nil)
var _ fhir.CanonicalResource = (*structure_map_go_proto.StructureMap)(nil)
var _ fhir.CanonicalResource = (*value_set_go_proto.ValueSet)(nil)

// MetadataResource types from https://www.hl7.org/fhir/metadataresource.html#bnr.
// The list of MetadataResources here is slightly smaller than on HL7's
// site since it includes trial-use types, and objects with trial-use fields
// which are not modeled in the google/fhir protos.

var _ fhir.MetadataResource = (*activity_definition_go_proto.ActivityDefinition)(nil)
var _ fhir.MetadataResource = (*event_definition_go_proto.EventDefinition)(nil)
var _ fhir.MetadataResource = (*library_go_proto.Library)(nil)
var _ fhir.MetadataResource = (*measure_go_proto.Measure)(nil)
var _ fhir.MetadataResource = (*plan_definition_go_proto.PlanDefinition)(nil)
var _ fhir.MetadataResource = (*research_definition_go_proto.ResearchDefinition)(nil)
var _ fhir.MetadataResource = (*research_element_definition_go_proto.ResearchElementDefinition)(nil)
