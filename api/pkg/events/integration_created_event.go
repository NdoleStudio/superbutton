package events

import (
	"time"

	"github.com/google/uuid"

	"github.com/NdoleStudio/superbutton/pkg/entities"
)

// IntegrationCreated is raised when an integration is created
const IntegrationCreated = "integration.created"

// IntegrationCreatedPayload stores the data for the IntegrationCreated  event
type IntegrationCreatedPayload struct {
	UserID               entities.UserID          `json:"user_id"`
	ProjectID            uuid.UUID                `json:"project_id"`
	IntegrationID        uuid.UUID                `json:"integration_id"`
	IntegrationType      entities.IntegrationType `json:"integration_type"`
	IntegrationName      string                   `json:"integration_name"`
	IntegrationCreatedAt time.Time                `json:"integration_created_at"`
}
