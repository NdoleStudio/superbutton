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

// ContentIntegrationHandler handles user http requests.
type ContentIntegrationHandler struct {
	handler
	logger    telemetry.Logger
	tracer    telemetry.Tracer
	validator *validators.ContentIntegrationHandlerValidator
	service   *services.ContentIntegrationService
}

// NewContentIntegrationHandler creates a new ContentIntegrationHandler
func NewContentIntegrationHandler(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	validator *validators.ContentIntegrationHandlerValidator,
	service *services.ContentIntegrationService,
) (h *ContentIntegrationHandler) {
	return &ContentIntegrationHandler{
		logger:    logger.WithService(fmt.Sprintf("%T", h)),
		tracer:    tracer,
		validator: validator,
		service:   service,
	}
}

// RegisterRoutes registers the routes for the ContentIntegrationHandler
func (h *ContentIntegrationHandler) RegisterRoutes(app *fiber.App, middlewares []fiber.Handler) {
	router := app.Group("/v1/projects/:projectID/content-integrations")
	router.Post("/", h.computeRoute(middlewares, h.create)...)
	router.Get("/:integrationID", h.computeRoute(middlewares, h.show)...)
	router.Put("/:integrationID", h.computeRoute(middlewares, h.update)...)
	router.Delete("/:integrationID", h.computeRoute(middlewares, h.delete)...)
}

