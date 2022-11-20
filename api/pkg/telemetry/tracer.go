package telemetry

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
)

const (
	// TracerContextKey stores the fiber trace context
	TracerContextKey = "tracer.context.key"
)

// Tracer is used for tracing
type Tracer interface {
	// StartFromFiberCtx creates a context.Context and trace.Span from fiber.Ctx.
	StartFromFiberCtx(c *fiber.Ctx, name ...string) (context.Context, trace.Span)

	// StartFromFiberCtxWithLogger creates a context.Context and trace.Span from fiber.Ctx.
	StartFromFiberCtxWithLogger(c *fiber.Ctx, logger Logger, name ...string) (context.Context, trace.Span, Logger)

	// Start creates a context.Context and trace.Span
	Start(c context.Context, name ...string) (context.Context, trace.Span)

	// StartWithLogger creates a context.Context and trace.Span from fiber.Ctx.
	StartWithLogger(c context.Context, logger Logger, name ...string) (context.Context, trace.Span, Logger)

	// CtxLogger creates a telemetry.Logger with spanContext attributes in the structured logger
	CtxLogger(logger Logger, span trace.Span) Logger

	// WrapErrorSpan sets a spanContext as error
	WrapErrorSpan(span trace.Span, err error) error

	// Span returns the trace.Span from context.Context
	Span(ctx context.Context) trace.Span
}
