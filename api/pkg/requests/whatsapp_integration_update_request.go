package requests

import (
	"github.com/NdoleStudio/superbutton/pkg/entities"
	"github.com/NdoleStudio/superbutton/pkg/services"
	"github.com/google/uuid"
)

// WhatsappIntegrationUpdateRequest is the payload for the /projects/:projectID/whatsapp-integrations/:integrationID endpoint
type WhatsappIntegrationUpdateRequest struct {
	request
	IntegrationID string `json:"integrationID" swaggerignore:"true"`
	Name          string `json:"name"`
	Text          string `json:"text"`
	PhoneNumber   string `json:"phone_number"`
}

// Sanitize the request by stripping whitespaces
func (request *WhatsappIntegrationUpdateRequest) Sanitize() *WhatsappIntegrationUpdateRequest {
	request.Name = request.sanitizeString(request.Name)
	request.Text = request.sanitizeString(request.Text)
	request.PhoneNumber = request.sanitizePhoneNumber(request.PhoneNumber)
	return request
}

// ToUpdateParams creates services.WhatsappIntegrationUpdateParams
func (request *WhatsappIntegrationUpdateRequest) ToUpdateParams(source string, userID entities.UserID) *services.WhatsappIntegrationUpdateParams {
	return &services.WhatsappIntegrationUpdateParams{
		Name:          request.Name,
		Text:          request.Text,
		PhoneNumber:   request.PhoneNumber,
		IntegrationID: uuid.MustParse(request.IntegrationID),
		Source:        source,
		UserID:        userID,
	}
}
