package ndoapi

import (
	"fmt"
)

// Translation maps provider values to API values and back again.
//
// The input table is copied when the translation is created. Unmapped values
// are returned unchanged to preserve forward compatibility with newer API
// values.
type Translation struct {
	schemaToAPI map[string]string
	apiToSchema map[string]string
}

// NewTranslationMap creates a translation for a hardcoded schema mapping.
// It panics when the mapping cannot be reversed unambiguously, which indicates
// a provider programming error.
func NewTranslationMap(values map[string]string) Translation {
	schemaToAPI := make(map[string]string, len(values))
	apiToSchema := make(map[string]string, len(values))

	for schemaValue, apiValue := range values {
		if existingSchemaValue, exists := apiToSchema[apiValue]; exists {
			panic(fmt.Sprintf("API value %q is mapped from both %q and %q", apiValue, existingSchemaValue, schemaValue))
		}
		schemaToAPI[schemaValue] = apiValue
		apiToSchema[apiValue] = schemaValue
	}

	return Translation{schemaToAPI: schemaToAPI, apiToSchema: apiToSchema}
}

// ToAPI translates a schema value to the API representation.
func (translation Translation) ToAPI(value string) string {
	if translatedValue, ok := translation.schemaToAPI[value]; ok {
		return translatedValue
	}
	return value
}

// ToSchema translates an API value to the schema representation.
func (translation Translation) ToSchema(value string) string {
	if translatedValue, ok := translation.apiToSchema[value]; ok {
		return translatedValue
	}
	return value
}
