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

func (validator *ProjectHandlerValidator) ValidateUpdate(ctx context.Context, request *requests.ProjectUpdateRequest) url.Values {
	_, span := validator.tracer.Start(ctx)
	defer span.End()

	v := govalidator.New(govalidator.Options{
		Data: request,
		Rules: govalidator.MapData{
			"icon": []string{
				"required",
				"in:chat,whatsapp,help-chat",
			},
			"color": []string{
				"required",
				"regex:^#[0-9A-F]{6}$",
			},
			"name": []string{
				"required",
				"min:1",
				"max:30",
			},
			"greeting": []string{
				"max:30",
			},
			"greeting_timeout": []string{
				"max:300",
			},
			"website": []string{
				"required",
				"url",
				"max:255",
			},
			"project_id": []string{
				"required",
				"uuid",
			},
		},
		Messages: map[string][]string{
			"color": {
				"regex:The color must be valid HEX color e.g #283593",
			},
		},
	})
	return v.ValidateStruct()
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
