package ndoapi

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrNotFound identifies an HTTP 404 from an NDO GET request.
var ErrNotFound = errors.New("NDO object not found")

// ErrTemplatePathNotFound identifies a missing path inside a fetched template.
var ErrTemplatePathNotFound = errors.New("template path not found")

// TemplatePathNotFoundError includes the field at which a template lookup stopped.
type TemplatePathNotFoundError struct {
	TemplateID string
	Field      string
}

func (err *TemplatePathNotFoundError) Error() string {
	if err.Field == "" {
		return fmt.Sprintf("template %q: %s", err.TemplateID, ErrTemplatePathNotFound)
	}
	return fmt.Sprintf("template %q: field %q: %s", err.TemplateID, err.Field, ErrTemplatePathNotFound)
}

func (err *TemplatePathNotFoundError) Unwrap() error {
	return ErrTemplatePathNotFound
}

func errorForObjectNotFound(endpoint string, response *http.Response, err error) error {
	if response != nil && response.StatusCode == http.StatusNotFound {
		return fmt.Errorf("GET %q returned HTTP 404: %w", endpoint, ErrNotFound)
	}
	return err
}
