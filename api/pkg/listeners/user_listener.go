package listeners

import (
	"context"
	"fmt"

	"github.com/NdoleStudio/superbutton/pkg/events"
	"github.com/NdoleStudio/superbutton/pkg/services"
	"github.com/NdoleStudio/superbutton/pkg/telemetry"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/palantir/stacktrace"
)

// UserListener listens for events with user handlers
type UserListener struct {
	tracer  telemetry.Tracer
	logger  telemetry.Logger
	service *services.UserService
}

// UserListeners returns the list of user listeners to events
func UserListeners(tracer telemetry.Tracer, logger telemetry.Logger, service *services.UserService) map[string]services.EventListener {
	listener := &UserListener{
		tracer:  tracer,
		logger:  logger.WithService(fmt.Sprintf("%T", &MarketingListener{})),
		service: service,
	}
	return map[string]services.EventListener{
		events.UserSubscriptionCreated:   listener.OnUserSubscriptionCreated,
		events.UserSubscriptionCancelled: listener.OnUserSubscriptionCancelled,
	}
}

// OnUserSubscriptionCreated handles the events.UserSubscriptionCreated event
func (listener *UserListener) OnUserSubscriptionCreated(ctx context.Context, event cloudevents.Event) error {
	ctx, span := listener.tracer.Start(ctx)
	defer span.End()

	var payload events.UserSubscriptionCreatedPayload
	if err := event.DataAs(&payload); err != nil {
		msg := fmt.Sprintf("cannot decode [%s] into [%T]", event.Data(), payload)
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	if err := listener.service.StartSubscription(ctx, event.Source(), &payload); err != nil {
		msg := fmt.Sprintf("cannot start subscription for user with ID [%s] for event with ID [%s]", payload.UserID, event.ID())
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

// OnUserSubscriptionCancelled handles the events.UserSubscriptionCancelled event
func (listener *UserListener) OnUserSubscriptionCancelled(ctx context.Context, event cloudevents.Event) error {
	ctx, span := listener.tracer.Start(ctx)
	defer span.End()

	var payload events.UserSubscriptionCancelledPayload
	if err := event.DataAs(&payload); err != nil {
		msg := fmt.Sprintf("cannot decode [%s] into [%T]", event.Data(), payload)
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	if err := listener.service.CancelSubscription(ctx, event.Source(), &payload); err != nil {
		msg := fmt.Sprintf("cannot cancell subscription for user with ID [%s] for event with ID [%s]", payload.UserID, event.ID())
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}
