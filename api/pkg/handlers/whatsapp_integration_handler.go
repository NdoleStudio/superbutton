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

// WhatsappIntegrationHandler handles user http requests.
type WhatsappIntegrationHandler struct {
	handler
	logger    telemetry.Logger
	tracer    telemetry.Tracer
	validator *validators.WhatsappIntegrationHandlerValidator
	service   *services.WhatsappIntegrationService
}

// NewWhatsappIntegrationHandler creates a new WhatsappIntegrationHandler
func NewWhatsappIntegrationHandler(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	validator *validators.WhatsappIntegrationHandlerValidator,
	service *services.WhatsappIntegrationService,
) (h *WhatsappIntegrationHandler) {
	return &WhatsappIntegrationHandler{
		logger:    logger.WithService(fmt.Sprintf("%T", h)),
		tracer:    tracer,
		validator: validator,
		service:   service,
	}
}

// RegisterRoutes registers the routes for the MessageHandler
func (h *WhatsappIntegrationHandler) RegisterRoutes(app *fiber.App, middlewares []fiber.Handler) {
	router := app.Group("/v1/projects/:projectID/whatsapp-integrations")
	router.Post("/", h.computeRoute(middlewares, h.create)...)
	router.Get("/:integrationID", h.computeRoute(middlewares, h.show)...)
	router.Put("/:integrationID", h.computeRoute(middlewares, h.update)...)
	router.Delete("/:integrationID", h.computeRoute(middlewares, h.delete)...)
}

