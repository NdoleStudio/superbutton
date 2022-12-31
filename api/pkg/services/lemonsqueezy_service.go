package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/NdoleStudio/superbutton/pkg/repositories"

	lemonsqueezy "github.com/NdoleStudio/lemonsqueezy-go"
	"github.com/NdoleStudio/superbutton/pkg/entities"
	"github.com/NdoleStudio/superbutton/pkg/events"
	"github.com/NdoleStudio/superbutton/pkg/telemetry"
	"github.com/palantir/stacktrace"
)

// LemonsqueezyService is responsible for managing lemonsqueezy events
type LemonsqueezyService struct {
	service
	logger          telemetry.Logger
	tracer          telemetry.Tracer
	eventDispatcher *EventDispatcher
	userRepository  repositories.UserRepository
}

// NewLemonsqueezyService creates a new LemonsqueezyService
func NewLemonsqueezyService(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	repository repositories.UserRepository,
	eventDispatcher *EventDispatcher,
) (s *LemonsqueezyService) {
	return &LemonsqueezyService{
		logger:          logger.WithService(fmt.Sprintf("%T", s)),
		tracer:          tracer,
		userRepository:  repository,
		eventDispatcher: eventDispatcher,
	}
}

// HandleSubscriptionCreatedEvent handles the subscription_created lemonsqueezy event
func (service *LemonsqueezyService) HandleSubscriptionCreatedEvent(ctx context.Context, source string, request *lemonsqueezy.WebHookRequestSubscription) error {
	ctx, span, ctxLogger := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	payload := &events.UserSubscriptionCreatedPayload{
		UserID:                entities.UserID(request.Meta.CustomData["user_id"].(string)),
		SubscriptionCreatedAt: request.Data.Attributes.CreatedAt,
		SubscriptionID:        request.Data.ID,
		SubscriptionName:      service.subscriptionName(request.Data.Attributes.VariantName),
		SubscriptionRenewsAt:  request.Data.Attributes.RenewsAt,
		SubscriptionStatus:    request.Data.Attributes.Status,
	}

	event, err := service.createEvent(events.UserSubscriptionCreated, source, payload)
	if err != nil {
		msg := fmt.Sprintf("cannot created [%s] event for user [%s]", events.UserSubscriptionCreated, payload.UserID)
		return stacktrace.Propagate(err, msg)
	}

	if err = service.eventDispatcher.Dispatch(ctx, event); err != nil {
		msg := fmt.Sprintf("cannot dispatch [%s] event for user [%s]", events.UserSubscriptionCreated, payload.UserID)
		return stacktrace.Propagate(err, msg)
	}
	ctxLogger.Info(fmt.Sprintf("[%s] subscription [%s] created for user [%s]", payload.SubscriptionName, payload.SubscriptionID, payload.UserID))
	return nil
}

// HandleSubscriptionCanceledEvent handles the subscription_cancelled lemonsqueezy event
func (service *LemonsqueezyService) HandleSubscriptionCanceledEvent(ctx context.Context, source string, request *lemonsqueezy.WebHookRequestSubscription) error {
	ctx, span, ctxLogger := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	user, err := service.userRepository.LoadBySubscriptionID(ctx, request.Data.ID)
	if err != nil {
		msg := fmt.Sprintf("cannot load user with subscription ID [%s]", request.Data.ID)
		return stacktrace.Propagate(err, msg)
	}

	payload := &events.UserSubscriptionCancelledPayload{
		UserID:                  user.ID,
		SubscriptionCancelledAt: request.Data.Attributes.CreatedAt,
		SubscriptionID:          request.Data.ID,
		SubscriptionName:        service.subscriptionName(request.Data.Attributes.VariantName),
		SubscriptionEndsAt:      *request.Data.Attributes.EndsAt,
		SubscriptionStatus:      request.Data.Attributes.Status,
	}

	event, err := service.createEvent(events.UserSubscriptionCancelled, source, payload)
	if err != nil {
		msg := fmt.Sprintf("cannot created [%s] event for user [%s]", events.UserSubscriptionCancelled, payload.UserID)
		return stacktrace.Propagate(err, msg)
	}

	if err = service.eventDispatcher.Dispatch(ctx, event); err != nil {
		msg := fmt.Sprintf("cannot dispatch [%s] event for user [%s]", events.UserSubscriptionCancelled, payload.UserID)
		return stacktrace.Propagate(err, msg)
	}
	ctxLogger.Info(fmt.Sprintf("[%s] subscription [%s] cancelled for user [%s]", payload.SubscriptionName, payload.SubscriptionID, payload.UserID))
	return nil
}

func (service *LemonsqueezyService) subscriptionName(variant string) entities.SubscriptionName {
	if strings.Contains(strings.ToLower(variant), "monthly") {
		return entities.SubscriptionNameProMonthly
	}
	return entities.SubscriptionNameProYearly
}