// @Summary      Get content integration
// @Description  Fetches a specific content integration
// @Security	 BearerAuth
// @Tags         ContentIntegration
// @Produce      json
// @Success      200 		{object}	responses.Ok[entities.ContentIntegration]
// @Failure      400		{object}	responses.BadRequest
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure 	 404    	{object}	responses.NotFound
// @Failure      422		{object}	responses.UnprocessableEntity
// @Failure      500		{object}	responses.InternalServerError
// @Router       /projects/{projectID}/content-integrations/{integrationID} 	[get]
func (h *ContentIntegrationHandler) show(c *fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	if errors := h.mergeErrors(h.validateUUID(c, "integrationID")); len(errors) != 0 {
		msg := fmt.Sprintf("validation errors [%s], while fetching content integration with ID [%s]", spew.Sdump(errors), c.Params("integrationID"))
		ctxLogger.Warn(stacktrace.NewError(msg))
		return h.responseUnprocessableEntity(c, errors, "validation errors while fetching content integration")
	}

	integrationID := uuid.MustParse(c.Params("integrationID"))
	authUser := h.userFromContext(c)

	integration, err := h.service.Get(ctx, authUser.ID, integrationID)
	if err != nil {
		msg := fmt.Sprintf("cannot fetch text intergration [%s] for user with ID [%s]", authUser.ID, integrationID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return h.responseInternalServerError(c)
	}

	return h.responseOK(c, "content integration fetched successfully", integration)
}

// @Summary      Create a content integration
// @Description  This endpoint creates a new content integration for a project
// @Security	 BearerAuth
// @Tags         ContentIntegration
// @Produce      json
// @Param        payload	body 		requests.ContentIntegrationCreateRequest	true 	"content integration create payload"
// @Success      200 		{object}	responses.Ok[entities.ContentIntegration]
// @Failure      400		{object}	responses.BadRequest
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure      422		{object}	responses.UnprocessableEntity
// @Failure      500		{object}	responses.InternalServerError
// @Router       /projects/{projectID}/content-integrations [post]
func (h *ContentIntegrationHandler) create(c *fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	var request requests.ContentIntegrationCreateRequest
	if err := c.BodyParser(&request); err != nil {
		msg := fmt.Sprintf("cannot marshall params [%s] into %T", c.OriginalURL(), request)
		ctxLogger.Warn(stacktrace.Propagate(err, msg))
		return h.responseBadRequest(c, err)
	}
	request.ProjectID = c.Params("projectID")

	if errors := h.validator.ValidateCreate(ctx, request.Sanitize()); len(errors) != 0 {
		msg := fmt.Sprintf("validation errors [%s], while content integration with request [%s]", spew.Sdump(errors), c.Body())
		ctxLogger.Warn(stacktrace.NewError(msg))
		return h.responseUnprocessableEntity(c, errors, "validation errors while integration")
	}

	authUser := h.userFromContext(c)
	integration, err := h.service.Create(ctx, request.ToCreateParams(c.OriginalURL(), authUser.ID))
	if err != nil {
		msg := fmt.Sprintf("cannot create content integration for user with ID [%s]", authUser.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return h.responseInternalServerError(c)
	}

	return h.responseOK(c, "content integration created successfully", integration)
}

// @Summary      Update a content integration
// @Description  This endpoint updates a content integration for a project
// @Security	 BearerAuth
// @Tags         ContentIntegration
// @Produce      json
// @Param        payload	body 		requests.ContentIntegrationUpdateRequest	true 	"content integration update payload"
// @Success      200 		{object}	responses.Ok[entities.ContentIntegration]
// @Failure      400		{object}	responses.BadRequest
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure      422		{object}	responses.UnprocessableEntity
// @Failure      500		{object}	responses.InternalServerError
// @Router       /projects/{projectID}/content-integrations/{integrationID} [put]
func (h *ContentIntegrationHandler) update(c *fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	var request requests.ContentIntegrationUpdateRequest
	if err := c.BodyParser(&request); err != nil {
		msg := fmt.Sprintf("cannot marshall params [%s] into %T", c.OriginalURL(), request)
		ctxLogger.Warn(stacktrace.Propagate(err, msg))
		return h.responseBadRequest(c, err)
	}
	request.IntegrationID = c.Params("integrationID")

	if errors := h.validator.ValidateUpdate(ctx, request.Sanitize()); len(errors) != 0 {
		msg := fmt.Sprintf("validation errors [%s], while updating content integration with request [%s]", spew.Sdump(errors), c.Body())
		ctxLogger.Warn(stacktrace.NewError(msg))
		return h.responseUnprocessableEntity(c, errors, "validation errors while updating integration")
	}

	authUser := h.userFromContext(c)
	integration, err := h.service.Update(ctx, request.ToUpdateParams(c.OriginalURL(), authUser.ID))
	if stacktrace.GetCode(err) == repositories.ErrCodeNotFound {
		msg := fmt.Sprintf("cannot find content integration with id [%s] for user [%s]", request.IntegrationID, authUser.ID)
		ctxLogger.Warn(stacktrace.Propagate(err, msg))
		return h.responseNotFound(c, msg)
	}

	if err != nil {
		msg := fmt.Sprintf("cannot update content integration [%s] for user with ID [%s]", request.IntegrationID, authUser.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return h.responseInternalServerError(c)
	}

	return h.responseOK(c, "integration updated successfully", integration)
}

// @Summary      Delete a content integration
// @Description  This endpoint deletes a content integration for a project
// @Security	 BearerAuth
// @Tags         ContentIntegration
// @Produce      json
// @Success      200 		{object}	responses.NoContent
// @Failure      400		{object}	responses.BadRequest
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure 	 404    	{object}	responses.NotFound
// @Failure      422		{object}	responses.UnprocessableEntity
// @Failure      500		{object}	responses.InternalServerError
// @Router       /projects/{projectID}/content-integrations/{integrationID} [delete]
func (h *ContentIntegrationHandler) delete(c *fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	if errors := h.mergeErrors(h.validateUUID(c, "projectID"), h.validateUUID(c, "integrationID")); len(errors) != 0 {
		msg := fmt.Sprintf("validation errors [%s], while deleting integration with request [%s]", spew.Sdump(errors), c.Body())
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
		msg := fmt.Sprintf("cannot find content integration with id [%s] for user [%s]", integrationID, authUser.ID)
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
