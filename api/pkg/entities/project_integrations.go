package entities

import "github.com/google/uuid"

type IntegrationType string

const (
	IntegrationTypeWhatsapp    = IntegrationType("whatsapp")
	IntegrationTypePhoneCall   = IntegrationType("phone-call")
	IntegrationTypeURL         = IntegrationType("url")
	IntegrationTypeText        = IntegrationType("text")
	IntegrationTypeContactForm = IntegrationType("contact-form")
)

type ProjectIntegration struct {
	ID            uuid.UUID       `json:"id" gorm:"primaryKey;type:string;" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	UserID        UserID          `json:"user_id" example:"WB7DRDWrJZRGbYrv2CKGkqbzvqdC"`
	ProjectID     uuid.UUID       `json:"project_id" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	IntegrationID uuid.UUID       `json:"integration_id" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	Type          IntegrationType `json:"type" example:"whatsapp"`
	Order         uint            `json:"order" example:"1"`
}
