package repositories

import (
	"context"

	"github.com/google/uuid"

	"github.com/NdoleStudio/superbutton/pkg/entities"
)

// LinkIntegrationRepository loads and persists an entities.LinkIntegration
type LinkIntegrationRepository interface {
	// Store a new entities.Project
	Store(ctx context.Context, integration *entities.LinkIntegration) error

	// Update a new entities.LinkIntegration
	Update(ctx context.Context, integration *entities.LinkIntegration) error

	// Fetch all entities.LinkIntegration for a user
	Fetch(ctx context.Context, userID entities.UserID, projectID uuid.UUID) ([]*entities.LinkIntegration, error)

	// FetchMultiple returns multiple entities.LinkIntegration by userID
	FetchMultiple(ctx context.Context, userID entities.UserID, projectIDs []uuid.UUID) ([]*entities.LinkIntegration, error)

	// Delete an entities.LinkIntegration
	Delete(ctx context.Context, userID entities.UserID, integrationID uuid.UUID) error

	// Load an entities.LinkIntegration by entities.UserID and integrationID
	Load(ctx context.Context, userID entities.UserID, integrationID uuid.UUID) (*entities.LinkIntegration, error)
}
