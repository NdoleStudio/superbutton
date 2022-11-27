package entities

import (
	"time"

	"github.com/google/uuid"
)

// ContentIntegration contains content integration settings
type ContentIntegration struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:string;" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	UserID    UserID    `json:"user_id" example:"WB7DRDWrJZRGbYrv2CKGkqbzvqdC"`
	ProjectID uuid.UUID `json:"project_id" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	Enabled   bool      `json:"enabled" example:"true"`
	Name      string    `json:"name" example:"FAQ"`
	Title     string    `json:"title" example:"What is SuperButton?"`
	Summary   string    `json:"summary" example:"Configurable floating button for your website"`
	Text      string    `json:"text" example:"SuperButton is the best app to create configurable floating buttons on your website."`
	CreatedAt time.Time `json:"created_at" example:"2022-06-05T14:26:02.302718+03:00"`
	UpdatedAt time.Time `json:"updated_at" example:"2022-06-05T14:26:10.303278+03:00"`
}

// Integration converts WhatsappIntegration to Integration
func (integration *ContentIntegration) Integration() *Integration {
	return &Integration{
		UserID:        integration.UserID,
		ProjectID:     integration.ProjectID,
		IntegrationID: integration.ID,
		Type:          IntegrationTypeContent,
		CreatedAt:     integration.CreatedAt,
		UpdatedAt:     integration.UpdatedAt,
	}
}

func (integration *ContentIntegration) NewProjectIntegration(position uint) *ProjectIntegration {
	return &ProjectIntegration{
		ID:            uuid.New(),
		UserID:        integration.UserID,
		ProjectID:     integration.ProjectID,
		IntegrationID: integration.ID,
		Type:          IntegrationTypeContent,
		Name:          integration.Name,
		Position:      position,
		CreatedAt:     integration.CreatedAt,
		UpdatedAt:     integration.UpdatedAt,
	}
}
