package validators

import (
	"fmt"

	"github.com/NdoleStudio/superbutton/pkg/telemetry"
)

// UserHandlerValidator validates models used in handlers.UserHandler
type UserHandlerValidator struct {
	logger telemetry.Logger
	tracer telemetry.Tracer
}

// NewUserHandlerValidator creates a new handlers.UserHandler validator
func NewUserHandlerValidator(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
) (v *UserHandlerValidator) {
	return &UserHandlerValidator{
		logger: logger.WithService(fmt.Sprintf("%T", v)),
		tracer: tracer,
	}
}
