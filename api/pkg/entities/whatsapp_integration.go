package entities

import (
	"time"

	"github.com/google/uuid"
)

// WhatsappIntegration contains whatsapp integration settings
type WhatsappIntegration struct {
	ID          uuid.UUID `json:"id" gorm:"primaryKey;type:string;" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	UserID      UserID    `json:"user_id" example:"WB7DRDWrJZRGbYrv2CKGkqbzvqdC"`
	ProjectID   uuid.UUID `json:"project_id" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	Enabled     bool      `json:"enabled" example:"true"`
	Name        string    `json:"name" example:"FAQ"`
	Text        string    `json:"text" example:"Contact us on WhatsApp"`
	PhoneNumber string    `json:"phone_number" example:"+18005550199"`
	Icon        string    `json:"icon" example:"https://cdn.superbutton.app/whatsapp-icon.svg"`
	CreatedAt   time.Time `json:"created_at" example:"2022-06-05T14:26:02.302718+03:00"`
	UpdatedAt   time.Time `json:"updated_at" example:"2022-06-05T14:26:10.303278+03:00"`
}

// Integration converts WhatsappIntegration to Integration
func (integration *WhatsappIntegration) Integration() *Integration {
	return &Integration{
		UserID:        integration.UserID,
		ProjectID:     integration.ProjectID,
		IntegrationID: integration.ID,
		Type:          IntegrationTypeWhatsapp,
		CreatedAt:     integration.CreatedAt,
		UpdatedAt:     integration.UpdatedAt,
	}
}

func (integration *WhatsappIntegration) NewProjectIntegration(order uint) *ProjectIntegration {
	return &ProjectIntegration{
		ID:            uuid.New(),
		UserID:        integration.UserID,
		ProjectID:     integration.ProjectID,
		IntegrationID: integration.ID,
		Type:          IntegrationTypeWhatsapp,
		Name:          integration.Name,
		Position:      order,
		CreatedAt:     integration.CreatedAt,
		UpdatedAt:     integration.UpdatedAt,
	}
}
