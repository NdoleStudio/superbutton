package repositories

import (
	"context"

	"github.com/google/uuid"

	"github.com/NdoleStudio/superbutton/pkg/entities"
)

// PhoneCallIntegrationRepository loads and persists an entities.PhoneIntegration
type PhoneCallIntegrationRepository interface {
	// Store a new entities.Project
	Store(ctx context.Context, integration *entities.PhoneCallIntegration) error

	// Update a new entities.PhoneCallIntegration
	Update(ctx context.Context, integration *entities.PhoneCallIntegration) error

	// Fetch all entities.PhoneCallIntegration for a user
	Fetch(ctx context.Context, userID entities.UserID, projectID uuid.UUID) ([]*entities.PhoneCallIntegration, error)

	// FetchMultiple returns multiple entities.PhoneCallIntegration by userID
	FetchMultiple(ctx context.Context, userID entities.UserID, projectIDs []uuid.UUID) ([]*entities.PhoneCallIntegration, error)

	// Delete an entities.PhoneCallIntegration
	Delete(ctx context.Context, userID entities.UserID, integrationID uuid.UUID) error

	// Load an entities.PhoneCallIntegration by entities.UserID and integrationID
	Load(ctx context.Context, userID entities.UserID, integrationID uuid.UUID) (*entities.PhoneCallIntegration, error)
}
