package requests

import (
	"github.com/NdoleStudio/superbutton/pkg/entities"
	"github.com/NdoleStudio/superbutton/pkg/services"
	"github.com/google/uuid"
)

// ContentIntegrationCreateRequest is the payload for the /projects/:projectID/text-integrations endpoint
type ContentIntegrationCreateRequest struct {
	request
	ProjectID string `json:"projectID" swaggerignore:"true"`
	Name      string `json:"name"`
	Text      string `json:"text"`
	Summary   string `json:"summary"`
	Title     string `json:"title"`
}

// Sanitize the request by stripping whitespaces
func (request *ContentIntegrationCreateRequest) Sanitize() *ContentIntegrationCreateRequest {
	request.Name = request.sanitizeString(request.Name)
	request.Text = request.sanitizeString(request.Text)
	request.Summary = request.sanitizeString(request.Summary)
	request.Title = request.sanitizeString(request.Title)
	return request
}

// ToCreateParams creates services.ContentIntegrationCreateParams
func (request *ContentIntegrationCreateRequest) ToCreateParams(source string, userID entities.UserID) *services.ContentIntegrationCreateParams {
	return &services.ContentIntegrationCreateParams{
		Name:      request.Name,
		Text:      request.Text,
		Summary:   request.Summary,
		Title:     request.Title,
		ProjectID: uuid.MustParse(request.ProjectID),
		Source:    source,
		UserID:    userID,
	}
}
