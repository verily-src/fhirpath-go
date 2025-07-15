package element

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/slices"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protopath"
	"google.golang.org/protobuf/reflect/protorange"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	ErrFhirPathNotImplemented error = errors.New("FHIRPath labeling not yet implemented")
)

// A ElementWithPath holds a FHIR element as a proto message and its corresponding
// FHIRPath. All FHIR elements (including scalar-ish elements like string
// and URI) are proto messages.
// See http://hl7.org/fhir/R4/fhirpath.html.
type ElementWithPath[elementT proto.Message] struct {
	Element  elementT
	FHIRPath string
}

// SortSliceOfElementWithPath sorts the given slice in-place according to the lexical
// ordering of the FHIR path of each element.
func SortSliceOfElementWithPath[elementT proto.Message](s []ElementWithPath[elementT]) {
	sort.SliceStable(s, func(i, j int) bool {
		return s[i].FHIRPath < s[j].FHIRPath
	})
}

// ExtractAll returns all the elements of type elementT within the resource.
// The returned elements are not copies: mutations of the returned elements
// will mutate the resource. The order of the returned list is unspecified.
func ExtractAll[elementT proto.Message](resource fhir.Resource) ([]elementT, error) {
	elements, err := extractAllImpl[elementT](resource, false)
	if err != nil {
		return nil, err
	}
	return slices.Map(elements, func(e ElementWithPath[elementT]) elementT {
		return e.Element
	}), nil
}

// ExtractAllWithPath returns all the elements of typeT within the resource,
// along with the FHIR path of each such element. The returned elements
// are not copies: mutations of the returned elements will mutate the resource.
// The returned list is sorted by the lexical order of the FHIR path of each
// element.
func ExtractAllWithPath[elementT proto.Message](resource fhir.Resource) ([]ElementWithPath[elementT], error) {
	elements, err := extractAllImpl[elementT](resource, true)
	if err != nil {
		return nil, err
	}
	SortSliceOfElementWithPath(elements)
	return elements, err
}

// extractAllImpl extracts all the instances of elementT from resource.
// If addPaths is true, all returns elements include the FHIR path
// of that element; otherwise each element's FHIR path is an empty string.
// The FHIR path labelling is optional because it may trigger error
// conditions that otherwise would not occur. (Also it is slower).
//
// Once FHIR path labelling is proven robust consider always doing it
// because it yields a nice determinstic sort.
func extractAllImpl[elementT proto.Message](resource fhir.Resource, addPaths bool) ([]ElementWithPath[elementT], error) {
	elements := []ElementWithPath[elementT]{}
	err := protorange.Range(resource.ProtoReflect(), func(pv protopath.Values) error {
		element, found := getElementOfProtoPath[elementT](pv)
		if found {
			var fhirpath string
			if addPaths {
				var err error
				fhirpath, err = computeFHIRPathOfProtoPath(pv.Path)
				if err != nil {
					return err
				}
			}
			elementWithPath := ElementWithPath[elementT]{Element: element, FHIRPath: fhirpath}
			elements = append(elements, elementWithPath)
		}
		return nil
	})
	return elements, err
}

