package funcs

import (
	"fmt"
)

// FunctionTable is the data structure used to store
// valid FHIRPath functions, and maps their case-sensitive
// names.
type FunctionTable map[string]Function

// Register attempts to add a given function to the FunctionTable t.
func (t FunctionTable) Register(name string, fn any) error {
	if _, ok := t[name]; ok {
		return fmt.Errorf("function '%s' already exists in default table", name)
	}
	fhirpathFunc, err := ToFunction(fn)
	if err != nil {
		return err
	}
	t[name] = fhirpathFunc
	return nil
}
