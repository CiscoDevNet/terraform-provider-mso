package ndoapi

import (
	"fmt"
	"strconv"
	"strings"
)

// ObjectIdentifier matches one string field of an object in an NDO collection.
type ObjectIdentifier struct {
	Field string
	Value string
}

// IsUnset reports whether an identifier has no usable field or value.
func (identifier ObjectIdentifier) IsUnset() bool {
	return identifier.Field == "" || identifier.Value == ""
}

// ObjectSelector uses UUID when present, otherwise every key identifier.
// A UUID that does not match never falls back to Keys.
type ObjectSelector struct {
	UUID ObjectIdentifier
	Keys []ObjectIdentifier
}

func (selector ObjectSelector) matches(object map[string]any) (bool, error) {
	if !selector.UUID.IsUnset() {
		return identifierMatches(object, selector.UUID)
	}
	if len(selector.Keys) == 0 {
		return false, nil
	}
	for _, identifier := range selector.Keys {
		if identifier.IsUnset() {
			return false, nil
		}
		matched, err := identifierMatches(object, identifier)
		if err != nil || !matched {
			return false, err
		}
	}
	return true, nil
}

func identifierMatches(object map[string]any, identifier ObjectIdentifier) (bool, error) {
	raw, exists := object[identifier.Field]
	if !exists {
		return false, nil
	}
	value, ok := raw.(string)
	if !ok {
		return false, unexpectedFieldType(identifier.Field, raw)
	}
	return value == identifier.Value, nil
}

// ResolvedObject contains the selected object and its JSON-Patch path.
type ResolvedObject struct {
	Object map[string]any
	path   string
}

func (resolved ResolvedObject) PatchPath() string {
	return resolved.path
}

// PathStep traverses a map field or selects an object in a collection. A nil
// selector is a map field; a non-nil selector identifies a collection.
type PathStep struct {
	Field    string
	Selector *ObjectSelector
}

// Path describes a location in an NDO response without fixed array indexes.
// NewPath copies the steps so callers cannot change a path after construction.
type Path struct {
	steps      []PathStep
	appendPath string
}

// NewPath constructs a path from model-defined steps. Collection selectors and
// their key identifiers are copied along with the steps. When all parent steps are
// plain map fields, it also prepares the collection's JSON-Patch append path.
func NewPath(steps ...PathStep) Path {
	path := Path{steps: make([]PathStep, len(steps))}
	for index, step := range steps {
		path.steps[index] = step
		if step.Selector != nil {
			selector := *step.Selector
			selector.Keys = append([]ObjectIdentifier(nil), selector.Keys...)
			path.steps[index].Selector = &selector
		}
	}

	if !path.endsInCollection() {
		return path
	}
	var patchPath strings.Builder
	for index, step := range path.steps {
		if step.Field == "" || (index < len(path.steps)-1 && step.Selector != nil) {
			return path
		}
		appendPathValue(&patchPath, step.Field)
	}
	appendPathValue(&patchPath, "-")
	path.appendPath = patchPath.String()
	return path
}

// WithSteps returns a new path with nested fields or collection selectors.
// The original path and the supplied selectors remain independent of it.
func (path Path) WithSteps(steps ...PathStep) Path {
	return NewPath(append(path.steps, steps...)...)
}

// ResolveObjectFromResponse finds an existing object in a decoded NDO response.
// Use it for reads, updates, or deletes: it returns the object and its
// JSON-Patch path, resolving selected collections to their current indexes.
// A missing selector returns found=false; malformed response fields return an error.
func (path Path) ResolveObjectFromResponse(response map[string]any) (ResolvedObject, bool, error) {
	resolved, _, found, err := path.resolveObjectFromResponse(response)
	return resolved, found, err
}

// resolveObjectFromResponse also identifies the field where traversal stopped.
func (path Path) resolveObjectFromResponse(response map[string]any) (ResolvedObject, string, bool, error) {
	if !path.endsInCollection() {
		return ResolvedObject{}, "", false, nil
	}
	object, patchPath, missingField, found, err := walkPath(response, path.steps)
	if err != nil || !found {
		return ResolvedObject{}, missingField, false, err
	}
	return ResolvedObject{Object: object, path: patchPath}, "", true, nil
}

