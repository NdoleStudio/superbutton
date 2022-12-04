package handlers

import (
	"fmt"

	"github.com/NdoleStudio/superbutton/pkg/entities"
	"github.com/NdoleStudio/superbutton/pkg/repositories"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"

	"github.com/NdoleStudio/superbutton/pkg/services"
	"github.com/NdoleStudio/superbutton/pkg/telemetry"
	"github.com/gofiber/fiber/v2"
	"github.com/palantir/stacktrace"
)

// ProjectSettingsHandler handles user http requests.
type ProjectSettingsHandler struct {
	handler
	logger  telemetry.Logger
	tracer  telemetry.Tracer
	service *services.ProjectSettingsService
}

// NewProjectSettingsHandler creates a new ProjectSettingsHandler
func NewProjectSettingsHandler(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	service *services.ProjectSettingsService,
) (h *ProjectSettingsHandler) {
	return &ProjectSettingsHandler{
		logger:  logger.WithService(fmt.Sprintf("%T", h)),
		tracer:  tracer,
		service: service,
	}
}

// RegisterRoutes registers the routes for the MessageHandler
func (h *ProjectSettingsHandler) RegisterRoutes(app *fiber.App, middlewares ...fiber.Handler) {
	router := app.Group("/v1/settings/:userID/projects")
	router.Get("/:projectID", h.computeRoute(middlewares, h.show)...)
}

// show returns all the settings for a project
// @Summary      Project Settings
// @Description  Fetches all the settings and integrations of a project
// @Security	 BearerAuth
// @Tags         ProjectSettings
// @Produce      json
// @Param 		 userID			path 		string true "User ID"
// @Param 		 projectID		path 		string true "Project ID"
// @Success      200 			{object}	responses.Ok[entities.ProjectSettings]
// @Failure      400			{object}	responses.BadRequest
// @Failure 	 401    		{object}	responses.Unauthorized
// @Failure 	 404    		{object}	responses.NotFound
// @Failure      422			{object}	responses.UnprocessableEntity
// @Failure      500			{object}	responses.InternalServerError
// @Router       /settings/{userID}/projects/{projectID} 	[get]
func (h *ProjectSettingsHandler) show(c *fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	if errors := h.mergeErrors(h.validateUUID(c, "projectID")); len(errors) != 0 {
		msg := fmt.Sprintf("validation errors [%s], while fetching settings for project with URL [%s]", spew.Sdump(errors), c.OriginalURL())
		ctxLogger.Warn(stacktrace.NewError(msg))
		return h.responseUnprocessableEntity(c, errors, "validation errors while fetching project settings")
	}

	userID := entities.UserID(c.Params("userID"))
	projectID := uuid.MustParse(c.Params("projectID"))

	settings, err := h.service.Get(ctx, userID, projectID)
	if stacktrace.GetCode(err) == repositories.ErrCodeNotFound {
		msg := fmt.Sprintf("cannot find settings for project with id [%s] for user [%s]", projectID, userID)
		ctxLogger.Warn(stacktrace.Propagate(err, msg))
		return h.responseNotFound(c, msg)
	}

	if err != nil {
		msg := fmt.Sprintf("cannot get settings for project [%s] and user [%s]", projectID, userID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return h.responseInternalServerError(c)
	}

	return h.responseOK(c, "project settings fetched successfully", settings)
}
