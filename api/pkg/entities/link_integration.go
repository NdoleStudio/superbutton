package entities

import (
	"time"

	"github.com/google/uuid"
)

// LinkIntegration are URL integration settings
type LinkIntegration struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:string;" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	UserID    UserID    `json:"user_id" example:"WB7DRDWrJZRGbYrv2CKGkqbzvqdC"`
	ProjectID uuid.UUID `json:"project_id" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	Enabled   bool      `json:"enabled" example:"true"`
	Name      string    `json:"name" example:"FAQ"`
	Text      string    `json:"text" example:"Visit our FAQ"`
	URL       string    `json:"url" example:"https://example.com"`
	Icon      string    `json:"icon" example:"url"`
	CreatedAt time.Time `json:"created_at" example:"2022-06-05T14:26:02.302718+03:00"`
	UpdatedAt time.Time `json:"updated_at" example:"2022-06-05T14:26:10.303278+03:00"`
}

// Integration converts LinkIntegration to Integration
func (integration *LinkIntegration) Integration() *Integration {
	return &Integration{
		UserID:        integration.UserID,
		ProjectID:     integration.ProjectID,
		IntegrationID: integration.ID,
		Type:          IntegrationTypeLink,
		Name:          integration.Name,
		CreatedAt:     integration.CreatedAt,
		UpdatedAt:     integration.UpdatedAt,
	}
}

func (integration *LinkIntegration) NewProjectIntegration(order uint) *ProjectIntegration {
	return &ProjectIntegration{
		ID:            uuid.New(),
		UserID:        integration.UserID,
		ProjectID:     integration.ProjectID,
		IntegrationID: integration.ID,
		Type:          IntegrationTypeLink,
		Name:          integration.Name,
		Position:      order,
		CreatedAt:     integration.CreatedAt,
		UpdatedAt:     integration.UpdatedAt,
	}
}