// AppendPath returns the JSON-Patch path for adding an object to the final
// collection. A nested collection needs the NDO response to resolve its parent
// selectors; a path with only map parents needs no response. A missing parent
// returns found=false, while malformed response fields return an error.
func (path Path) AppendPath(response ...map[string]any) (string, bool, error) {
	if len(response) > 1 {
		return "", false, fmt.Errorf("append path accepts at most one response")
	}
	var object map[string]any
	if len(response) == 1 {
		object = response[0]
	}
	appendPath, _, found, err := path.appendPathFromResponse(object)
	return appendPath, found, err
}

// appendPathFromResponse also identifies the unresolved parent field.
func (path Path) appendPathFromResponse(response map[string]any) (string, string, bool, error) {
	if path.appendPath != "" {
		return path.appendPath, "", true, nil
	}
	if response == nil || !path.endsInCollection() || len(path.steps) < 2 {
		return "", "", false, nil
	}
	_, parentPath, missingField, found, err := walkPath(response, path.steps[:len(path.steps)-1])
	if err != nil || !found {
		return "", missingField, false, err
	}
	if path.steps[len(path.steps)-1].Field == "" {
		return "", "", false, nil
	}
	var appendPath strings.Builder
	appendPath.WriteString(parentPath)
	appendPathValue(&appendPath, path.steps[len(path.steps)-1].Field)
	appendPathValue(&appendPath, "-")
	return appendPath.String(), "", true, nil
}

func (path Path) endsInCollection() bool {
	return len(path.steps) > 0 && path.steps[len(path.steps)-1].Selector != nil
}

// walkPath resolves each selected collection within the previously matched
// object and records the indexes used in the JSON-Patch path.
func walkPath(response map[string]any, steps []PathStep) (map[string]any, string, string, bool, error) {
	current := response
	var patchPath strings.Builder
	for _, step := range steps {
		// Each step is relative to the object selected by the previous step.
		// Missing fields mean the requested path does not exist. Invalid field
		// types are response errors and must not be treated as deleted objects.
		if step.Field == "" {
			return nil, "", "", false, fmt.Errorf("NDO path contains an empty field")
		}
		value, exists := current[step.Field]
		if !exists || value == nil {
			return nil, "", step.Field, false, nil
		}
		appendPathValue(&patchPath, step.Field)
		// No selector means this field is a nested object; a selector means
		// it is a collection whose matching child must be found below.
		if step.Selector == nil {
			object, ok := value.(map[string]any)
			if !ok {
				return nil, "", "", false, unexpectedFieldType(step.Field, value)
			}
			current = object
			continue
		}
		objects, ok := value.([]any)
		if !ok {
			return nil, "", "", false, unexpectedFieldType(step.Field, value)
		}
		index, object, found, err := findSelectedObject(step.Field, objects, *step.Selector)
		if err != nil {
			return nil, "", "", false, err
		}
		if !found {
			return nil, "", step.Field, false, nil
		}
		// JSON-Patch addresses a collection child by its current array index,
		// even though the model selected it by UUID or key identifiers.
		appendPathValue(&patchPath, strconv.Itoa(index))
		current = object
	}
	return current, patchPath.String(), "", true, nil
}

func appendPathValue(path *strings.Builder, value string) {
	path.WriteByte('/')
	path.WriteString(value)
}

func findSelectedObject(field string, objects []any, selector ObjectSelector) (int, map[string]any, bool, error) {
	var selected map[string]any
	selectedIndex := 0
	for index, raw := range objects {
		object, ok := raw.(map[string]any)
		if !ok {
			return 0, nil, false, fmt.Errorf("NDO field %q[%d] has unexpected type %T", field, index, raw)
		}
		matched, err := selector.matches(object)
		if err != nil {
			return 0, nil, false, fmt.Errorf("NDO field %q[%d]: %w", field, index, err)
		}
		if selected == nil && matched {
			selected = object
			selectedIndex = index
		}
	}
	if selected == nil {
		return 0, nil, false, nil
	}
	return selectedIndex, selected, true, nil
}
