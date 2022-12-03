package entities

import (
	"time"

	"github.com/google/uuid"
)

// PhoneCallIntegration contains phone call integration settings
type PhoneCallIntegration struct {
	ID          uuid.UUID `json:"id" gorm:"primaryKey;type:string;" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	UserID      UserID    `json:"user_id" example:"WB7DRDWrJZRGbYrv2CKGkqbzvqdC"`
	ProjectID   uuid.UUID `json:"project_id" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	Enabled     bool      `json:"enabled" example:"true"`
	Text        string    `json:"text" example:"Call us on +18005550199"`
	Name        string    `json:"name" example:"Customer Service"`
	PhoneNumber string    `json:"phone_number" example:"+18005550199"`
	Icon        string    `json:"icon" example:"phone-call"`
	CreatedAt   time.Time `json:"created_at" example:"2022-06-05T14:26:02.302718+03:00"`
	UpdatedAt   time.Time `json:"updated_at" example:"2022-06-05T14:26:10.303278+03:00"`
}

// Integration converts WhatsappIntegration to Integration
func (integration *PhoneCallIntegration) Integration() *Integration {
	return &Integration{
		UserID:        integration.UserID,
		ProjectID:     integration.ProjectID,
		IntegrationID: integration.ID,
		Type:          IntegrationTypePhoneCall,
		Name:          integration.Name,
		CreatedAt:     integration.CreatedAt,
		UpdatedAt:     integration.UpdatedAt,
	}
}

func (integration *PhoneCallIntegration) NewProjectIntegration(order uint) *ProjectIntegration {
	return &ProjectIntegration{
		ID:            uuid.New(),
		UserID:        integration.UserID,
		ProjectID:     integration.ProjectID,
		IntegrationID: integration.ID,
		Type:          IntegrationTypePhoneCall,
		Name:          integration.Name,
		Position:      order,
		CreatedAt:     integration.CreatedAt,
		UpdatedAt:     integration.UpdatedAt,
	}
}
