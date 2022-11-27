package requests

import (
	"github.com/NdoleStudio/superbutton/pkg/entities"
	"github.com/NdoleStudio/superbutton/pkg/services"
	"github.com/google/uuid"
)

// ContentIntegrationUpdateRequest is the payload for the /projects/:projectID/text-integrations/:integrationID endpoint
type ContentIntegrationUpdateRequest struct {
	request
	IntegrationID string `json:"integrationID" swaggerignore:"true"`
	Name          string `json:"name"`
	Text          string `json:"text"`
	Summary       string `json:"summary"`
	Title         string `json:"title"`
}

// Sanitize the request by stripping whitespaces
func (request *ContentIntegrationUpdateRequest) Sanitize() *ContentIntegrationUpdateRequest {
	request.Name = request.sanitizeString(request.Name)
	request.Text = request.sanitizeString(request.Text)
	request.Title = request.sanitizeString(request.Title)
	request.Summary = request.sanitizeString(request.Summary)
	return request
}

// ToUpdateParams creates services.ContentIntegrationUpdateParams
func (request *ContentIntegrationUpdateRequest) ToUpdateParams(source string, userID entities.UserID) *services.ContentIntegrationUpdateParams {
	return &services.ContentIntegrationUpdateParams{
		Name:          request.Name,
		Text:          request.Text,
		Title:         request.Title,
		Summary:       request.Summary,
		IntegrationID: uuid.MustParse(request.IntegrationID),
		Source:        source,
		UserID:        userID,
	}
}
