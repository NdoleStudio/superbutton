package handlers

import (
	"fmt"
	"net/url"

	"github.com/google/uuid"

	"github.com/NdoleStudio/superbutton/pkg/entities"
	"github.com/NdoleStudio/superbutton/pkg/middlewares"

	"github.com/gofiber/fiber/v2"
)

// handler is the base struct for handling requests
type handler struct{}

func (h *handler) responseBadRequest(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"status":  "error",
		"message": "The request isn't properly formed",
		"data":    err,
	})
}

func (h *handler) computeRoute(middlewares []fiber.Handler, route fiber.Handler) []fiber.Handler {
	return append(append([]fiber.Handler{}, middlewares...), route)
}

func (h *handler) responseInternalServerError(c *fiber.Ctx) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"status":  "error",
		"message": "We ran into an internal error while handling the request.",
	})
}

//func (h *handler) responseUnauthorized(c *fiber.Ctx) error {
//	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
//		"status":  "error",
//		"message": "You are not authorized to carry out this request.",
//		"data":    "Make sure your API key is set in the [X-API-Key] header in the request",
//	})
//}

//func (h *handler) responseForbidden(c *fiber.Ctx) error {
//	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
//		"status":  "error",
//		"message": fiber.ErrForbidden.Message,
//	})
//}

func (h *handler) responseUnprocessableEntity(c *fiber.Ctx, errors url.Values, message string) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
		"status":  "error",
		"message": message,
		"data":    errors,
	})
}

func (h *handler) responseNotFound(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"status":  "error",
		"message": message,
	})
}

func (h *handler) responseNoContent(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"status":  "success",
		"message": message,
	})
}

func (h *handler) responseOK(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": message,
		"data":    data,
	})
}

func (h *handler) mergeErrors(errors ...url.Values) url.Values {
	result := url.Values{}
	for _, item := range errors {
		for key, values := range item {
			for _, value := range values {
				result.Add(key, value)
			}
		}
	}
	return result
}

func (h *handler) validateUUID(c *fiber.Ctx, param string) url.Values {
	_, err := uuid.Parse(c.Params(param))
	if err != nil {
		return url.Values{
			param: []string{
				fmt.Sprintf("%s is not a valid UUID string e.g b05b8cc4-6e13-11ed-a1eb-0242ac120002", param),
			},
		}
	}
	return nil
}

//func (h *handler) responseCreated(c *fiber.Ctx, message string, data interface{}) error {
//	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
//		"status":  "success",
//		"message": message,
//		"data":    data,
//	})
//}

//func (h *handler) pluralize(value string, count int) string {
//	if count == 1 {
//		return value
//	}
//	return value + "s"
//}

func (h *handler) userFromContext(c *fiber.Ctx) entities.AuthUser {
	if tokenUser, ok := c.Locals(middlewares.ContextKeyAuthUserID).(entities.AuthUser); ok && !tokenUser.IsNoop() {
		return tokenUser
	}
	panic("user does not exist in context.")
}

func (h *handler) userIDFomContext(c *fiber.Ctx) entities.UserID {
	return h.userFromContext(c).ID
}
