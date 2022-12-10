package services

import (
	"context"
	"fmt"

	"github.com/NdoleStudio/superbutton/pkg/entities"
	"github.com/NdoleStudio/superbutton/pkg/repositories"
	"github.com/NdoleStudio/superbutton/pkg/telemetry"
	"github.com/google/uuid"
	"github.com/palantir/stacktrace"
)

type ProjectIntegrationService struct {
	tracer          telemetry.Tracer
	logger          telemetry.Logger
	eventDispatcher *EventDispatcher
	repository      repositories.ProjectIntegrationRepository
}

// NewProjectIntegrationService creates a new ProjectIntegrationService
func NewProjectIntegrationService(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	eventDispatcher *EventDispatcher,
	repository repositories.ProjectIntegrationRepository,
) (s *ProjectIntegrationService) {
	return &ProjectIntegrationService{
		logger:          logger.WithService(fmt.Sprintf("%T", s)),
		tracer:          tracer,
		eventDispatcher: eventDispatcher,
		repository:      repository,
	}
}

// Index fetches all entities.ProjectIntegration for an authenticated user
func (service *ProjectIntegrationService) Index(ctx context.Context, userID entities.UserID, projectID uuid.UUID) ([]*entities.ProjectIntegration, error) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	integrations, err := service.repository.Fetch(ctx, userID, projectID)
	if err != nil {
		msg := fmt.Sprintf("could fetch project integrations for user with ID [%s] and project [%s]", userID, projectID)
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return integrations, nil
}

// Update updates the positions for entities.ProjectIntegration for an authenticated user
func (service *ProjectIntegrationService) Update(ctx context.Context, userID entities.UserID, integrationIDs []uuid.UUID) error {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	if err := service.repository.UpdatePositions(ctx, userID, integrationIDs); err != nil {
		msg := fmt.Sprintf("could update project integrations for user with ID [%s] and project [%s]", userID, integrationIDs)
		return service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}
