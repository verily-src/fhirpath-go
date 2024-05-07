package fhirtest

import (
	"time"

	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/element/meta"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// TouchMeta touches the meta to generate a new version ID and use now as the
// time. In most cases, this is likely what is required for tests -- since code
// should seldomly rely on what the version-id or update times discretely are, and
// this helps to ensure proper uniqueness just as the real FHIRStore would.
func TouchMeta(resource fhir.Resource) {
	UpdateMeta(resource, randomVersionID(), time.Now())
}

// UpdateID will update the resource's ID to the specified resourceID string.
func UpdateID(resource fhir.Resource, resourceID string) {
	reflect := resource.ProtoReflect()
	field := reflect.Descriptor().Fields().ByName("id")

	id := &datatypes_go_proto.Id{
		Value: resourceID,
	}
	reflect.Set(field, protoreflect.ValueOfMessage(id.ProtoReflect()))
}

// UpdateMeta updates the meta contents of the fhir resource to use the new
// version-ID and update-time.
func UpdateMeta(resource fhir.Resource, versionID string, updateTime time.Time) {
	reflect := resource.ProtoReflect()
	metaField := getMetaField(reflect)
	time := fhir.Instant(updateTime)
	version := fhir.ID(versionID)

	if resource.GetMeta() == nil {
		meta := &datatypes_go_proto.Meta{
			LastUpdated: time,
			VersionId:   version,
		}
		reflect.Set(metaField, protoreflect.ValueOfMessage(meta.ProtoReflect()))
	}

	message := reflect.Get(metaField).Message()
	descriptor := message.Descriptor()
	fields := descriptor.Fields()
	updateField := fields.ByName("last_updated")
	versionField := fields.ByName("version_id")

	message.Set(updateField, protoreflect.ValueOfMessage(time.ProtoReflect()))
	message.Set(versionField, protoreflect.ValueOfMessage(version.ProtoReflect()))
}

// NOTE: This method is deprecated and should use the production one in
// "github.com/verily-src/fhirpath-go/internal/element/meta"
func ReplaceMeta(resource fhir.Resource, m *datatypes_go_proto.Meta) {
	meta.ReplaceInResource(resource, m)
}

func getMetaField(reflect protoreflect.Message) protoreflect.FieldDescriptor {
	return reflect.Descriptor().Fields().ByName("meta")
}
