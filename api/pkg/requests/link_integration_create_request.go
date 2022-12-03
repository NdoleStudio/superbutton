package requests

import (
	"github.com/NdoleStudio/superbutton/pkg/entities"
	"github.com/NdoleStudio/superbutton/pkg/services"
	"github.com/google/uuid"
)

// LinkIntegrationCreateRequest is the payload for the /projects/:projectID/link-integrations endpoint
type LinkIntegrationCreateRequest struct {
	request
	ProjectID string `json:"projectID" swaggerignore:"true"`
	Name      string `json:"name"`
	Text      string `json:"text"`
	Website   string `json:"website"`
}

// Sanitize the request by stripping whitespaces
func (request *LinkIntegrationCreateRequest) Sanitize() *LinkIntegrationCreateRequest {
	request.Name = request.sanitizeString(request.Name)
	request.Text = request.sanitizeString(request.Text)
	request.Website = request.sanitizeString(request.Website)
	return request
}

// ToCreateParams creates services.PhoneCallIntegrationCreateParams
func (request *LinkIntegrationCreateRequest) ToCreateParams(source string, userID entities.UserID) *services.LinkIntegrationCreateParams {
	return &services.LinkIntegrationCreateParams{
		Name:      request.Name,
		Text:      request.Text,
		URL:       request.Website,
		ProjectID: uuid.MustParse(request.ProjectID),
		Source:    source,
		UserID:    userID,
	}
}
