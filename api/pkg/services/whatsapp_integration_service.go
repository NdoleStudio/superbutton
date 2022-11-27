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

// WhatsappIntegrationService is responsible for managing entities.WhatsappIntegration
type WhatsappIntegrationService struct {
	integrationService
	projectRepository repositories.ProjectRepository
	repository        repositories.WhatsappIntegrationRepository
}

// NewWhatsappIntegrationService creates a new WhatsappIntegrationService
func NewWhatsappIntegrationService(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	eventDispatcher *EventDispatcher,
	repository repositories.WhatsappIntegrationRepository,
) (s *WhatsappIntegrationService) {
	return &WhatsappIntegrationService{
		repository: repository,
		integrationService: integrationService{
			tracer:          tracer,
			logger:          logger.WithService(fmt.Sprintf("%T", s)),
			eventDispatcher: eventDispatcher,
		},
	}
}

// Index fetches all entities.Project for an authenticated user
func (service *WhatsappIntegrationService) Index(ctx context.Context, userID entities.UserID, projectID uuid.UUID) ([]*entities.WhatsappIntegration, error) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	integrations, err := service.repository.Fetch(ctx, userID, projectID)
	if err != nil {
		msg := fmt.Sprintf("could fetch whatsapp integrations for user with ID [%s] and projectID [%s]", userID, projectID)
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return integrations, nil
}

// WhatsappIntegrationCreateParams are the parameters for creating a new whatsapp integration.
type WhatsappIntegrationCreateParams struct {
	Name        string
	Text        string
	PhoneNumber string
	Source      string
	ProjectID   uuid.UUID
	UserID      entities.UserID
}

// Create a new entities.WhatsappIntegration
func (service *WhatsappIntegrationService) Create(ctx context.Context, params *WhatsappIntegrationCreateParams) (*entities.WhatsappIntegration, error) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	if _, err := service.projectRepository.Load(ctx, params.UserID, params.ProjectID); err != nil {
		msg := fmt.Sprintf("cannot load project [%s] for user ID [%s]", params.ProjectID, params.UserID)
		return nil, stacktrace.PropagateWithCode(err, stacktrace.GetCode(err), msg)
	}

	integration := &entities.WhatsappIntegration{
		ID:          uuid.New(),
		UserID:      params.UserID,
		ProjectID:   params.ProjectID,
		Text:        params.Text,
		PhoneNumber: params.PhoneNumber,
		Name:        params.Name,
		Icon:        "whatsapp",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if err := service.repository.Store(ctx, integration); err != nil {
		msg := fmt.Sprintf("could store whatsapp integration for user with ID [%s] and project [%s]", params.UserID, params.ProjectID)
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	service.dispatchIntegrationCreatedEvent(ctx, params.Source, integration.Integration())

	return integration, nil
}

// WhatsappIntegrationUpdateParams are the parameters for updating a new whatsapp integration.
type WhatsappIntegrationUpdateParams struct {
	Name          string
	Text          string
	PhoneNumber   string
	Source        string
	IntegrationID uuid.UUID
	UserID        entities.UserID
}

// Update a new entities.WhatsappIntegration
func (service *WhatsappIntegrationService) Update(ctx context.Context, params *WhatsappIntegrationUpdateParams) (*entities.WhatsappIntegration, error) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	integration, err := service.repository.Load(ctx, params.UserID, params.IntegrationID)
	if err != nil {
		msg := fmt.Sprintf("cannot load integrtion [%s] for user ID [%s]", params.IntegrationID, params.UserID)
		return nil, stacktrace.PropagateWithCode(err, stacktrace.GetCode(err), msg)
	}

	integration.Name = params.Name
	integration.Text = params.Text
	integration.PhoneNumber = params.PhoneNumber

	if err = service.repository.Update(ctx, integration); err != nil {
		msg := fmt.Sprintf("could update whatsapp integration for user with ID [%s] and id [%s]", params.UserID, params.IntegrationID)
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	service.dispatchIntegrationUpdatedEvent(ctx, params.Source, integration.Integration())

	return integration, nil
}

// WhatsappIntegrationDeleteParams are the parameters for updating a new whatsapp integration.
type WhatsappIntegrationDeleteParams struct {
	Source        string
	IntegrationID uuid.UUID
	ProjectID     uuid.UUID
	UserID        entities.UserID
}

// Delete a entities.WhatsappIntegration
func (service *WhatsappIntegrationService) Delete(ctx context.Context, params *WhatsappIntegrationDeleteParams) error {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	integration, err := service.repository.Load(ctx, params.UserID, params.IntegrationID)
	if err != nil {
		msg := fmt.Sprintf("cannot load integrtion [%s] for user ID [%s]", params.IntegrationID, params.UserID)
		return stacktrace.PropagateWithCode(err, stacktrace.GetCode(err), msg)
	}

	if err = service.repository.Delete(ctx, params.UserID, params.IntegrationID); err != nil {
		msg := fmt.Sprintf("cannot delete integrtion [%s] for user ID [%s]", params.IntegrationID, params.UserID)
		return stacktrace.PropagateWithCode(err, stacktrace.GetCode(err), msg)
	}

	service.dispatchIntegrationDeletedEvent(ctx, params.Source, integration.Integration())
	return nil
}
