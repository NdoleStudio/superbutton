package handlers

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/NdoleStudio/superbutton/pkg/repositories"
	"github.com/NdoleStudio/superbutton/pkg/requests"
	"github.com/davecgh/go-spew/spew"

	"github.com/NdoleStudio/superbutton/pkg/services"
	"github.com/NdoleStudio/superbutton/pkg/telemetry"
	"github.com/NdoleStudio/superbutton/pkg/validators"
	"github.com/gofiber/fiber/v2"
	"github.com/palantir/stacktrace"
)

// PhoneCallIntegrationHandler handles user http requests.
type PhoneCallIntegrationHandler struct {
	handler
	logger    telemetry.Logger
	tracer    telemetry.Tracer
	validator *validators.PhoneCallIntegrationHandlerValidator
	service   *services.PhoneCallIntegrationService
}

// NewPhoneCallIntegrationHandler creates a new PhoneCallIntegrationHandler
func NewPhoneCallIntegrationHandler(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	validator *validators.PhoneCallIntegrationHandlerValidator,
	service *services.PhoneCallIntegrationService,
) (h *PhoneCallIntegrationHandler) {
	return &PhoneCallIntegrationHandler{
		logger:    logger.WithService(fmt.Sprintf("%T", h)),
		tracer:    tracer,
		validator: validator,
		service:   service,
	}
}

// RegisterRoutes registers the routes for the ContentIntegrationHandler
func (h *PhoneCallIntegrationHandler) RegisterRoutes(app *fiber.App, middlewares []fiber.Handler) {
	router := app.Group("/v1/projects/:projectID/phone-call-integrations")
	router.Post("/", h.computeRoute(middlewares, h.create)...)
	router.Get("/:integrationID", h.computeRoute(middlewares, h.show)...)
	router.Put("/:integrationID", h.computeRoute(middlewares, h.update)...)
	router.Delete("/:integrationID", h.computeRoute(middlewares, h.delete)...)
}

