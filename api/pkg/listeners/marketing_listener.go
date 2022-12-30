package listeners

import (
	"context"
	"fmt"

	"github.com/NdoleStudio/superbutton/pkg/events"
	"github.com/NdoleStudio/superbutton/pkg/services"
	"github.com/NdoleStudio/superbutton/pkg/telemetry"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/davecgh/go-spew/spew"
	"github.com/palantir/stacktrace"
)

// MarketingListener listens for events with marketing handlers
type MarketingListener struct {
	tracer  telemetry.Tracer
	logger  telemetry.Logger
	service *services.MarketingService
}

// MarketingListeners returns the list of marketing listeners to events
func MarketingListeners(tracer telemetry.Tracer, logger telemetry.Logger, service *services.MarketingService) map[string]services.EventListener {
	listener := &MarketingListener{
		tracer:  tracer,
		logger:  logger.WithService(fmt.Sprintf("%T", &MarketingListener{})),
		service: service,
	}
	return map[string]services.EventListener{
		events.UserCreated: listener.OnUserCreated,
	}
}

// OnUserCreated handles the events.UserCreated event
func (listener *MarketingListener) OnUserCreated(ctx context.Context, event cloudevents.Event) error {
	ctx, span := listener.tracer.Start(ctx)
	defer span.End()

	var payload events.UserCreatedPayload
	if err := event.DataAs(&payload); err != nil {
		msg := fmt.Sprintf("cannot decode [%s] into [%T]", event.Data(), payload)
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	if err := listener.service.AddToList(ctx, payload.UserID); err != nil {
		msg := fmt.Sprintf("cannot add user with ID [%s] to list for event with ID [%s]", spew.Sdump(payload.UserID), event.ID())
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}
