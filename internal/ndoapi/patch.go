package ndoapi

import (
	"fmt"
	"maps"
	"reflect"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
)

// PatchOperations builds JSON-Patch operations by comparing a desired object
// with the object currently returned by NDO.
type PatchOperations struct {
	existing   map[string]any
	pathPrefix string
	operations []*PatchPayload
}

// NewPatchOperations creates an operation builder for a top-level NDO object
// at the supplied JSON-Patch path prefix.
func NewPatchOperations(existing map[string]any, pathPrefix string) *PatchOperations {
	existingCopy := make(map[string]any, len(existing))
	maps.Copy(existingCopy, existing)

	return &PatchOperations{
		existing:   existingCopy,
		pathPrefix: pathPrefix,
		operations: make([]*PatchPayload, 0),
	}
}

// Set records the desired value for a top-level NDO object attribute.
func (operations *PatchOperations) Set(attribute string, value any) {
	existingValue, exists := operations.existing[attribute]
	if exists && reflect.DeepEqual(existingValue, value) {
		return
	}

	operation := "add"
	if exists {
		operation = "replace"
	}
	operations.operations = append(operations.operations, NewPatchPayload(operation, operations.attributePath(attribute), value))
}

// Remove records that a top-level NDO object attribute should be removed.
func (operations *PatchOperations) Remove(attribute string) {
	if _, exists := operations.existing[attribute]; !exists {
		return
	}
	operations.operations = append(operations.operations, NewRemovePatchPayload(operations.attributePath(attribute)))
}

// SetString records a known string value after applying an optional
// transformation. A nil value is unmanaged and ignored. When clearEmpty is
// true, an empty transformed value records a remove operation instead.
func (operations *PatchOperations) SetString(attribute string, value *string, transform func(string) string, clearEmpty bool) {
	if value == nil {
		return
	}

	transformedValue := *value
	if transform != nil {
		transformedValue = transform(transformedValue)
	}
	if clearEmpty && transformedValue == "" {
		operations.Remove(attribute)
		return
	}
	operations.Set(attribute, transformedValue)
}

// Operations returns the minimal JSON-Patch operations recorded in the
// builder.
func (operations *PatchOperations) Operations() []*PatchPayload {
	return append([]*PatchPayload(nil), operations.operations...)
}

// Append adds operations built for child paths to the same template patch.
func (operations *PatchOperations) Append(patches ...*PatchPayload) {
	operations.operations = append(operations.operations, patches...)
}

func (operations *PatchOperations) attributePath(attribute string) string {
	return operations.pathPrefix + "/" + attribute
}

// Apply applies the recorded JSON-Patch operations to an NDO endpoint.
func (operations *PatchOperations) Apply(msoClient *client.Client, endpoint string) error {
	patchOperations, err := operations.toModels()
	if err != nil {
		return err
	}
	if len(patchOperations) == 0 {
		return nil
	}

	_, err = msoClient.PatchbyID(endpoint, patchOperations...)
	return err
}

func (operations *PatchOperations) toModels() ([]models.Model, error) {
	patchOperations := make([]models.Model, 0, len(operations.operations))
	for index, operation := range operations.Operations() {
		if operation == nil {
			return nil, fmt.Errorf("patch operation %d is nil", index)
		}

		modelOperation := *operation
		patchOperations = append(patchOperations, &modelOperation)
	}
	return patchOperations, nil
}

// PatchPayload is a generic NDO JSON-Patch operation. Value accepts any JSON
// value, including scalar values that the mso-go-client models.PatchPayload
// type cannot represent.
type PatchPayload struct {
	Operation string
	Path      string
	Value     any
	HasValue  bool
}

// NewPatchPayload creates a JSON-Patch operation with a value.
func NewPatchPayload(operation, path string, value any) *PatchPayload {
	return &PatchPayload{
		Operation: operation,
		Path:      path,
		Value:     value,
		HasValue:  true,
	}
}

// NewRemovePatchPayload creates a JSON-Patch remove operation.
func NewRemovePatchPayload(path string) *PatchPayload {
	return &PatchPayload{
		Operation: "remove",
		Path:      path,
	}
}

// ToMap implements models.Model for use with mso-go-client PatchbyID.
func (payload *PatchPayload) ToMap() (map[string]any, error) {
	result := map[string]any{
		"op":   payload.Operation,
		"path": payload.Path,
	}
	if payload.HasValue {
		result["value"] = payload.Value
	}
	return result, nil
}
