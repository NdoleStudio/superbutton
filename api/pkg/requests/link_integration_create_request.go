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
	Icon      string `json:"icon"`
	Color     string `json:"color"`
	Website   string `json:"website"`
}

// Sanitize the request by stripping whitespaces
func (request *LinkIntegrationCreateRequest) Sanitize() *LinkIntegrationCreateRequest {
	request.Name = request.sanitizeString(request.Name)
	request.Text = request.sanitizeString(request.Text)
	request.Icon = request.sanitizeString(request.Icon)
	request.Website = request.sanitizeString(request.Website)

	request.Color = request.sanitizeString(request.Color)
	if request.Color == "" {
		request.Color = "#1E88E5"
	}

	return request
}

// ToCreateParams creates services.PhoneCallIntegrationCreateParams
func (request *LinkIntegrationCreateRequest) ToCreateParams(source string, userID entities.UserID) *services.LinkIntegrationCreateParams {
	return &services.LinkIntegrationCreateParams{
		Name:      request.Name,
		Text:      request.Text,
		URL:       request.Website,
		Color:     request.Color,
		Icon:      request.Icon,
		ProjectID: uuid.MustParse(request.ProjectID),
		Source:    source,
		UserID:    userID,
	}
}
