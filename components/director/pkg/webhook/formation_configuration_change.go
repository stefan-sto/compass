package webhook

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/kyma-incubator/compass/components/director/internal/model"
)

// ApplicationWithLabels represents an application with its corresponding labels
type ApplicationWithLabels struct {
	*model.Application
	Labels map[string]interface{}
}

// ApplicationTemplateWithLabels represents an application template with its corresponding labels
type ApplicationTemplateWithLabels struct {
	*model.ApplicationTemplate
	Labels map[string]interface{}
}

// RuntimeWithLabels represents a runtime with its corresponding labels
type RuntimeWithLabels struct {
	*model.Runtime
	Labels map[string]interface{}
}

// RuntimeContextWithLabels represents runtime context with its corresponding labels
type RuntimeContextWithLabels struct {
	*model.RuntimeContext
	Labels map[string]interface{}
}

// FormationConfigurationChangeInput struct contains the input for a formation notification
type FormationConfigurationChangeInput struct {
	Operation           model.FormationOperation
	FormationID         string
	ApplicationTemplate *ApplicationTemplateWithLabels
	Application         *ApplicationWithLabels
	Runtime             *RuntimeWithLabels
	RuntimeContext      *RuntimeContextWithLabels
}

// ParseURLTemplate missing godoc
func (rd *FormationConfigurationChangeInput) ParseURLTemplate(tmpl *string) (*URL, error) {
	var url URL
	return &url, parseTemplate(tmpl, *rd, &url)
}

// ParseInputTemplate missing godoc
func (rd *FormationConfigurationChangeInput) ParseInputTemplate(tmpl *string) ([]byte, error) {
	res := json.RawMessage{}
	if err := parseTemplate(tmpl, *rd, &res); err != nil {
		return nil, err
	}
	res = bytes.ReplaceAll(res, []byte("<nil>"), nil)
	return res, nil
}

// ParseHeadersTemplate missing godoc
func (rd *FormationConfigurationChangeInput) ParseHeadersTemplate(tmpl *string) (http.Header, error) {
	var headers http.Header
	return headers, parseTemplate(tmpl, *rd, &headers)
}
