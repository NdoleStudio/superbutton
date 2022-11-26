package services

import (
	"context"
	"fmt"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"

	"github.com/NdoleStudio/superbutton/pkg/entities"
	"github.com/NdoleStudio/superbutton/pkg/events"
	"github.com/NdoleStudio/superbutton/pkg/repositories"
	"github.com/NdoleStudio/superbutton/pkg/telemetry"
	"github.com/google/uuid"
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

// ProjectCreateParams are the parameters for creating a new project.
type ProjectCreateParams struct {
	Name   string
	Source string
	URL    string
	UserID entities.UserID
}

// Create a new entities.Project
func (service *ProjectService) Create(ctx context.Context, params *ProjectCreateParams) (*entities.Project, error) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	project := &entities.Project{
		ID:                     uuid.New(),
		UserID:                 params.UserID,
		URL:                    params.URL,
		CreatedAt:              time.Now().UTC(),
		UpdatedAt:              time.Now().UTC(),
		Name:                   params.Name,
		Icon:                   "chat",
		Greeting:               "Need some help?",
		GreetingTimeoutSeconds: 10,
		Color:                  "#283593",
	}

	err := service.repository.Store(ctx, project)
	if err != nil {
		msg := fmt.Sprintf("could store project for user with ID [%s]", params.UserID)
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	service.dispatchProjectCreatedEvent(ctx, params.Source, project)

	return project, nil
}

// ProjectUpdateParams are the parameters for updating a project.
type ProjectUpdateParams struct {
	UserID                 entities.UserID
	ProjectID              uuid.UUID
	Name                   string
	URL                    string
	Icon                   string
	Greeting               string
	Source                 string
	GreetingTimeoutSeconds uint
	Color                  string
}

// Update an entities.Project
func (service *ProjectService) Update(ctx context.Context, params *ProjectUpdateParams) (*entities.Project, error) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	project, err := service.repository.Load(ctx, params.UserID, params.ProjectID)
	if err != nil {
		msg := fmt.Sprintf("cannot load project for user ID [%s] and project [%s]", params.UserID, params.ProjectID)
		return nil, stacktrace.PropagateWithCode(err, stacktrace.GetCode(err), msg)
	}

	project.Name = params.Name
	project.URL = params.URL
	project.UpdatedAt = time.Now().UTC()
	project.Icon = params.Icon
	project.GreetingTimeoutSeconds = params.GreetingTimeoutSeconds
	project.Greeting = params.Greeting

	if err = service.repository.Update(ctx, project); err != nil {
		msg := fmt.Sprintf("could update project [%s] for user with ID [%s]", project.ID, project.UserID)
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	service.dispatchProjectUpdatedEvent(ctx, params.Source, project)

	return project, nil
}

func (service *ProjectService) dispatchProjectUpdatedEvent(ctx context.Context, source string, project *entities.Project) {
	ctx, span, ctxLogger := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	event, err := service.createEvent(events.ProjectUpdated, source, &events.ProjectUpdatedPayload{
		UserID:                 project.UserID,
		ProjectID:              project.ID,
		ProjectName:            project.Name,
		ProjectURL:             project.URL,
		ProjectIcon:            project.Icon,
		ProjectGreeting:        project.Greeting,
		ProjectColor:           project.Color,
		ProjectGreetingTimeout: project.GreetingTimeoutSeconds,
		ProjectUpdatedAt:       project.UpdatedAt,
	})
	if err != nil {
		msg := fmt.Sprintf("cannot create [%s] event for project [%s]", events.ProjectUpdated, project.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return
	}

	service.dispatchEvent(ctx, project, event)
}

func (service *ProjectService) dispatchProjectCreatedEvent(ctx context.Context, source string, project *entities.Project) {
	ctx, span, ctxLogger := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	event, err := service.createEvent(events.ProjectCreated, source, &events.ProjectCreatedPayload{
		UserID:           project.UserID,
		ProjectCreatedAt: project.CreatedAt,
		ProjectID:        project.ID,
		ProjectName:      project.Name,
		ProjectURL:       project.URL,
	})
	if err != nil {
		msg := fmt.Sprintf("cannot create [%s] event for project [%s]", events.ProjectCreated, project.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return
	}

	service.dispatchEvent(ctx, project, event)
}

func (service *ProjectService) dispatchEvent(ctx context.Context, project *entities.Project, event cloudevents.Event) {
	ctx, span, ctxLogger := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	if err := service.eventDispatcher.Dispatch(ctx, event); err != nil {
		msg := fmt.Sprintf("cannot dispatch [%s] event for project [%s]", event.Type(), project.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return
	}
}
