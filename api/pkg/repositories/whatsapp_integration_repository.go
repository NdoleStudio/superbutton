package repositories

import (
	"context"

	"github.com/google/uuid"

	"github.com/NdoleStudio/superbutton/pkg/entities"
)

// WhatsappIntegrationRepository loads and persists an entities.WhatsappIntegration
type WhatsappIntegrationRepository interface {
	// Store a new entities.Project
	Store(ctx context.Context, integration *entities.WhatsappIntegration) error

	// Update a new entities.WhatsappIntegration
	Update(ctx context.Context, integration *entities.WhatsappIntegration) error

	// Fetch all entities.WhatsappIntegration for a user
	Fetch(ctx context.Context, userID entities.UserID, projectID uuid.UUID) ([]*entities.WhatsappIntegration, error)

	// Delete an entities.WhatsappIntegration
	Delete(ctx context.Context, userID entities.UserID, integrationID uuid.UUID) error

	// Load an entities.WhatsappIntegration by entities.UserID and integrationID
	Load(ctx context.Context, userID entities.UserID, integrationID uuid.UUID) (*entities.WhatsappIntegration, error)
}
