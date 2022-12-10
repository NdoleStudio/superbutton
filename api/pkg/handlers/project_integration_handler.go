package handlers

import (
	"fmt"

	"github.com/NdoleStudio/superbutton/pkg/repositories"
	"github.com/NdoleStudio/superbutton/pkg/requests"

	"github.com/NdoleStudio/superbutton/pkg/services"
	"github.com/NdoleStudio/superbutton/pkg/telemetry"
	"github.com/davecgh/go-spew/spew"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/palantir/stacktrace"
)

// ProjectIntegrationHandler handles /integration http requests.
type ProjectIntegrationHandler struct {
	handler
	logger  telemetry.Logger
	tracer  telemetry.Tracer
	service *services.ProjectIntegrationService
}

// NewIntegrationHandler creates a new ProjectIntegrationHandler
func NewIntegrationHandler(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	service *services.ProjectIntegrationService,
) (h *ProjectIntegrationHandler) {
	return &ProjectIntegrationHandler{
		logger:  logger.WithService(fmt.Sprintf("%T", h)),
		tracer:  tracer,
		service: service,
	}
}

// RegisterRoutes registers the routes for the MessageHandler
func (h *ProjectIntegrationHandler) RegisterRoutes(app *fiber.App, middlewares []fiber.Handler) {
	router := app.Group("/v1/projects/:projectID/integrations")
	router.Get("/", h.computeRoute(middlewares, h.index)...)
	router.Put("/", h.computeRoute(middlewares, h.update)...)
}

// @Summary      List of project integrations
// @Description  Fetches the list of all integrations for a project
// @Security	 BearerAuth
// @Tags         ProjectIntegrations
// @Produce      json
// @Param 		 projectID	path 		string true "Project ID"
// @Success      200 		{object}	responses.Ok[[]entities.ProjectIntegration]
// @Failure      400		{object}	responses.BadRequest
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure      422		{object}	responses.UnprocessableEntity
// @Failure      500		{object}	responses.InternalServerError
// @Router       /projects/{projectID}/integrations 	[get]
func (h *ProjectIntegrationHandler) index(c *fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	if errors := h.mergeErrors(h.validateUUID(c, "projectID")); len(errors) != 0 {
		msg := fmt.Sprintf("validation errors [%s], while fetchng integrations for project [%s]", spew.Sdump(errors), c.Params("projectID"))
		ctxLogger.Warn(stacktrace.NewError(msg))
		return h.responseUnprocessableEntity(c, errors, "validation errors while fetching integrations")
	}

	authUser := h.userFromContext(c)
	projects, err := h.service.Index(ctx, authUser.ID, uuid.MustParse(c.Params("projectID")))
	if err != nil {
		msg := fmt.Sprintf("cannot fetch projects intergrations for user with ID [%s]", authUser.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return h.responseInternalServerError(c)
	}

	return h.responseOK(c, "integrations fetched successfully", projects)
}

// @Summary      Update project integrations
// @Description  This endpoint updates project integrations for a user
// @Security	 BearerAuth
// @Tags         ProjectIntegrations
// @Produce      json
// @Param 		 projectID	path 		string true "Project ID"
// @Param        payload	body 		requests.ProjectIntegrationsUpdateRequest	true 	"project integrations update payload"
// @Success      200 		{object}	responses.NoContent
// @Failure      400		{object}	responses.BadRequest
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure      422		{object}	responses.UnprocessableEntity
// @Failure      500		{object}	responses.InternalServerError
// @Router       /projects/{projectID}/integrations 	[put]
func (h *ProjectIntegrationHandler) update(c *fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	var request requests.ProjectIntegrationsUpdateRequest
	if err := c.BodyParser(&request); err != nil {
		msg := fmt.Sprintf("cannot marshall params [%s] into %T", c.OriginalURL(), request)
		ctxLogger.Warn(stacktrace.Propagate(err, msg))
		return h.responseBadRequest(c, err)
	}

	err := h.service.Update(ctx, h.userIDFomContext(c), request.Order)
	if stacktrace.GetCode(err) == repositories.ErrCodeNotFound {
		msg := fmt.Sprintf("cannot find project with id [%s] for user [%s]", c.Params("projectID"), h.userIDFomContext(c))
		ctxLogger.Warn(stacktrace.Propagate(err, msg))
		return h.responseNotFound(c, msg)
	}

	if err != nil {
		msg := fmt.Sprintf("cannot update project [%s] integrations for  user with ID [%s]", c.Params("projectID"), h.userIDFomContext(c))
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return h.responseInternalServerError(c)
	}

	return h.responseNoContent(c, "project integrations updated successfully")
}
