package validators

import (
	"context"
	"fmt"
	"net/url"

	"github.com/NdoleStudio/superbutton/pkg/requests"
	"github.com/NdoleStudio/superbutton/pkg/telemetry"
	"github.com/thedevsaddam/govalidator"
)

// ContentIntegrationHandlerValidator validates models used in handlers.ContentIntegrationHandler
type ContentIntegrationHandlerValidator struct {
	logger telemetry.Logger
	tracer telemetry.Tracer
}

// NewContentIntegrationHandlerValidator creates a new handlers.ContentIntegrationHandler validator
func NewContentIntegrationHandlerValidator(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
) (v *ContentIntegrationHandlerValidator) {
	return &ContentIntegrationHandlerValidator{
		logger: logger.WithService(fmt.Sprintf("%T", v)),
		tracer: tracer,
	}
}

func (validator *ContentIntegrationHandlerValidator) ValidateUpdate(ctx context.Context, request *requests.ContentIntegrationUpdateRequest) url.Values {
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
			"title": []string{
				"required",
				"min:1",
				"max:50",
			},
			"text": []string{
				"required",
				"min:1",
				"max:1000",
			},
			"summary": []string{
				"required",
				"min:1",
				"max:100",
			},
			"integrationID": []string{
				"required",
				"uuid",
			},
		},
	})
	return v.ValidateStruct()
}

func (validator *ContentIntegrationHandlerValidator) ValidateCreate(ctx context.Context, request *requests.ContentIntegrationCreateRequest) url.Values {
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
			"title": []string{
				"required",
				"min:1",
				"max:50",
			},
			"text": []string{
				"required",
				"min:1",
				"max:1000",
			},
			"summary": []string{
				"required",
				"min:1",
				"max:100",
			},
			"projectID": []string{
				"required",
				"uuid",
			},
		},
	})
	return v.ValidateStruct()
}
