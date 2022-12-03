package handlers

import (
	"fmt"

	"github.com/NdoleStudio/superbutton/pkg/entities"

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

// LinkIntegrationHandler handles user http requests.
type LinkIntegrationHandler struct {
	handler
	integrationType entities.IntegrationType
	logger          telemetry.Logger
	tracer          telemetry.Tracer
	validator       *validators.LinkIntegrationHandlerValidator
	service         *services.LinkIntegrationService
}

// NewLinkIntegrationHandler creates a new LinkIntegrationHandler
func NewLinkIntegrationHandler(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	validator *validators.LinkIntegrationHandlerValidator,
	service *services.LinkIntegrationService,
) (h *LinkIntegrationHandler) {
	return &LinkIntegrationHandler{
		logger:          logger.WithService(fmt.Sprintf("%T", h)),
		tracer:          tracer,
		integrationType: entities.IntegrationTypeLink,
		validator:       validator,
		service:         service,
	}
}

// RegisterRoutes registers the routes for the ContentIntegrationHandler
func (h *LinkIntegrationHandler) RegisterRoutes(app *fiber.App, middlewares []fiber.Handler) {
	router := app.Group("/v1/projects/:projectID/link-integrations")
	router.Post("/", h.computeRoute(middlewares, h.create)...)
	router.Get("/:integrationID", h.computeRoute(middlewares, h.show)...)
	router.Put("/:integrationID", h.computeRoute(middlewares, h.update)...)
	router.Delete("/:integrationID", h.computeRoute(middlewares, h.delete)...)
}

// @Summary      Get link integration
// @Description  Fetches a specific link integration
// @Security	 BearerAuth
// @Tags         LinkIntegration
// @Produce      json
// @Success      200 		{object}	responses.Ok[entities.LinkIntegration]
// @Failure      400		{object}	responses.BadRequest
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure 	 404    	{object}	responses.NotFound
// @Failure      422		{object}	responses.UnprocessableEntity
// @Failure      500		{object}	responses.InternalServerError
// @Router       /projects/:projectID/link-integrations/:integrationID 	[get]
func (h *LinkIntegrationHandler) show(c *fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	if errors := h.mergeErrors(h.validateUUID(c, "integrationID")); len(errors) != 0 {
		msg := fmt.Sprintf("validation errors [%s], while fetching [%s] integration with ID [%s]", spew.Sdump(errors), h.integrationType, c.Params("integrationID"))
		ctxLogger.Warn(stacktrace.NewError(msg))
		return h.responseUnprocessableEntity(c, errors, "validation errors while fetching link integration")
	}

	integrationID := uuid.MustParse(c.Params("integrationID"))
	authUser := h.userFromContext(c)

	integration, err := h.service.Get(ctx, authUser.ID, integrationID)
	if err != nil {
		msg := fmt.Sprintf("cannot fetch [%s] intergration [%s] for user with ID [%s]", h.integrationType, authUser.ID, integrationID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return h.responseInternalServerError(c)
	}

	return h.responseOK(c, "link integration fetched successfully", integration)
}

// @Summary      Create a link integration
// @Description  This endpoint creates a new link integration for a project
// @Security	 BearerAuth
// @Tags         LinkIntegration
// @Produce      json
// @Param        payload	body 		requests.LinkIntegrationCreateRequest	true 	"link integration create payload"
// @Success      200 		{object}	responses.Ok[entities.LinkIntegration]
// @Failure      400		{object}	responses.BadRequest
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure      422		{object}	responses.UnprocessableEntity
// @Failure      500		{object}	responses.InternalServerError
// @Router       /projects/{projectID}/link-integrations [post]
func (h *LinkIntegrationHandler) create(c *fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	var request requests.LinkIntegrationCreateRequest
	if err := c.BodyParser(&request); err != nil {
		msg := fmt.Sprintf("cannot marshall params [%s] into %T", c.OriginalURL(), request)
		ctxLogger.Warn(stacktrace.Propagate(err, msg))
		return h.responseBadRequest(c, err)
	}
	request.ProjectID = c.Params("projectID")

	if errors := h.validator.ValidateCreate(ctx, request.Sanitize()); len(errors) != 0 {
		msg := fmt.Sprintf("validation errors [%s], while creating [%s] integration with request [%s]", h.integrationType, spew.Sdump(errors), c.Body())
		ctxLogger.Warn(stacktrace.NewError(msg))
		return h.responseUnprocessableEntity(c, errors, "validation errors while creating integration")
	}

	authUser := h.userFromContext(c)
	integration, err := h.service.Create(ctx, request.ToCreateParams(c.OriginalURL(), authUser.ID))
	if err != nil {
		msg := fmt.Sprintf("cannot create [%s] integration for user with ID [%s]", h.integrationType, authUser.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return h.responseInternalServerError(c)
	}

	return h.responseOK(c, "link integration created successfully", integration)
}

// @Summary      Update a link integration
// @Description  This endpoint updates a link integration for a project
// @Security	 BearerAuth
// @Tags         LinkIntegration
// @Produce      json
// @Param        payload	body 		requests.LinkIntegrationUpdateRequest	true 	"link integration update payload"
// @Success      200 		{object}	responses.Ok[entities.LinkIntegration]
// @Failure      400		{object}	responses.BadRequest
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure      422		{object}	responses.UnprocessableEntity
// @Failure      500		{object}	responses.InternalServerError
// @Router       /projects/{projectID}/link-integrations/{integrationID} [put]
func (h *LinkIntegrationHandler) update(c *fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	var request requests.LinkIntegrationUpdateRequest
	if err := c.BodyParser(&request); err != nil {
		msg := fmt.Sprintf("cannot marshall params [%s] into %T", c.OriginalURL(), request)
		ctxLogger.Warn(stacktrace.Propagate(err, msg))
		return h.responseBadRequest(c, err)
	}
	request.IntegrationID = c.Params("integrationID")

	if errors := h.validator.ValidateUpdate(ctx, request.Sanitize()); len(errors) != 0 {
		msg := fmt.Sprintf("validation errors [%s], while updating [%s] integration with request [%s]", spew.Sdump(errors), h.integrationType, c.Body())
		ctxLogger.Warn(stacktrace.NewError(msg))
		return h.responseUnprocessableEntity(c, errors, "validation errors while updating integration")
	}

	authUser := h.userFromContext(c)
	integration, err := h.service.Update(ctx, request.ToUpdateParams(c.OriginalURL(), authUser.ID))
	if stacktrace.GetCode(err) == repositories.ErrCodeNotFound {
		msg := fmt.Sprintf("cannot find [%s] integration with id [%s] for user [%s]", h.integrationType, request.IntegrationID, authUser.ID)
		ctxLogger.Warn(stacktrace.Propagate(err, msg))
		return h.responseNotFound(c, msg)
	}

	if err != nil {
		msg := fmt.Sprintf("cannot update [%s] integration [%s] for user with ID [%s]", h.integrationType, request.IntegrationID, authUser.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return h.responseInternalServerError(c)
	}

	return h.responseOK(c, "integration updated successfully", integration)
}

// @Summary      Delete a link integration
// @Description  This endpoint deletes a link integration for a project
// @Security	 BearerAuth
// @Tags         LinkIntegration
// @Produce      json
// @Success      200 		{object}	responses.NoContent
// @Failure      400		{object}	responses.BadRequest
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure 	 404    	{object}	responses.NotFound
// @Failure      422		{object}	responses.UnprocessableEntity
// @Failure      500		{object}	responses.InternalServerError
// @Router       /projects/{projectID}/link-integrations/{integrationID} [delete]
func (h *LinkIntegrationHandler) delete(c *fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	if errors := h.mergeErrors(h.validateUUID(c, "projectID"), h.validateUUID(c, "integrationID")); len(errors) != 0 {
		msg := fmt.Sprintf("validation errors [%s], while deleting [%s] integration with request [%s]", h.integrationType, spew.Sdump(errors), c.Body())
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
		msg := fmt.Sprintf("cannot find [%s] integration with id [%s] for user [%s]", h.integrationType, integrationID, authUser.ID)
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
