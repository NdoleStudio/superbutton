package events

import (
	"time"

	"github.com/google/uuid"

	"github.com/NdoleStudio/superbutton/pkg/entities"
)

// IntegrationDeleted is raised when an integration is created
const IntegrationDeleted = "integration.deleted"

// IntegrationDeletedPayload stores the data for the IntegrationDeleted  event
type IntegrationDeletedPayload struct {
	UserID               entities.UserID          `json:"user_id"`
	ProjectID            uuid.UUID                `json:"project_id"`
	IntegrationID        uuid.UUID                `json:"integration_id"`
	IntegrationType      entities.IntegrationType `json:"integration_type"`
	IntegrationName      string                   `json:"integration_name"`
	IntegrationDeletedAt time.Time                `json:"integration_deleted_at"`
}
