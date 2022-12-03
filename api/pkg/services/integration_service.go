package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/NdoleStudio/superbutton/pkg/entities"
	"github.com/NdoleStudio/superbutton/pkg/events"
	"github.com/NdoleStudio/superbutton/pkg/telemetry"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/palantir/stacktrace"
)

type integrationService struct {
	service
	integrationType entities.IntegrationType
	tracer          telemetry.Tracer
	logger          telemetry.Logger
	eventDispatcher *EventDispatcher
}

// IntegrationDeleteParams are the parameters for updating a new whatsapp integration.
type IntegrationDeleteParams struct {
	Source        string
	IntegrationID uuid.UUID
	ProjectID     uuid.UUID
	UserID        entities.UserID
}

func (service *integrationService) dispatchIntegrationDeletedEvent(ctx context.Context, source string, integration *entities.Integration) {
	ctx, span, ctxLogger := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	event, err := service.createIntegrationDeletedEvent(source, integration)
	if err != nil {
		msg := fmt.Sprintf("cannot create [%s] event for integration [%s]", events.IntegrationDeleted, integration.IntegrationID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return
	}
	service.dispatchEvent(ctx, integration, event)
}

func (service *integrationService) dispatchIntegrationUpdatedEvent(ctx context.Context, source string, integration *entities.Integration) {
	ctx, span, ctxLogger := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	event, err := service.createIntegrationUpdatedEvent(source, integration)
	if err != nil {
		msg := fmt.Sprintf("cannot create [%s] event for integration [%s]", events.IntegrationUpdated, integration.IntegrationID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return
	}
	service.dispatchEvent(ctx, integration, event)
}

func (service *integrationService) dispatchIntegrationCreatedEvent(ctx context.Context, source string, integration *entities.Integration) {
	ctx, span, ctxLogger := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	event, err := service.createIntegrationCreatedEvent(source, integration)
	if err != nil {
		msg := fmt.Sprintf("cannot create [%s] event for integration [%s]", events.IntegrationCreated, integration.IntegrationID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return
	}
	service.dispatchEvent(ctx, integration, event)
}

func (service *integrationService) dispatchEvent(ctx context.Context, integration *entities.Integration, event *cloudevents.Event) {
	ctx, span, ctxLogger := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	if err := service.eventDispatcher.Dispatch(ctx, event); err != nil {
		msg := fmt.Sprintf("cannot dispatch [%s] event for integration [%s]", event.Type(), integration.IntegrationID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return
	}
}

func (service *integrationService) createIntegrationDeletedEvent(source string, integration *entities.Integration) (*cloudevents.Event, error) {
	event, err := service.createEvent(events.IntegrationDeleted, source, &events.IntegrationDeletedPayload{
		UserID:               integration.UserID,
		ProjectID:            integration.ProjectID,
		IntegrationID:        integration.IntegrationID,
		IntegrationType:      integration.Type,
		IntegrationName:      integration.Name,
		IntegrationDeletedAt: time.Now().UTC(),
	})
	if err != nil {
		msg := fmt.Sprintf("cannot create [%s] event for [%s] integration [%s]", events.IntegrationDeleted, integration.Type, integration.IntegrationID)
		return nil, stacktrace.Propagate(err, msg)
	}

	return event, nil
}

func (service *integrationService) createIntegrationUpdatedEvent(source string, integration *entities.Integration) (*cloudevents.Event, error) {
	event, err := service.createEvent(events.IntegrationUpdated, source, &events.IntegrationUpdatedPayload{
		UserID:               integration.UserID,
		ProjectID:            integration.ProjectID,
		IntegrationID:        integration.IntegrationID,
		IntegrationType:      integration.Type,
		IntegrationName:      integration.Name,
		IntegrationUpdatedAt: integration.UpdatedAt,
	})
	if err != nil {
		msg := fmt.Sprintf("cannot create [%s] event for [%s] integration [%s]", events.IntegrationUpdated, integration.Type, integration.IntegrationID)
		return nil, stacktrace.Propagate(err, msg)
	}

	return event, nil
}

func (service *integrationService) createIntegrationCreatedEvent(source string, integration *entities.Integration) (*cloudevents.Event, error) {
	event, err := service.createEvent(events.IntegrationCreated, source, &events.IntegrationCreatedPayload{
		UserID:               integration.UserID,
		ProjectID:            integration.ProjectID,
		IntegrationID:        integration.IntegrationID,
		IntegrationType:      integration.Type,
		IntegrationName:      integration.Name,
		IntegrationCreatedAt: integration.CreatedAt,
	})
	if err != nil {
		msg := fmt.Sprintf("cannot create [%s] event for [%s] integration [%s]", events.IntegrationCreated, integration.Type, integration.IntegrationID)
		return nil, stacktrace.Propagate(err, msg)
	}

	return event, nil
}