// @Summary      Get whatsapp integration
// @Description  Fetches a specific whatsapp integration
// @Security	 BearerAuth
// @Tags         WhatsappIntegration
// @Produce      json
// @Param 		 projectID		path 		string true "Project ID"
// @Param 		 integrationID	path 		string true "Integration ID"
// @Success      200 			{object}	responses.Ok[entities.WhatsappIntegration]
// @Failure      400			{object}	responses.BadRequest
// @Failure 	 401    		{object}	responses.Unauthorized
// @Failure 	 404    		{object}	responses.NotFound
// @Failure      422			{object}	responses.UnprocessableEntity
// @Failure      500			{object}	responses.InternalServerError
// @Router       /projects/{projectID}/whatsapp-integrations/{integrationID} 	[get]
func (h *WhatsappIntegrationHandler) show(c *fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	if errors := h.mergeErrors(h.validateUUID(c, "integrationID")); len(errors) != 0 {
		msg := fmt.Sprintf("validation errors [%s], while deleting integration with request [%s]", spew.Sdump(errors), c.Body())
		ctxLogger.Warn(stacktrace.NewError(msg))
		return h.responseUnprocessableEntity(c, errors, "validation errors while deleting integration")
	}

	integrationID := uuid.MustParse(c.Params("integrationID"))
	authUser := h.userFromContext(c)

	integration, err := h.service.Get(ctx, authUser.ID, integrationID)
	if err != nil {
		msg := fmt.Sprintf("cannot fetch intergration [%s] for user with ID [%s]", authUser.ID, integrationID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return h.responseInternalServerError(c)
	}

	return h.responseOK(c, "integration fetched successfully", integration)
}

// @Summary      Create a WhatsappIntegration
// @Description  This endpoint creates a new whatsapp integration for a project
// @Security	 BearerAuth
// @Tags         WhatsappIntegration
// @Produce      json
// @Param 		 projectID	path 		string true "Project ID"
// @Param        payload	body 		requests.WhatsappIntegrationCreateRequest	true 	"whatsapp integration create payload"
// @Success      200 		{object}	responses.Ok[entities.WhatsappIntegration]
// @Failure      400		{object}	responses.BadRequest
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure      422		{object}	responses.UnprocessableEntity
// @Failure      500		{object}	responses.InternalServerError
// @Router       /projects/{projectID}/whatsapp-integrations [post]
func (h *WhatsappIntegrationHandler) create(c *fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	var request requests.WhatsappIntegrationCreateRequest
	if err := c.BodyParser(&request); err != nil {
		msg := fmt.Sprintf("cannot marshall params [%s] into %T", c.OriginalURL(), request)
		ctxLogger.Warn(stacktrace.Propagate(err, msg))
		return h.responseBadRequest(c, err)
	}
	request.ProjectID = c.Params("projectID")

	if errors := h.validator.ValidateCreate(ctx, request.Sanitize()); len(errors) != 0 {
		msg := fmt.Sprintf("validation errors [%s], while creating project with request [%s]", spew.Sdump(errors), c.Body())
		ctxLogger.Warn(stacktrace.NewError(msg))
		return h.responseUnprocessableEntity(c, errors, "validation errors while creating project")
	}

	authUser := h.userFromContext(c)
	integration, err := h.service.Create(ctx, request.ToCreateParams(c.OriginalURL(), authUser.ID))
	if err != nil {
		msg := fmt.Sprintf("cannot create whatsapp integration for user with ID [%s]", authUser.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return h.responseInternalServerError(c)
	}

	return h.responseOK(c, "whatsapp integration created successfully", integration)
}

// @Summary      Update a whatsapp integration
// @Description  This endpoint updates a whatsapp integration for a project
// @Security	 BearerAuth
// @Tags         WhatsappIntegration
// @Produce      json
// @Param 		 projectID		path 		string true "Project ID"
// @Param 		 integrationID	path 		string true "Integration ID"
// @Param        payload		body 		requests.WhatsappIntegrationUpdateRequest	true 	"whatsapp integration update payload"
// @Success      200 			{object}	responses.Ok[entities.WhatsappIntegration]
// @Failure      400			{object}	responses.BadRequest
// @Failure 	 401    		{object}	responses.Unauthorized
// @Failure      422			{object}	responses.UnprocessableEntity
// @Failure      500			{object}	responses.InternalServerError
// @Router       /projects/{projectID}/whatsapp-integrations/{integrationID} [put]
func (h *WhatsappIntegrationHandler) update(c *fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	var request requests.WhatsappIntegrationUpdateRequest
	if err := c.BodyParser(&request); err != nil {
		msg := fmt.Sprintf("cannot marshall params [%s] into %T", c.OriginalURL(), request)
		ctxLogger.Warn(stacktrace.Propagate(err, msg))
		return h.responseBadRequest(c, err)
	}
	request.IntegrationID = c.Params("integrationID")

	if errors := h.validator.ValidateUpdate(ctx, request.Sanitize()); len(errors) != 0 {
		msg := fmt.Sprintf("validation errors [%s], while updating project with request [%s]", spew.Sdump(errors), c.Body())
		ctxLogger.Warn(stacktrace.NewError(msg))
		return h.responseUnprocessableEntity(c, errors, "validation errors while updating integration")
	}

	authUser := h.userFromContext(c)
	integration, err := h.service.Update(ctx, request.ToUpdateParams(c.OriginalURL(), authUser.ID))
	if stacktrace.GetCode(err) == repositories.ErrCodeNotFound {
		msg := fmt.Sprintf("cannot find integration with id [%s] for user [%s]", request.IntegrationID, authUser.ID)
		ctxLogger.Warn(stacktrace.Propagate(err, msg))
		return h.responseNotFound(c, msg)
	}

	if err != nil {
		msg := fmt.Sprintf("cannot update integration [%s] for user with ID [%s]", request.IntegrationID, authUser.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return h.responseInternalServerError(c)
	}

	return h.responseOK(c, "integration updated successfully", integration)
}

// @Summary      Delete a whatsapp integration
// @Description  This endpoint deletes a whatsapp integration for a project
// @Security	 BearerAuth
// @Tags         WhatsappIntegration
// @Produce      json
// @Param 		 projectID		path 		string true "Project ID"
// @Param 		 integrationID	path 		string true "Integration ID"
// @Success      200 			{object}	responses.NoContent
// @Failure      400			{object}	responses.BadRequest
// @Failure 	 401    		{object}	responses.Unauthorized
// @Failure 	 404    		{object}	responses.NotFound
// @Failure      422			{object}	responses.UnprocessableEntity
// @Failure      500			{object}	responses.InternalServerError
// @Router       /projects/{projectID}/whatsapp-integrations/{integrationID} [delete]
func (h *WhatsappIntegrationHandler) delete(c *fiber.Ctx) error {
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
		msg := fmt.Sprintf("cannot find whatsapp integration with id [%s] for user [%s]", integrationID, authUser.ID)
		ctxLogger.Warn(stacktrace.Propagate(err, msg))
		return h.responseNotFound(c, msg)
	}

	if err != nil {
		msg := fmt.Sprintf("cannot delete whatsapp integration [%s] for user with ID [%s]", integrationID, authUser.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return h.responseInternalServerError(c)
	}

	return h.responseNoContent(c, "integration deleted successfully")
}
