package requests

import (
	"github.com/NdoleStudio/superbutton/pkg/entities"
	"github.com/NdoleStudio/superbutton/pkg/services"
	"github.com/google/uuid"
)

// PhoneCallIntegrationCreateRequest is the payload for the /projects/:projectID/phone-call-integration endpoint
type PhoneCallIntegrationCreateRequest struct {
	request
	ProjectID   string `json:"projectID" swaggerignore:"true"`
	Name        string `json:"name"`
	Text        string `json:"text"`
	PhoneNumber string `json:"phone_number"`
}

// Sanitize the request by stripping whitespaces
func (request *PhoneCallIntegrationCreateRequest) Sanitize() *PhoneCallIntegrationCreateRequest {
	request.Name = request.sanitizeString(request.Name)
	request.Text = request.sanitizeString(request.Text)
	request.PhoneNumber = request.sanitizePhoneNumber(request.PhoneNumber)
	return request
}

// ToCreateParams creates services.PhoneCallIntegrationCreateParams
func (request *PhoneCallIntegrationCreateRequest) ToCreateParams(source string, userID entities.UserID) *services.PhoneCallIntegrationCreateParams {
	return &services.PhoneCallIntegrationCreateParams{
		Name:        request.Name,
		Text:        request.Text,
		PhoneNumber: request.PhoneNumber,
		ProjectID:   uuid.MustParse(request.ProjectID),
		Source:      source,
		UserID:      userID,
	}
}
