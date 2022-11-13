package middlewares

import (
	"context"
	"fmt"
	"strings"

	"github.com/palantir/stacktrace"

	"github.com/NdoleStudio/superbutton/pkg/entities"
	"github.com/NdoleStudio/superbutton/pkg/telemetry"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/idtoken"
)

// GoogleAuth authenticates a user based on the bearer token
func GoogleAuth(logger telemetry.Logger, tracer telemetry.Tracer, audience string, subject string) fiber.Handler {
	logger = logger.WithService("middlewares.GoogleAuth")
	return func(c *fiber.Ctx) error {
		_, span := tracer.StartFromFiberCtx(c, "middlewares.GoogleAuth")
		defer span.End()

		authToken := c.Get(authHeaderBearer)
		if !strings.HasPrefix(authToken, bearerScheme) {
			span.AddEvent(fmt.Sprintf("The request header has no [%s] token", bearerScheme))
			return c.Next()
		}

		if len(authToken) > len(bearerScheme)+1 {
			authToken = authToken[len(bearerScheme)+1:]
		}

		ctxLogger := tracer.CtxLogger(logger, span)

		payload, err := idtoken.Validate(context.Background(), authToken, audience)
		if err != nil {
			msg := fmt.Sprintf("invalid google auto token [%s]", authToken)
			ctxLogger.Error(tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg)))
			return c.Next()
		}

		if payload.Subject != subject {
			msg := fmt.Sprintf("invalid subject [%s] for google auth token [%s]", payload.Subject, authToken)
			ctxLogger.Error(tracer.WrapErrorSpan(span, stacktrace.NewError(msg)))
			return c.Next()
		}

		span.AddEvent(fmt.Sprintf("[%s] google auth token is valid", bearerScheme))

		authUser := entities.AuthUser{
			Email: payload.Claims["email"].(string),
			Name:  strings.Split(payload.Claims["email"].(string), "@")[0],
			ID:    entities.UserID(payload.Claims["sub"].(string)),
		}

		c.Locals(ContextKeyAuthUserID, authUser)

		ctxLogger.Info(fmt.Sprintf("[%T] set successfully for user with ID [%s]", authUser, authUser.ID))
		return c.Next()
	}
}
