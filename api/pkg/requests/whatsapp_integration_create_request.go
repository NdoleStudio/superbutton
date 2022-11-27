package requests

import (
	"github.com/NdoleStudio/superbutton/pkg/entities"
	"github.com/NdoleStudio/superbutton/pkg/services"
	"github.com/google/uuid"
)

// WhatsappIntegrationCreateRequest is the payload for the /projects/:projectID/whatsapp-integration endpoint
type WhatsappIntegrationCreateRequest struct {
	request
	ProjectID   string `json:"projectID" swaggerignore:"true"`
	Name        string `json:"name"`
	Text        string `json:"text"`
	PhoneNumber string `json:"phone_number"`
}

// Sanitize the request by stripping whitespaces
func (request *WhatsappIntegrationCreateRequest) Sanitize() *WhatsappIntegrationCreateRequest {
	request.Name = request.sanitizeString(request.Name)
	request.Text = request.sanitizeString(request.Text)
	request.PhoneNumber = request.sanitizePhoneNumber(request.PhoneNumber)
	return request
}

// ToCreateParams creates services.WhatsappIntegrationUpdateParams
func (request *WhatsappIntegrationCreateRequest) ToCreateParams(source string, userID entities.UserID) *services.WhatsappIntegrationCreateParams {
	return &services.WhatsappIntegrationCreateParams{
		Name:        request.Name,
		Text:        request.Text,
		PhoneNumber: request.PhoneNumber,
		ProjectID:   uuid.MustParse(request.ProjectID),
		Source:      source,
		UserID:      userID,
	}
}
