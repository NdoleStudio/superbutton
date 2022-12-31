package handlers

import (
	"fmt"

	"github.com/NdoleStudio/superbutton/pkg/services"
	"github.com/NdoleStudio/superbutton/pkg/telemetry"
	"github.com/NdoleStudio/superbutton/pkg/validators"
	"github.com/gofiber/fiber/v2"
	"github.com/palantir/stacktrace"
)

// UserHandler handles user http requests.
type UserHandler struct {
	handler
	logger    telemetry.Logger
	tracer    telemetry.Tracer
	validator *validators.UserHandlerValidator
	service   *services.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	validator *validators.UserHandlerValidator,
	service *services.UserService,
) (h *UserHandler) {
	return &UserHandler{
		logger:    logger.WithService(fmt.Sprintf("%T", h)),
		tracer:    tracer,
		validator: validator,
		service:   service,
	}
}

// RegisterRoutes registers the routes for the MessageHandler
func (h *UserHandler) RegisterRoutes(app *fiber.App, middlewares []fiber.Handler) {
	router := app.Group("/v1/users")
	router.Get("/me", h.computeRoute(middlewares, h.me)...)
	router.Get("/subscription-update-url", h.computeRoute(middlewares, h.subscriptionUpdateURL)...)
	router.Delete("/subscription", h.computeRoute(middlewares, h.cancelSubscription)...)
}

// me returns the currently authenticated entities.User
// @Summary      Currently authenticated user
// @Description  Fetches the currently authenticated user. This method creates the user if one doesn't exist
// @Security	 BearerAuth
// @Tags         Users
// @Produce      json
// @Success      200 		{object}	responses.Ok[entities.User]
// @Failure      400		{object}	responses.BadRequest
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure      422		{object}	responses.UnprocessableEntity
// @Failure      500		{object}	responses.InternalServerError
// @Router       /users/me 	[get]
func (h *UserHandler) me(c *fiber.Ctx) error {
	ctx, span := h.tracer.StartFromFiberCtx(c)
	defer span.End()

	ctxLogger := h.tracer.CtxLogger(h.logger, span)

	authUser := h.userFromContext(c)

	user, err := h.service.Get(ctx, c.OriginalURL(), authUser)
	if err != nil {
		msg := fmt.Sprintf("cannot get user with ID [%s]", authUser.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return h.responseInternalServerError(c)
	}

	return h.responseOK(c, "user fetched successfully", user)
}

// subscriptionUpdateURL returns the subscription update URL for the authenticated entities.User
// @Summary      Currently authenticated user subscription update URL
// @Description  Fetches the subscription URL of the authenticated user.
// @Security	 BearerAuth
// @Tags         Users
// @Produce      json
// @Success      200 		{object}	responses.Ok[string]
// @Failure      400		{object}	responses.BadRequest
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure      422		{object}	responses.UnprocessableEntity
// @Failure      500		{object}	responses.InternalServerError
// @Router       /users/subscription-update-url 	[get]
func (h *UserHandler) subscriptionUpdateURL(c *fiber.Ctx) error {
	ctx, span := h.tracer.StartFromFiberCtx(c)
	defer span.End()

	ctxLogger := h.tracer.CtxLogger(h.logger, span)
	authUser := h.userFromContext(c)

	url, err := h.service.GetSubscriptionUpdateURL(ctx, authUser.ID)
	if err != nil {
		msg := fmt.Sprintf("cannot get user with ID [%s]", authUser.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return h.responseInternalServerError(c)
	}

	return h.responseOK(c, "Subscription update URL fetched successfully", url)
}

// cancelSubscription cancels the subscription for the authenticated entities.User
// @Summary      Cancel the user's subscription
// @Description  Cancel the subscription of the authenticated user.
// @Security	 BearerAuth
// @Tags         Users
// @Produce      json
// @Success      200 		{object}	responses.NoContent
// @Failure      400		{object}	responses.BadRequest
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure      422		{object}	responses.UnprocessableEntity
// @Failure      500		{object}	responses.InternalServerError
// @Router       /users/subscription 	[delete]
func (h *UserHandler) cancelSubscription(c *fiber.Ctx) error {
	ctx, span := h.tracer.StartFromFiberCtx(c)
	defer span.End()

	ctxLogger := h.tracer.CtxLogger(h.logger, span)
	authUser := h.userFromContext(c)

	err := h.service.InitiateSubscriptionCancel(ctx, authUser.ID)
	if err != nil {
		msg := fmt.Sprintf("cannot get user with ID [%s]", authUser.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return h.responseInternalServerError(c)
	}

	return h.responseNoContent(c, "Subscription cancelled successfully")
}
