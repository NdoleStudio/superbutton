package handlers

import (
	"fmt"

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
}

// @Summary      List of project integrations
// @Description  Fetches the list of all integrations for a project
// @Security	 BearerAuth
// @Tags         ProjectIntegrations
// @Produce      json
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
