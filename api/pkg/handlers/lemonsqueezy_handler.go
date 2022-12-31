package handlers

import (
	"fmt"

	"github.com/NdoleStudio/superbutton/pkg/telemetry"
	"github.com/davecgh/go-spew/spew"
	"github.com/gofiber/fiber/v2"
)

// LemonsqueezyHandler handles lemonsqueezy events
type LemonsqueezyHandler struct {
	handler
	logger telemetry.Logger
	tracer telemetry.Tracer
}

// NewLemonsqueezyHandlerHandler creates a new LemonsqueezyHandler
func NewLemonsqueezyHandlerHandler(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
) (h *LemonsqueezyHandler) {
	return &LemonsqueezyHandler{
		logger: logger.WithService(fmt.Sprintf("%T", h)),
		tracer: tracer,
	}
}

// RegisterRoutes registers the routes for the MessageHandler
func (h *LemonsqueezyHandler) RegisterRoutes(app *fiber.App, middlewares ...fiber.Handler) {
	router := app.Group("/v1/lemonsqueezy")
	router.Post("/event", h.computeRoute(middlewares, h.Event)...)
}

// Event consumes a lemonsqueezy event
// @Summary      Consume a lemonsqueezy event
// @Description  Publish a lemonsqueezy event to the registered listeners
// @Security	 BearerAuth
// @Tags         Lemonsqueezy
// @Accept       json
// @Produce      json
// @Param        payload	body 		requests.CloudEvent				true 	"cloud event payload"
// @Success      204 		{object}	responses.NoContent
// @Failure      400		{object}	responses.BadRequest
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure      422		{object}	responses.UnprocessableEntity
// @Failure      500		{object}	responses.InternalServerError
// @Router       /lemonsqueezy/event [post]
func (h *LemonsqueezyHandler) Event(c *fiber.Ctx) error {
	_, span := h.tracer.StartFromFiberCtx(c)
	defer span.End()

	spew.Dump(string(c.Body()))
	spew.Dump(c.GetReqHeaders())

	return h.responseNoContent(c, "event consumed successfully")
}
