package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/NdoleStudio/superbutton/pkg/entities"
	"github.com/NdoleStudio/superbutton/pkg/events"
	"github.com/NdoleStudio/superbutton/pkg/repositories"
	"github.com/NdoleStudio/superbutton/pkg/telemetry"
	"github.com/palantir/stacktrace"
)

// ProjectService is responsible for managing entities.Project
type ProjectService struct {
	service
	logger          telemetry.Logger
	tracer          telemetry.Tracer
	repository      repositories.ProjectRepository
	eventDispatcher *EventDispatcher
}

// NewProjectService creates a new ProjectService
func NewProjectService(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	eventDispatcher *EventDispatcher,
	repository repositories.ProjectRepository,
) (s *ProjectService) {
	return &ProjectService{
		logger:          logger.WithService(fmt.Sprintf("%T", s)),
		tracer:          tracer,
		eventDispatcher: eventDispatcher,
		repository:      repository,
	}
}

// Index fetches all entities.Project for an authenticated user
func (service *ProjectService) Index(ctx context.Context, userID entities.UserID) ([]*entities.Project, error) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	projects, err := service.repository.Fetch(ctx, userID)
	if err != nil {
		msg := fmt.Sprintf("could fetch projects for user with ID [%s]", userID)
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return projects, nil
}

type ProjectCreateParams struct {
	Name    string
	Source  string
	Website string
	UserID  entities.UserID
}

// Create a new entities.Project
func (service *ProjectService) Create(ctx context.Context, params ProjectCreateParams) (*entities.Project, error) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	project := &entities.Project{
		ID:        uuid.New(),
		UserID:    params.UserID,
		URL:       params.Website,
		Settings:  nil,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err := service.repository.Store(ctx, project)
	if err != nil {
		msg := fmt.Sprintf("could store projects for user with ID [%s]", params.UserID)
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	service.dispatchProjectCreatedEvent(ctx, params.Source, project)

	return project, nil
}

func (service *ProjectService) dispatchProjectCreatedEvent(ctx context.Context, source string, project *entities.Project) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	ctxLogger := service.tracer.CtxLogger(service.logger, span)

	event, err := service.createEvent(events.ProjectCreated, source, &events.ProjectCreatedPayload{
		UserID:           project.UserID,
		ProjectCreatedAt: project.CreatedAt,
		ProjectID:        project.ID,
		ProjectName:      project.Name,
		ProjectURL:       project.URL,
	})
	if err != nil {
		msg := fmt.Sprintf("cannot created [%s] event for project [%s]", events.ProjectCreated, project.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return
	}

	if err = service.eventDispatcher.Dispatch(ctx, event); err != nil {
		msg := fmt.Sprintf("cannot dispatch [%s] event for project [%s]", events.ProjectCreated, project.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return
	}
}
