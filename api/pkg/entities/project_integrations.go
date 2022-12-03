package entities

import (
	"time"

	"github.com/google/uuid"
)

type IntegrationType string

const (
	IntegrationTypeWhatsapp    = IntegrationType("whatsapp")
	IntegrationTypePhoneCall   = IntegrationType("phone-call")
	IntegrationTypeLink        = IntegrationType("link")
	IntegrationTypeContent     = IntegrationType("content")
	IntegrationTypeContactForm = IntegrationType("contact-form")
)

type ProjectIntegration struct {
	ID            uuid.UUID       `json:"id" gorm:"primaryKey;type:string;" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	UserID        UserID          `json:"user_id" example:"WB7DRDWrJZRGbYrv2CKGkqbzvqdC"`
	ProjectID     uuid.UUID       `json:"project_id" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	IntegrationID uuid.UUID       `json:"integration_id" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	Type          IntegrationType `json:"type" example:"whatsapp"`
	Name          string          `json:"name"`
	Position      uint            `json:"position" example:"1"`
	CreatedAt     time.Time       `json:"created_at" example:"2022-06-05T14:26:02.302718+03:00"`
	UpdatedAt     time.Time       `json:"updated_at" example:"2022-06-05T14:26:10.303278+03:00"`
}
