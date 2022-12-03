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

// LinkIntegrationService is responsible for managing entities.LinkIntegration
type LinkIntegrationService struct {
	integrationService
	projectRepository repositories.ProjectRepository
	repository        repositories.LinkIntegrationRepository
}

// NewLinkIntegrationService creates a new LinkIntegrationService
func NewLinkIntegrationService(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	eventDispatcher *EventDispatcher,
	projectRepository repositories.ProjectRepository,
	repository repositories.LinkIntegrationRepository,
) (s *LinkIntegrationService) {
	return &LinkIntegrationService{
		repository:        repository,
		projectRepository: projectRepository,
		integrationService: integrationService{
			tracer:          tracer,
			integrationType: entities.IntegrationTypeLink,
			logger:          logger.WithService(fmt.Sprintf("%T", s)),
			eventDispatcher: eventDispatcher,
		},
	}
}

// Get returns an entities.WhatsappIntegration for an authenticated user
func (service *LinkIntegrationService) Get(ctx context.Context, userID entities.UserID, integrationID uuid.UUID) (*entities.LinkIntegration, error) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	integration, err := service.repository.Load(ctx, userID, integrationID)
	if err != nil {
		msg := fmt.Sprintf("could not get [%s] integrations for user with ID [%s] and ID [%s]", service.integrationType, userID, integrationID)
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return integration, nil
}

// Index fetches all entities.Project for an authenticated user
func (service *LinkIntegrationService) Index(ctx context.Context, userID entities.UserID, projectID uuid.UUID) ([]*entities.LinkIntegration, error) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	integrations, err := service.repository.Fetch(ctx, userID, projectID)
	if err != nil {
		msg := fmt.Sprintf("could fetch [%s] integrations for user with ID [%s] and projectID [%s]", service.integrationType, userID, projectID)
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return integrations, nil
}

// LinkIntegrationCreateParams are the parameters for creating a new entities.LinkIntegration.
type LinkIntegrationCreateParams struct {
	Name        string
	PhoneNumber string
	Text        string
	Source      string
	URL         string
	ProjectID   uuid.UUID
	UserID      entities.UserID
}

// Create a new entities.PhoneCallIntegration
func (service *LinkIntegrationService) Create(ctx context.Context, params *LinkIntegrationCreateParams) (*entities.LinkIntegration, error) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	if _, err := service.projectRepository.Load(ctx, params.UserID, params.ProjectID); err != nil {
		msg := fmt.Sprintf("cannot load project [%s] for user ID [%s]", params.ProjectID, params.UserID)
		return nil, stacktrace.PropagateWithCode(err, stacktrace.GetCode(err), msg)
	}

	integration := &entities.LinkIntegration{
		ID:        uuid.New(),
		UserID:    params.UserID,
		ProjectID: params.ProjectID,
		Enabled:   true,
		Icon:      string(service.integrationType),
		Name:      params.Name,
		URL:       params.URL,
		Text:      params.Text,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	if err := service.repository.Store(ctx, integration); err != nil {
		msg := fmt.Sprintf("could store [%s] integration for user with ID [%s] and project [%s]", service.integrationType, params.UserID, params.ProjectID)
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	service.dispatchIntegrationCreatedEvent(ctx, params.Source, integration.Integration())

	return integration, nil
}

// LinkIntegrationUpdateParams are the parameters for updating an entities.LinkIntegration
type LinkIntegrationUpdateParams struct {
	Name          string
	Text          string
	URL           string
	Source        string
	IntegrationID uuid.UUID
	UserID        entities.UserID
}

// Update a new entities.WhatsappIntegration
func (service *LinkIntegrationService) Update(ctx context.Context, params *LinkIntegrationUpdateParams) (*entities.LinkIntegration, error) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	integration, err := service.repository.Load(ctx, params.UserID, params.IntegrationID)
	if err != nil {
		msg := fmt.Sprintf("cannot load integrtion [%s] for user ID [%s]", params.IntegrationID, params.UserID)
		return nil, stacktrace.PropagateWithCode(err, stacktrace.GetCode(err), msg)
	}

	integration.Name = params.Name
	integration.Text = params.Text
	integration.URL = params.URL

	if err = service.repository.Update(ctx, integration); err != nil {
		msg := fmt.Sprintf("could update [%s] integration for user with ID [%s] and id [%s]", service.integrationType, params.UserID, params.IntegrationID)
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	service.dispatchIntegrationUpdatedEvent(ctx, params.Source, integration.Integration())

	return integration, nil
}

// Delete a entities.ContentIntegration
func (service *LinkIntegrationService) Delete(ctx context.Context, params *IntegrationDeleteParams) error {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	integration, err := service.repository.Load(ctx, params.UserID, params.IntegrationID)
	if err != nil {
		msg := fmt.Sprintf("cannot load [%s] integrtion [%s] for user ID [%s]", service.integrationType, params.IntegrationID, params.UserID)
		return stacktrace.PropagateWithCode(err, stacktrace.GetCode(err), msg)
	}

	if err = service.repository.Delete(ctx, params.UserID, params.IntegrationID); err != nil {
		msg := fmt.Sprintf("cannot delete  [%s] integrtion [%s] for user ID [%s]", service.integrationType, params.IntegrationID, params.UserID)
		return stacktrace.PropagateWithCode(err, stacktrace.GetCode(err), msg)
	}

	service.dispatchIntegrationDeletedEvent(ctx, params.Source, integration.Integration())
	return nil
}