// getElementOfProtoPath returns FHIR element referenced by pv.
//
// Returns (element, true) if pv references an element of type elementT;
// else returns (typed-nil, false).
func getElementOfProtoPath[elementT proto.Message](pv protopath.Values) (elementT, bool) {
	currStep := pv.Path.Index(-1)
	currV := pv.Values[pv.Len()-1]
	var prm protoreflect.Message
	switch currStep.Kind() {
	case protopath.FieldAccessStep:
		// FieldAccess describes access of a field within a message.
		// The type of the current step value is determined by the field descriptor.
		fd := currStep.FieldDescriptor()
		if fd.Kind() == protoreflect.MessageKind && fd.Cardinality() != protoreflect.Repeated {
			prm = currV.Message()
		}
	case protopath.ListIndexStep:
		// ListIndex describes index of an element within a list.
		// The previous step value is always a list.
		// The previous step type is a field access for lists with message kind.
		prevStep := pv.Path.Index(-2)
		prevV := pv.Values[pv.Len()-2]
		if prevStep.Kind() != protopath.FieldAccessStep {
			break
		}
		if prevStep.FieldDescriptor().Kind() == protoreflect.MessageKind {
			prm = prevV.List().Get(currStep.ListIndex()).Message()
		}
	case protopath.MapIndexStep:
		// MapIndex describes index of an entry within a map.
		// The previous step value is always a map.
		// The previous step type is a field access for maps with message kind.
		prevStep := pv.Path.Index(-2)
		prevV := pv.Values[pv.Len()-2]
		if prevStep.Kind() != protopath.FieldAccessStep {
			break
		}
		if prevStep.FieldDescriptor().Kind() == protoreflect.MessageKind {
			prm = prevV.Map().Get(currStep.MapIndex()).Message()
		}
	}
	if prm != nil {
		if element, ok := prm.Interface().(elementT); ok {
			return element, true
		}
	}
	var nilElement elementT
	return nilElement, false
}

// leafElementsByMsgFullName is set of Descriptor full names of elements
// that have proto implementation with subfields that don't matter
// at the FHIR element level. When computing the FHIR path of a proto path,
// proto fields of these message types are ignored.
var leafElementsByMsgFullName = map[string]struct{}{
	"google.fhir.r4.core.Date":     {},
	"google.fhir.r4.core.DateTime": {},
	"google.fhir.r4.core.Instant":  {},
	"google.fhir.r4.core.Time":     {},
}

// computeFHIRPathOfProtoPath returns the FHIR path of p.
//
// The following cases are supported:
//   - Typical single and repeated elements.
//   - extensions. Extensions are expessed as "Resource.extension[k]" instead
//     of as "Resource.extension('<extension-url>')" because the former is unique
//     even with repeated extensions of the same URL. (And because
//     it is simpler to implement here.)
//   - "choice" fields.
//   - Special date fields: see above leafElementsByMsgFullName.
//
// The following cases will return an error:
//   - Any element inside a ContainedResource. Typically this happens
//     for a Bundle.
//
// WATCHOUT: There are almost certainly cases that do not return an error
// but return an incorrect FHIRPath.
func computeFHIRPathOfProtoPath(p protopath.Path) (string, error) {
	fhirpath := []string{}
	for _, step := range p {
		switch step.Kind() {
		case protopath.RootStep:
			fhirpath = append(fhirpath, string(step.MessageDescriptor().Name()))

		case protopath.FieldAccessStep:
			fd := step.FieldDescriptor()
			cfn := string(fd.ContainingMessage().FullName())
			if cfn == "google.fhir.r4.core.ContainedResource" {
				// Only elements of Resource have well defined FHIRpath within
				// Bundle.entry.resource. See ticket. Punt on this for now.
				return "", fmt.Errorf("%w: for %s", ErrFhirPathNotImplemented, cfn)
			}
			elementName := fd.JSONName()
			if cof := fd.ContainingOneof(); cof != nil && cof.Name() == "choice" {
				cappedName := strings.ToUpper(elementName[0:1]) + elementName[1:]
				fhirpath[len(fhirpath)-1] += cappedName
			} else {
				fhirpath = append(fhirpath, elementName)
			}
			if _, found := leafElementsByMsgFullName[cfn]; found {
				// The fhirpath component above is suffient. The remaining protopath
				// steps are implementation fields that do not appear in the spec.
				break
			}

		case protopath.ListIndexStep:
			fhirpath[len(fhirpath)-1] += fmt.Sprintf("[%d]", step.ListIndex())

		default:
			return "", fmt.Errorf("%w: for step %v", ErrFhirPathNotImplemented, step)
		}
	}
	return strings.Join(fhirpath, "."), nil
}
