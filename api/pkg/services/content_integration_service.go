package services

import (
	"context"
	"fmt"
	"time"

	"github.com/NdoleStudio/superbutton/pkg/entities"
	"github.com/NdoleStudio/superbutton/pkg/repositories"
	"github.com/NdoleStudio/superbutton/pkg/telemetry"
	"github.com/google/uuid"
	"github.com/palantir/stacktrace"
)

// ContentIntegrationService is responsible for managing entities.ContentIntegration
type ContentIntegrationService struct {
	integrationService
	projectRepository repositories.ProjectRepository
	repository        repositories.ContentIntegrationRepository
}

// NewContentIntegrationService creates a new ContentIntegrationService
func NewContentIntegrationService(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	eventDispatcher *EventDispatcher,
	projectRepository repositories.ProjectRepository,
	repository repositories.ContentIntegrationRepository,
) (s *ContentIntegrationService) {
	return &ContentIntegrationService{
		repository:        repository,
		projectRepository: projectRepository,
		integrationService: integrationService{
			tracer:          tracer,
			logger:          logger.WithService(fmt.Sprintf("%T", s)),
			eventDispatcher: eventDispatcher,
		},
	}
}

// Get returns an entities.WhatsappIntegration for an authenticated user
func (service *ContentIntegrationService) Get(ctx context.Context, userID entities.UserID, integrationID uuid.UUID) (*entities.ContentIntegration, error) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	integration, err := service.repository.Load(ctx, userID, integrationID)
	if err != nil {
		msg := fmt.Sprintf("could whatsapp integrations for user with ID [%s] and ID [%s]", userID, integrationID)
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return integration, nil
}

// Index fetches all entities.Project for an authenticated user
func (service *ContentIntegrationService) Index(ctx context.Context, userID entities.UserID, projectID uuid.UUID) ([]*entities.ContentIntegration, error) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	integrations, err := service.repository.Fetch(ctx, userID, projectID)
	if err != nil {
		msg := fmt.Sprintf("could fetch whatsapp integrations for user with ID [%s] and projectID [%s]", userID, projectID)
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return integrations, nil
}

// ContentIntegrationCreateParams are the parameters for creating a new content integration.
type ContentIntegrationCreateParams struct {
	Name      string
	Summary   string
	Text      string
	Source    string
	Title     string
	ProjectID uuid.UUID
	UserID    entities.UserID
}

// Create a new entities.WhatsappIntegration
func (service *ContentIntegrationService) Create(ctx context.Context, params *ContentIntegrationCreateParams) (*entities.ContentIntegration, error) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	if _, err := service.projectRepository.Load(ctx, params.UserID, params.ProjectID); err != nil {
		msg := fmt.Sprintf("cannot load project [%s] for user ID [%s]", params.ProjectID, params.UserID)
		return nil, stacktrace.PropagateWithCode(err, stacktrace.GetCode(err), msg)
	}

	integration := &entities.ContentIntegration{
		ID:        uuid.New(),
		UserID:    params.UserID,
		ProjectID: params.ProjectID,
		Enabled:   true,
		Name:      params.Name,
		Title:     params.Title,
		Summary:   params.Summary,
		Text:      params.Text,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	if err := service.repository.Store(ctx, integration); err != nil {
		msg := fmt.Sprintf("could store content integration for user with ID [%s] and project [%s]", params.UserID, params.ProjectID)
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	service.dispatchIntegrationCreatedEvent(ctx, params.Source, integration.Integration())

	return integration, nil
}

// ContentIntegrationUpdateParams are the parameters for updating a content integration.
type ContentIntegrationUpdateParams struct {
	Name          string
	Summary       string
	Text          string
	Title         string
	Source        string
	IntegrationID uuid.UUID
	UserID        entities.UserID
}

// Update a new entities.WhatsappIntegration
func (service *ContentIntegrationService) Update(ctx context.Context, params *ContentIntegrationUpdateParams) (*entities.ContentIntegration, error) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	integration, err := service.repository.Load(ctx, params.UserID, params.IntegrationID)
	if err != nil {
		msg := fmt.Sprintf("cannot load integrtion [%s] for user ID [%s]", params.IntegrationID, params.UserID)
		return nil, stacktrace.PropagateWithCode(err, stacktrace.GetCode(err), msg)
	}

	integration.Name = params.Name
	integration.Text = params.Text
	integration.Summary = params.Summary
	integration.Title = params.Title

	if err = service.repository.Update(ctx, integration); err != nil {
		msg := fmt.Sprintf("could update content integration for user with ID [%s] and id [%s]", params.UserID, params.IntegrationID)
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	service.dispatchIntegrationUpdatedEvent(ctx, params.Source, integration.Integration())

	return integration, nil
}

// Delete a entities.ContentIntegration
func (service *ContentIntegrationService) Delete(ctx context.Context, params *IntegrationDeleteParams) error {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	integration, err := service.repository.Load(ctx, params.UserID, params.IntegrationID)
	if err != nil {
		msg := fmt.Sprintf("cannot load text integrtion [%s] for user ID [%s]", params.IntegrationID, params.UserID)
		return stacktrace.PropagateWithCode(err, stacktrace.GetCode(err), msg)
	}

	if err = service.repository.Delete(ctx, params.UserID, params.IntegrationID); err != nil {
		msg := fmt.Sprintf("cannot delete text integrtion [%s] for user ID [%s]", params.IntegrationID, params.UserID)
		return stacktrace.PropagateWithCode(err, stacktrace.GetCode(err), msg)
	}

	service.dispatchIntegrationDeletedEvent(ctx, params.Source, integration.Integration())
	return nil
}
