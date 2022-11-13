package handlers

import (
	"fmt"

	cloudevents "github.com/cloudevents/sdk-go/v2"

	"github.com/NdoleStudio/superbutton/pkg/services"
	"github.com/NdoleStudio/superbutton/pkg/telemetry"
	"github.com/davecgh/go-spew/spew"
	"github.com/gofiber/fiber/v2"
	"github.com/palantir/stacktrace"
)

// EventsHandler handles heartbeat http requests.
type EventsHandler struct {
	handler
	logger  telemetry.Logger
	tracer  telemetry.Tracer
	service *services.EventDispatcher
}

// NewEventsHandler creates a new EventsHandler
func NewEventsHandler(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	service *services.EventDispatcher,
) (h *EventsHandler) {
	return &EventsHandler{
		logger:  logger.WithService(fmt.Sprintf("%T", h)),
		tracer:  tracer,
		service: service,
	}
}

// RegisterRoutes registers the routes for the MessageHandler
func (h *EventsHandler) RegisterRoutes(router fiber.Router) {
	router.Post("/events/consume", h.Consume)
}

// Consume a cloudevents.Event
// @Summary      Consume a cloud event
// @Description  Publish a cloud event to the registered listeners
// @Security	 BearerAuth
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        payload	body 		requests.CloudEvent				true 	"cloud event payload"
// @Success      204 		{object}	responses.NoContent
// @Failure      400		{object}	responses.BadRequest
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure      422		{object}	responses.UnprocessableEntity
// @Failure      500		{object}	responses.InternalServerError
// @Router       /events/consume [post]
func (h *EventsHandler) Consume(c *fiber.Ctx) error {
	ctx, span := h.tracer.StartFromFiberCtx(c)
	defer span.End()

	ctxLogger := h.tracer.CtxLogger(h.logger, span)

	var request cloudevents.Event
	if err := c.BodyParser(&request); err != nil {
		msg := fmt.Sprintf("cannot marshall params [%s] into %T", c.OriginalURL(), request)
		ctxLogger.Warn(stacktrace.Propagate(err, msg))
		return h.responseBadRequest(c, err)
	}

	if err := request.Validate(); err != nil {
		msg := fmt.Sprintf("validation errors [%s], while dispatching event [%+#v]", spew.Sdump(err.Error()), request)
		ctxLogger.Warn(stacktrace.NewError(msg))
		return h.responseUnprocessableEntity(c, map[string][]string{"event": {err.Error()}}, "validation errors while consuming event")
	}

	h.service.Publish(ctx, request)

	return h.responseNoContent(c, "event consumed successfully")
}
