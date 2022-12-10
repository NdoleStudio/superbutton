package requests

import (
	"github.com/NdoleStudio/superbutton/pkg/entities"
	"github.com/NdoleStudio/superbutton/pkg/services"
	"github.com/google/uuid"
)

// LinkIntegrationUpdateRequest struct { is the payload for the /projects/:projectID/links-integrations/:integrationID endpoint
type LinkIntegrationUpdateRequest struct {
	request
	IntegrationID string `json:"integrationID" swaggerignore:"true"`
	Name          string `json:"name"`
	Text          string `json:"text"`
	Website       string `json:"website"`
	Icon          string `json:"icon"`
	PhoneNumber   string `json:"phone_number"`
}

// Sanitize the request by stripping whitespaces
func (request *LinkIntegrationUpdateRequest) Sanitize() *LinkIntegrationUpdateRequest {
	request.Name = request.sanitizeString(request.Name)
	request.Text = request.sanitizeString(request.Text)
	request.Icon = request.sanitizeString(request.Icon)
	request.PhoneNumber = request.sanitizePhoneNumber(request.PhoneNumber)
	return request
}

// ToUpdateParams creates services.PhoneCallIntegrationUpdateParams
func (request *LinkIntegrationUpdateRequest) ToUpdateParams(source string, userID entities.UserID) *services.LinkIntegrationUpdateParams {
	return &services.LinkIntegrationUpdateParams{
		Name:          request.Name,
		Text:          request.Text,
		Icon:          request.Icon,
		URL:           request.Website,
		IntegrationID: uuid.MustParse(request.IntegrationID),
		Source:        source,
		UserID:        userID,
	}
}