// @Summary      Get phone call integration
// @Description  Fetches a specific phone call integration
// @Security	 BearerAuth
// @Tags         PhoneCallIntegration
// @Produce      json
// @Success      200 		{object}	responses.Ok[entities.PhoneCallIntegration]
// @Failure      400		{object}	responses.BadRequest
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure 	 404    	{object}	responses.NotFound
// @Failure      422		{object}	responses.UnprocessableEntity
// @Failure      500		{object}	responses.InternalServerError
// @Router       /projects/:projectID/phone-call-integrations/:integrationID 	[get]
func (h *PhoneCallIntegrationHandler) show(c *fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	if errors := h.mergeErrors(h.validateUUID(c, "integrationID")); len(errors) != 0 {
		msg := fmt.Sprintf("validation errors [%s], while fetching [phone-call] integration with ID [%s]", spew.Sdump(errors), c.Params("integrationID"))
		ctxLogger.Warn(stacktrace.NewError(msg))
		return h.responseUnprocessableEntity(c, errors, "validation errors while fetching phone call integration")
	}

	integrationID := uuid.MustParse(c.Params("integrationID"))
	authUser := h.userFromContext(c)

	integration, err := h.service.Get(ctx, authUser.ID, integrationID)
	if err != nil {
		msg := fmt.Sprintf("cannot fetch [phone-call] intergration [%s] for user with ID [%s]", authUser.ID, integrationID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return h.responseInternalServerError(c)
	}

	return h.responseOK(c, "phone call integration fetched successfully", integration)
}

// @Summary      Create a phone call integration
// @Description  This endpoint creates a new phone call integration for a project
// @Security	 BearerAuth
// @Tags         PhoneCallIntegration
// @Produce      json
// @Param        payload	body 		requests.PhoneCallIntegrationCreateRequest	true 	"phone call integration create payload"
// @Success      200 		{object}	responses.Ok[entities.PhoneCallIntegration]
// @Failure      400		{object}	responses.BadRequest
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure      422		{object}	responses.UnprocessableEntity
// @Failure      500		{object}	responses.InternalServerError
// @Router       /projects/{projectID}/phone-call-integrations [post]
func (h *PhoneCallIntegrationHandler) create(c *fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	var request requests.PhoneCallIntegrationCreateRequest
	if err := c.BodyParser(&request); err != nil {
		msg := fmt.Sprintf("cannot marshall params [%s] into %T", c.OriginalURL(), request)
		ctxLogger.Warn(stacktrace.Propagate(err, msg))
		return h.responseBadRequest(c, err)
	}
	request.ProjectID = c.Params("projectID")

	if errors := h.validator.ValidateCreate(ctx, request.Sanitize()); len(errors) != 0 {
		msg := fmt.Sprintf("validation errors [%s], while creating [phone-call] integration with request [%s]", spew.Sdump(errors), c.Body())
		ctxLogger.Warn(stacktrace.NewError(msg))
		return h.responseUnprocessableEntity(c, errors, "validation errors while creating integration")
	}

	authUser := h.userFromContext(c)
	integration, err := h.service.Create(ctx, request.ToCreateParams(c.OriginalURL(), authUser.ID))
	if err != nil {
		msg := fmt.Sprintf("cannot create [phone-call] integration for user with ID [%s]", authUser.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return h.responseInternalServerError(c)
	}

	return h.responseOK(c, "phone call integration created successfully", integration)
}

// @Summary      Update a phone call integration
// @Description  This endpoint updates a phone call integration for a project
// @Security	 BearerAuth
// @Tags         PhoneCallIntegration
// @Produce      json
// @Param        payload	body 		requests.PhoneCallIntegrationUpdateRequest	true 	"phone call integration update payload"
// @Success      200 		{object}	responses.Ok[entities.PhoneCallIntegration]
// @Failure      400		{object}	responses.BadRequest
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure      422		{object}	responses.UnprocessableEntity
// @Failure      500		{object}	responses.InternalServerError
// @Router       /projects/{projectID}/phone-call-integrations/{integrationID} [put]
func (h *PhoneCallIntegrationHandler) update(c *fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	var request requests.PhoneCallIntegrationUpdateRequest
	if err := c.BodyParser(&request); err != nil {
		msg := fmt.Sprintf("cannot marshall params [%s] into %T", c.OriginalURL(), request)
		ctxLogger.Warn(stacktrace.Propagate(err, msg))
		return h.responseBadRequest(c, err)
	}
	request.IntegrationID = c.Params("integrationID")

	if errors := h.validator.ValidateUpdate(ctx, request.Sanitize()); len(errors) != 0 {
		msg := fmt.Sprintf("validation errors [%s], while updating [phone-call] integration with request [%s]", spew.Sdump(errors), c.Body())
		ctxLogger.Warn(stacktrace.NewError(msg))
		return h.responseUnprocessableEntity(c, errors, "validation errors while updating integration")
	}

	authUser := h.userFromContext(c)
	integration, err := h.service.Update(ctx, request.ToUpdateParams(c.OriginalURL(), authUser.ID))
	if stacktrace.GetCode(err) == repositories.ErrCodeNotFound {
		msg := fmt.Sprintf("cannot find [phone-call] integration with id [%s] for user [%s]", request.IntegrationID, authUser.ID)
		ctxLogger.Warn(stacktrace.Propagate(err, msg))
		return h.responseNotFound(c, msg)
	}

	if err != nil {
		msg := fmt.Sprintf("cannot update [phone-call] integration [%s] for user with ID [%s]", request.IntegrationID, authUser.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return h.responseInternalServerError(c)
	}

	return h.responseOK(c, "integration updated successfully", integration)
}

// @Summary      Delete a phone call integration
// @Description  This endpoint deletes a phone call integration for a project
// @Security	 BearerAuth
// @Tags         PhoneCallIntegration
// @Produce      json
// @Success      200 		{object}	responses.NoContent
// @Failure      400		{object}	responses.BadRequest
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure 	 404    	{object}	responses.NotFound
// @Failure      422		{object}	responses.UnprocessableEntity
// @Failure      500		{object}	responses.InternalServerError
// @Router       /projects/{projectID}/phone-call-integrations/{integrationID} [delete]
func (h *PhoneCallIntegrationHandler) delete(c *fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	if errors := h.mergeErrors(h.validateUUID(c, "projectID"), h.validateUUID(c, "integrationID")); len(errors) != 0 {
		msg := fmt.Sprintf("validation errors [%s], while deleting [phone-call] integration with request [%s]", spew.Sdump(errors), c.Body())
		ctxLogger.Warn(stacktrace.NewError(msg))
		return h.responseUnprocessableEntity(c, errors, "validation errors while deleting integration")
	}

	authUser := h.userFromContext(c)
	integrationID := uuid.MustParse(c.Params("integrationID"))
	projectID := uuid.MustParse(c.Params("projectID"))

	err := h.service.Delete(ctx, &services.IntegrationDeleteParams{
		Source:        c.OriginalURL(),
		IntegrationID: integrationID,
		ProjectID:     projectID,
		UserID:        authUser.ID,
	})
	if stacktrace.GetCode(err) == repositories.ErrCodeNotFound {
		msg := fmt.Sprintf("cannot find [phone-call] integration with id [%s] for user [%s]", integrationID, authUser.ID)
		ctxLogger.Warn(stacktrace.Propagate(err, msg))
		return h.responseNotFound(c, msg)
	}

	if err != nil {
		msg := fmt.Sprintf("cannot delete integration [%s] for user with ID [%s]", integrationID, authUser.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return h.responseInternalServerError(c)
	}

	return h.responseNoContent(c, "integration deleted successfully")
}
