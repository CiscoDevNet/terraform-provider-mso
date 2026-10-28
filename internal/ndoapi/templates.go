// Package ndoapi contains shared helpers for the Nexus Dashboard Orchestrator API.
package ndoapi

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/ciscoecosystem/mso-go-client/client"
)

const templatesEndpoint = "api/v1/templates"

// TemplateEndpoint returns the NDO endpoint for a template UUID.
func TemplateEndpoint(templateID string) string {
	return fmt.Sprintf("%s/%s", templatesEndpoint, url.PathEscape(templateID))
}

// Template is one decoded response from the NDO template endpoint.
type Template struct {
	id     string
	object map[string]any
}

// Find resolves a model path in the fetched template. Absence is allowed;
// malformed response fields return an error.
func (template *Template) Find(path Path) (ResolvedObject, bool, error) {
	resolved, found, err := path.ResolveObjectFromResponse(template.object)
	if err != nil {
		return ResolvedObject{}, false, fmt.Errorf("template %q: %w", template.id, err)
	}
	return resolved, found, nil
}

// FindRequired resolves a model path. Missing paths return a typed not-found
// error; malformed response fields return a decoding error.
func (template *Template) FindRequired(path Path) (ResolvedObject, error) {
	resolved, missingField, found, err := path.resolveObjectFromResponse(template.object)
	if err != nil {
		return ResolvedObject{}, fmt.Errorf("template %q: %w", template.id, err)
	}
	if !found {
		return ResolvedObject{}, &TemplatePathNotFoundError{TemplateID: template.id, Field: missingField}
	}
	return resolved, nil
}

// AppendPath resolves the parent of a collection and returns its patch append
// path. Missing parents and malformed response fields return distinct errors.
func (template *Template) AppendPath(path Path) (string, error) {
	appendPath, missingField, found, err := path.appendPathFromResponse(template.object)
	if err != nil {
		return "", fmt.Errorf("template %q: %w", template.id, err)
	}
	if !found {
		return "", &TemplatePathNotFoundError{TemplateID: template.id, Field: missingField}
	}
	return appendPath, nil
}

// GetTemplate retrieves and decodes an NDO template by UUID.
func GetTemplate(msoClient *client.Client, templateID string) (*Template, error) {
	endpoint := TemplateEndpoint(templateID)
	request, err := msoClient.MakeRestRequest(http.MethodGet, endpoint, nil, true)
	if err != nil {
		return nil, fmt.Errorf("template %q: %w", templateID, err)
	}
	response, httpResponse, err := msoClient.Do(request)
	if err = errorForObjectNotFound(endpoint, httpResponse, err); err != nil {
		return nil, fmt.Errorf("template %q: %w", templateID, err)
	}
	if response == nil {
		return nil, fmt.Errorf("template %q: empty response body", templateID)
	}
	if err := client.CheckForErrors(response, http.MethodGet); err != nil {
		return nil, fmt.Errorf("template %q: %w", templateID, err)
	}
	object, ok := response.Data().(map[string]any)
	if !ok {
		return nil, fmt.Errorf("unexpected template response type %T", response.Data())
	}
	return &Template{id: templateID, object: object}, nil
}
