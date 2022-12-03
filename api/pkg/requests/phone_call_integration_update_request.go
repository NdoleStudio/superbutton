package requests

import (
	"github.com/NdoleStudio/superbutton/pkg/entities"
	"github.com/NdoleStudio/superbutton/pkg/services"
	"github.com/google/uuid"
)

// PhoneCallIntegrationUpdateRequest is the payload for the /projects/:projectID/phone-call-integrations/:integrationID endpoint
type PhoneCallIntegrationUpdateRequest struct {
	request
	IntegrationID string `json:"integrationID" swaggerignore:"true"`
	Name          string `json:"name"`
	Text          string `json:"text"`
	PhoneNumber   string `json:"phone_number"`
}

// Sanitize the request by stripping whitespaces
func (request *PhoneCallIntegrationUpdateRequest) Sanitize() *PhoneCallIntegrationUpdateRequest {
	request.Name = request.sanitizeString(request.Name)
	request.Text = request.sanitizeString(request.Text)
	request.PhoneNumber = request.sanitizePhoneNumber(request.PhoneNumber)
	return request
}

// ToUpdateParams creates services.PhoneCallIntegrationUpdateParams
func (request *PhoneCallIntegrationUpdateRequest) ToUpdateParams(source string, userID entities.UserID) *services.PhoneCallIntegrationUpdateParams {
	return &services.PhoneCallIntegrationUpdateParams{
		Name:          request.Name,
		Text:          request.Text,
		PhoneNumber:   request.PhoneNumber,
		IntegrationID: uuid.MustParse(request.IntegrationID),
		Source:        source,
		UserID:        userID,
	}
}
