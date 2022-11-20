package handlers

import (
	"fmt"

	"github.com/NdoleStudio/superbutton/pkg/requests"
	"github.com/davecgh/go-spew/spew"

	"github.com/NdoleStudio/superbutton/pkg/services"
	"github.com/NdoleStudio/superbutton/pkg/telemetry"
	"github.com/NdoleStudio/superbutton/pkg/validators"
	"github.com/gofiber/fiber/v2"
	"github.com/palantir/stacktrace"
)

// ProjectHandler handles user http requests.
type ProjectHandler struct {
	handler
	logger    telemetry.Logger
	tracer    telemetry.Tracer
	validator *validators.ProjectHandlerValidator
	service   *services.ProjectService
}

// NewProjectHandler creates a new ProjectHandler
func NewProjectHandler(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	validator *validators.ProjectHandlerValidator,
	service *services.ProjectService,
) (h *ProjectHandler) {
	return &ProjectHandler{
		logger:    logger.WithService(fmt.Sprintf("%T", h)),
		tracer:    tracer,
		validator: validator,
		service:   service,
	}
}

// RegisterRoutes registers the routes for the MessageHandler
func (h *ProjectHandler) RegisterRoutes(app *fiber.App, middlewares []fiber.Handler) {
	router := app.Group("/v1/projects")
	router.Get("/", h.computeRoute(middlewares, h.index)...)
	router.Post("/", h.computeRoute(middlewares, h.create)...)
}

// @Summary      List of projects
// @Description  Fetches the list of all projects available to the currently authenticated user
// @Security	 BearerAuth
// @Tags         Projects
// @Produce      json
// @Success      200 		{object}	responses.Ok[[]entities.Project]
// @Failure      400		{object}	responses.BadRequest
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure      422		{object}	responses.UnprocessableEntity
// @Failure      500		{object}	responses.InternalServerError
// @Router       /projects 	[get]
func (h *ProjectHandler) index(c *fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	authUser := h.userFromContext(c)
	projects, err := h.service.Index(ctx, authUser.ID)
	if err != nil {
		msg := fmt.Sprintf("cannot fetch projects for user with ID [%s]", authUser.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return h.responseInternalServerError(c)
	}

	return h.responseOK(c, "projects fetched successfully", projects)
}

// @Summary      Create a project
// @Description  This endpoint creates a new project for a user
// @Security	 BearerAuth
// @Tags         Projects
// @Produce      json
// @Success      200 		{object}	responses.Ok[entities.Project]
// @Failure      400		{object}	responses.BadRequest
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure      422		{object}	responses.UnprocessableEntity
// @Failure      500		{object}	responses.InternalServerError
// @Router       /projects 	[post]
func (h *ProjectHandler) create(c *fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	var request requests.ProjectCreateRequest
	if err := c.BodyParser(&request); err != nil {
		msg := fmt.Sprintf("cannot marshall params [%s] into %T", c.OriginalURL(), request)
		ctxLogger.Warn(stacktrace.Propagate(err, msg))
		return h.responseBadRequest(c, err)
	}

	if errors := h.validator.ValidateCreate(ctx, request.Sanitize()); len(errors) != 0 {
		msg := fmt.Sprintf("validation errors [%s], while creating project with request [%s]", spew.Sdump(errors), c.Body())
		ctxLogger.Warn(stacktrace.NewError(msg))
		return h.responseUnprocessableEntity(c, errors, "validation errors while creating project")
	}

	authUser := h.userFromContext(c)
	project, err := h.service.Create(ctx, request.ToProjectCreateParams(c.OriginalURL(), authUser.ID))
	if err != nil {
		msg := fmt.Sprintf("cannot get user with ID [%s]", authUser.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return h.responseInternalServerError(c)
	}

	return h.responseOK(c, "project created successfully", project)
}
