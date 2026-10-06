package ndoapi

import (
	"fmt"
	"math"
)

// FieldRequirement controls whether a decoded NDO field may be absent.
type FieldRequirement uint8

const (
	OptionalField FieldRequirement = iota
	RequiredField
	// RequiredNonEmptyField applies to string fields only.
	RequiredNonEmptyField
)

func fieldValue(object map[string]any, key string, requirement FieldRequirement) (any, bool, error) {
	value, exists := object[key]
	if !exists || value == nil {
		if requirement != OptionalField {
			return nil, false, fmt.Errorf("NDO field %q is required", key)
		}
		return nil, false, nil
	}
	return value, true, nil
}

func unexpectedFieldType(key string, value any) error {
	return fmt.Errorf("NDO field %q has unexpected type %T", key, value)
}

// StringField reads a string without coercing other JSON types.
func StringField(object map[string]any, key string, requirement FieldRequirement) (string, bool, error) {
	value, exists, err := fieldValue(object, key, requirement)
	if err != nil || !exists {
		return "", exists, err
	}
	stringValue, ok := value.(string)
	if !ok {
		return "", true, unexpectedFieldType(key, value)
	}
	if requirement == RequiredNonEmptyField && stringValue == "" {
		return "", true, fmt.Errorf("NDO field %q must be nonempty", key)
	}
	return stringValue, true, nil
}

// BoolField reads a boolean from a decoded NDO object.
func BoolField(object map[string]any, key string, requirement FieldRequirement) (bool, bool, error) {
	value, exists, err := fieldValue(object, key, requirement)
	if err != nil || !exists {
		return false, exists, err
	}
	boolean, ok := value.(bool)
	if !ok {
		return false, true, unexpectedFieldType(key, value)
	}
	return boolean, true, nil
}

// MapField reads a decoded NDO object field.
func MapField(object map[string]any, key string, requirement FieldRequirement) (map[string]any, bool, error) {
	value, exists, err := fieldValue(object, key, requirement)
	if err != nil || !exists {
		return nil, exists, err
	}
	field, ok := value.(map[string]any)
	if !ok {
		return nil, true, unexpectedFieldType(key, value)
	}
	return field, true, nil
}

// ListField reads a decoded NDO JSON array.
func ListField(object map[string]any, key string, requirement FieldRequirement) ([]any, bool, error) {
	value, exists, err := fieldValue(object, key, requirement)
	if err != nil || !exists {
		return nil, exists, err
	}
	field, ok := value.([]any)
	if !ok {
		return nil, true, unexpectedFieldType(key, value)
	}
	return field, true, nil
}

// Int64Field reads a whole number representable as int64.
func Int64Field(object map[string]any, key string, requirement FieldRequirement) (int64, bool, error) {
	value, exists, err := fieldValue(object, key, requirement)
	if err != nil || !exists {
		return 0, exists, err
	}
	switch number := value.(type) {
	case int:
		return int64(number), true, nil
	case int32:
		return int64(number), true, nil
	case int64:
		return number, true, nil
	case float64:
		if !math.IsNaN(number) && !math.IsInf(number, 0) && number >= -0x1p63 && number < 0x1p63 && math.Trunc(number) == number {
			return int64(number), true, nil
		}
		return 0, true, fmt.Errorf("NDO field %q must be a whole number within int64 range", key)
	default:
		return 0, true, unexpectedFieldType(key, value)
	}
}
