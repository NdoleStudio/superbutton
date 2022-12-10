package repositories

import (
	"context"

	"github.com/google/uuid"

	"github.com/NdoleStudio/superbutton/pkg/entities"
)

// ProjectIntegrationRepository loads and persists an entities.ProjectIntegration
type ProjectIntegrationRepository interface {
	// Fetch all entities.ProjectIntegration for a project
	Fetch(ctx context.Context, userID entities.UserID, projectID uuid.UUID) ([]*entities.ProjectIntegration, error)

	// UpdatePositions updates the positions of multiple project integrations
	UpdatePositions(ctx context.Context, userID entities.UserID, integrationIDs []uuid.UUID) error
}
