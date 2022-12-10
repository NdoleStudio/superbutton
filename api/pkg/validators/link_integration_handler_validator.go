package validators

import (
	"context"
	"fmt"
	"net/url"

	"github.com/NdoleStudio/superbutton/pkg/requests"
	"github.com/NdoleStudio/superbutton/pkg/telemetry"
	"github.com/thedevsaddam/govalidator"
)

// LinkIntegrationHandlerValidator validates models used in handlers.LinkIntegrationHandler
type LinkIntegrationHandlerValidator struct {
	logger telemetry.Logger
	tracer telemetry.Tracer
}

// NewLinkIntegrationHandlerValidator creates a new handlers.LinkIntegrationHandler validator
func NewLinkIntegrationHandlerValidator(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
) (v *LinkIntegrationHandlerValidator) {
	return &LinkIntegrationHandlerValidator{
		logger: logger.WithService(fmt.Sprintf("%T", v)),
		tracer: tracer,
	}
}

func (validator *LinkIntegrationHandlerValidator) ValidateUpdate(ctx context.Context, request *requests.LinkIntegrationUpdateRequest) url.Values {
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
			"icon": []string{
				"required",
				"in:link,documentation,mail",
			},
			"text": []string{
				"required",
				"min:1",
				"max:30",
			},
			"website": []string{
				"required",
				"url",
				"max:255",
			},
			"integrationID": []string{
				"required",
				"uuid",
			},
		},
	})
	return v.ValidateStruct()
}

func (validator *LinkIntegrationHandlerValidator) ValidateCreate(ctx context.Context, request *requests.LinkIntegrationCreateRequest) url.Values {
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
			"icon": []string{
				"required",
				"in:link,documentation,mail",
			},
			"text": []string{
				"required",
				"min:1",
				"max:30",
			},
			"website": []string{
				"required",
				"url",
				"max:255",
			},
			"projectID": []string{
				"required",
				"uuid",
			},
		},
	})
	return v.ValidateStruct()
}
