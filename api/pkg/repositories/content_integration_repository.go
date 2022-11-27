package repositories

import (
	"context"

	"github.com/google/uuid"

	"github.com/NdoleStudio/superbutton/pkg/entities"
)

// ContentIntegrationRepository loads and persists an entities.WhatsappIntegration
type ContentIntegrationRepository interface {
	// Store a new entities.Project
	Store(ctx context.Context, integration *entities.ContentIntegration) error

	// Update a new entities.ContentIntegration
	Update(ctx context.Context, integration *entities.ContentIntegration) error

	// Fetch all entities.ContentIntegration for a user
	Fetch(ctx context.Context, userID entities.UserID, projectID uuid.UUID) ([]*entities.ContentIntegration, error)

	// FetchMultiple returns multiple entities.ContentIntegration by userID
	FetchMultiple(ctx context.Context, userID entities.UserID, projectIDs []uuid.UUID) ([]*entities.ContentIntegration, error)

	// Delete an entities.WhatsappIntegration
	Delete(ctx context.Context, userID entities.UserID, integrationID uuid.UUID) error

	// Load an entities.ContentIntegration by entities.UserID and integrationID
	Load(ctx context.Context, userID entities.UserID, integrationID uuid.UUID) (*entities.ContentIntegration, error)
}
