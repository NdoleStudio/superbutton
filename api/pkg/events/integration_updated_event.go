package events

import (
	"time"

	"github.com/google/uuid"

	"github.com/NdoleStudio/superbutton/pkg/entities"
)

// IntegrationUpdated is raised when an integration is updated
const IntegrationUpdated = "integration.updated"

// IntegrationUpdatedPayload stores the data for the IntegrationUpdated  event
type IntegrationUpdatedPayload struct {
	UserID               entities.UserID          `json:"user_id"`
	ProjectID            uuid.UUID                `json:"project_id"`
	IntegrationID        uuid.UUID                `json:"integration_id"`
	IntegrationType      entities.IntegrationType `json:"integration_type"`
	IntegrationName      string                   `json:"integration_name"`
	IntegrationUpdatedAt time.Time                `json:"integration_updated_at"`
}
