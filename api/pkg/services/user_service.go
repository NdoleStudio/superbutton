package services

import (
	"context"
	"fmt"

	lemonsqueezy "github.com/NdoleStudio/lemonsqueezy-go"

	"github.com/NdoleStudio/superbutton/pkg/entities"
	"github.com/NdoleStudio/superbutton/pkg/events"
	"github.com/NdoleStudio/superbutton/pkg/repositories"
	"github.com/NdoleStudio/superbutton/pkg/telemetry"
	"github.com/palantir/stacktrace"
)

// UserService is responsible for managing entities.User
type UserService struct {
	service
	logger             telemetry.Logger
	tracer             telemetry.Tracer
	repository         repositories.UserRepository
	lemonsqueezyClient *lemonsqueezy.Client
	eventDispatcher    *EventDispatcher
}

// NewUserService creates a new UserService
func NewUserService(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	eventDispatcher *EventDispatcher,
	repository repositories.UserRepository,
	lemonsqueezyClient *lemonsqueezy.Client,
) (s *UserService) {
	return &UserService{
		logger:             logger.WithService(fmt.Sprintf("%T", s)),
		tracer:             tracer,
		lemonsqueezyClient: lemonsqueezyClient,
		eventDispatcher:    eventDispatcher,
		repository:         repository,
	}
}

// Get fetches or creates an entities.User
func (service *UserService) Get(ctx context.Context, source string, authUser entities.AuthUser) (*entities.User, error) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	user, created, err := service.repository.LoadOrStore(ctx, authUser)
	if err != nil {
		msg := fmt.Sprintf("could not get [%T] with from [%+#v]", user, authUser)
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	if created {
		service.dispatchUserCreatedEvent(ctx, source, user)
	}

	return user, nil
}

// StartSubscription starts a subscription for an entities.User
func (service *UserService) StartSubscription(ctx context.Context, source string, params *events.UserSubscriptionCreatedPayload) error {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	user, err := service.repository.Load(ctx, params.UserID)
	if err != nil {
		msg := fmt.Sprintf("could not get [%T] with with ID [%s]", user, params.UserID)
		return service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	user.SubscriptionID = params.SubscriptionID
	user.SubscriptionName = params.SubscriptionName
	user.SubscriptionRenewsAt = &params.SubscriptionRenewsAt
	user.SubscriptionStatus = params.SubscriptionStatus
	user.SubscriptionEndsAt = nil

	if err = service.repository.Update(ctx, user); err != nil {
		msg := fmt.Sprintf("could not update [%T] with with ID [%s] after update", user, params.UserID)
		return service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	service.dispatchUserUpdatedEvent(ctx, source, user)
	return nil
}

// InitiateSubscriptionCancel initiates the cancelling of a subscription on lemonsqueezy
func (service *UserService) InitiateSubscriptionCancel(ctx context.Context, userID entities.UserID) error {
	ctx, span, ctxLogger := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	user, err := service.repository.Load(ctx, userID)
	if err != nil {
		msg := fmt.Sprintf("could not get [%T] with with ID [%s]", user, userID)
		return service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	if _, _, err = service.lemonsqueezyClient.Subscriptions.Cancel(ctx, user.SubscriptionID); err != nil {
		msg := fmt.Sprintf("could not cancel subscription [%s] for [%T] with with ID [%s]", user.SubscriptionID, user, user.ID)
		return service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	ctxLogger.Info(fmt.Sprintf("cancelled subscription [%s] for user [%s]", user.SubscriptionID, user.ID))
	return nil
}

// GetSubscriptionUpdateURL initiates the cancelling of a subscription on lemonsqueezy
func (service *UserService) GetSubscriptionUpdateURL(ctx context.Context, userID entities.UserID) (url string, err error) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	user, err := service.repository.Load(ctx, userID)
	if err != nil {
		msg := fmt.Sprintf("could not get [%T] with with ID [%s]", user, userID)
		return "", service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	subscription, _, err := service.lemonsqueezyClient.Subscriptions.Get(ctx, user.SubscriptionID)
	if err != nil {
		msg := fmt.Sprintf("could not get subscription [%s] for [%T] with with ID [%s]", user.SubscriptionID, user, user.ID)
		return url, service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return subscription.Data.Attributes.Urls.UpdatePaymentMethod, nil
}

// CancelSubscription starts a subscription for an entities.User
func (service *UserService) CancelSubscription(ctx context.Context, source string, params *events.UserSubscriptionCancelledPayload) error {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	user, err := service.repository.Load(ctx, params.UserID)
	if err != nil {
		msg := fmt.Sprintf("could not get [%T] with with ID [%s]", user, params.UserID)
		return service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	user.SubscriptionID = params.SubscriptionID
	user.SubscriptionName = params.SubscriptionName
	user.SubscriptionRenewsAt = nil
	user.SubscriptionStatus = params.SubscriptionStatus
	user.SubscriptionEndsAt = &params.SubscriptionEndsAt

	if err = service.repository.Update(ctx, user); err != nil {
		msg := fmt.Sprintf("could not update [%T] with with ID [%s] after update", user, params.UserID)
		return service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	service.dispatchUserUpdatedEvent(ctx, source, user)
	return nil
}

func (service *UserService) dispatchUserUpdatedEvent(ctx context.Context, source string, user *entities.User) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	ctxLogger := service.tracer.CtxLogger(service.logger, span)

	event, err := service.createEvent(events.UserUpdated, source, &events.UserUpdatedPayload{
		UserID:           user.ID,
		UserUpdatedAt:    user.UpdatedAt,
		SubscriptionName: user.SubscriptionName,
		UserEmail:        user.Email,
	})
	if err != nil {
		msg := fmt.Sprintf("cannot created [%s] event for user [%s]", events.UserUpdated, user.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return
	}

	if err = service.eventDispatcher.Dispatch(ctx, event); err != nil {
		msg := fmt.Sprintf("cannot dispatch [%s] event for user [%s]", events.UserUpdated, user.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return
	}
}

func (service *UserService) dispatchUserCreatedEvent(ctx context.Context, source string, user *entities.User) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	ctxLogger := service.tracer.CtxLogger(service.logger, span)

	event, err := service.createEvent(events.UserCreated, source, &events.UserCreatedPayload{
		UserID:        user.ID,
		UserCreatedAt: user.CreatedAt,
		UserName:      user.Name,
		UserEmail:     user.Email,
	})
	if err != nil {
		msg := fmt.Sprintf("cannot created [%s] event for user [%s]", events.UserCreated, user.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return
	}

	if err = service.eventDispatcher.Dispatch(ctx, event); err != nil {
		msg := fmt.Sprintf("cannot dispatch [%s] event for user [%s]", events.UserCreated, user.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return
	}
}
