package validators

import (
	"context"
	"fmt"
	"net/url"

	"github.com/NdoleStudio/superbutton/pkg/requests"

	"github.com/NdoleStudio/superbutton/pkg/telemetry"
	"github.com/thedevsaddam/govalidator"
)

// ProjectHandlerValidator validates models used in handlers.ProjectHandler
type ProjectHandlerValidator struct {
	logger telemetry.Logger
	tracer telemetry.Tracer
}

// NewProjectHandlerValidator creates a new handlers.ProjectHandler validator
func NewProjectHandlerValidator(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
) (v *ProjectHandlerValidator) {
	return &ProjectHandlerValidator{
		logger: logger.WithService(fmt.Sprintf("%T", v)),
		tracer: tracer,
	}
}

func (validator *ProjectHandlerValidator) ValidateCreate(ctx context.Context, request *requests.ProjectCreateRequest) url.Values {
	_, span := validator.tracer.Start(ctx)
	defer span.End()

	v := govalidator.New(govalidator.Options{
		Data: request,
		Rules: govalidator.MapData{
			"name": []string{
				"required",
				"min:1",
				"max:30",
			},
			"website": []string{
				"required",
				"url",
				"max:255",
			},
		},
	})
	return v.ValidateStruct()
}
