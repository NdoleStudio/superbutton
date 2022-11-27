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
	router.Put("/:projectID", h.computeRoute(middlewares, h.update)...)
	router.Delete("/:projectID", h.computeRoute(middlewares, h.delete)...)
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
// @Param        payload	body 		requests.ProjectCreateRequest	true 	"project create payload"
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

// @Summary      Update a project
// @Description  This endpoint updates a project for a user
// @Security	 BearerAuth
// @Tags         Projects
// @Produce      json
// @Param        payload	body 		requests.ProjectUpdateRequest	true 	"project update payload"
// @Success      200 		{object}	responses.Ok[entities.Project]
// @Failure      400		{object}	responses.BadRequest
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure      422		{object}	responses.UnprocessableEntity
// @Failure      500		{object}	responses.InternalServerError
// @Router       /projects/{projectID} 	[put]
func (h *ProjectHandler) update(c *fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	var request requests.ProjectUpdateRequest
	if err := c.BodyParser(&request); err != nil {
		msg := fmt.Sprintf("cannot marshall params [%s] into %T", c.OriginalURL(), request)
		ctxLogger.Warn(stacktrace.Propagate(err, msg))
		return h.responseBadRequest(c, err)
	}
	request.ProjectID = c.Params("projectID")

	if errors := h.validator.ValidateUpdate(ctx, request.Sanitize()); len(errors) != 0 {
		msg := fmt.Sprintf("validation errors [%s], while creating project with request [%s]", spew.Sdump(errors), c.Body())
		ctxLogger.Warn(stacktrace.NewError(msg))
		return h.responseUnprocessableEntity(c, errors, "validation errors while creating project")
	}

	authUser := h.userFromContext(c)
	project, err := h.service.Update(ctx, request.ToProjectUpdatePrams(c.OriginalURL(), authUser.ID))
	if stacktrace.GetCode(err) == repositories.ErrCodeNotFound {
		msg := fmt.Sprintf("cannot find project with id [%s] for user [%s]", request.ProjectID, authUser.ID)
		ctxLogger.Warn(stacktrace.Propagate(err, msg))
		return h.responseNotFound(c, msg)
	}

	if err != nil {
		msg := fmt.Sprintf("cannot update project [%s] user with ID [%s]", request.ProjectID, authUser.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return h.responseInternalServerError(c)
	}

	return h.responseOK(c, "project updated successfully", project)
}

// @Summary      Delete a project
// @Description  This endpoint deletes a project
// @Security	 BearerAuth
// @Tags         Projects
// @Produce      json
// @Success      200 		{object}	responses.NoContent
// @Failure      400		{object}	responses.BadRequest
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure 	 404    	{object}	responses.NotFound
// @Failure      422		{object}	responses.UnprocessableEntity
// @Failure      500		{object}	responses.InternalServerError
// @Router       /projects/{projectID} [delete]
func (h *ProjectHandler) delete(c *fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	if errors := h.mergeErrors(h.validateUUID(c, "projectID")); len(errors) != 0 {
		msg := fmt.Sprintf("validation errors [%s], while deleting project with url [%s]", spew.Sdump(errors), c.OriginalURL())
		ctxLogger.Warn(stacktrace.NewError(msg))
		return h.responseUnprocessableEntity(c, errors, "validation errors while deleting project")
	}

	authUser := h.userFromContext(c)
	projectID := uuid.MustParse(c.Params("projectID"))

	err := h.service.Delete(ctx, c.OriginalURL(), authUser.ID, projectID)
	if stacktrace.GetCode(err) == repositories.ErrCodeNotFound {
		msg := fmt.Sprintf("cannot delete project with id [%s] for user [%s]", projectID, authUser.ID)
		ctxLogger.Warn(stacktrace.Propagate(err, msg))
		return h.responseNotFound(c, msg)
	}

	if err != nil {
		msg := fmt.Sprintf("cannot delete project [%s] for user with ID [%s]", projectID, authUser.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return h.responseInternalServerError(c)
	}

	return h.responseNoContent(c, "project deleted successfully")
}
