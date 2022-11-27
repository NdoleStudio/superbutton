package validators

import (
	"context"
	"fmt"
	"net/url"

	"github.com/NdoleStudio/superbutton/pkg/requests"
	"github.com/NdoleStudio/superbutton/pkg/telemetry"
	"github.com/thedevsaddam/govalidator"
)

// WhatsappIntegrationHandlerValidator validates models used in handlers.WhatsappIntegrationHandler
type WhatsappIntegrationHandlerValidator struct {
	logger telemetry.Logger
	tracer telemetry.Tracer
}

// NewWhatsappIntegrationHandlerValidator creates a new handlers.WhatsappIntegrationHandlerValidator validator
func NewWhatsappIntegrationHandlerValidator(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
) (v *WhatsappIntegrationHandlerValidator) {
	return &WhatsappIntegrationHandlerValidator{
		logger: logger.WithService(fmt.Sprintf("%T", v)),
		tracer: tracer,
	}
}

func (validator *WhatsappIntegrationHandlerValidator) ValidateUpdate(ctx context.Context, request *requests.WhatsappIntegrationUpdateRequest) url.Values {
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
			"text": []string{
				"required",
				"min:1",
				"max:30",
			},
			"phone_number": []string{
				"required",
				phoneNumberRule,
			},
			"integrationID": []string{
				"required",
				"uuid",
			},
		},
	})
	return v.ValidateStruct()
}

func (validator *WhatsappIntegrationHandlerValidator) ValidateCreate(ctx context.Context, request *requests.WhatsappIntegrationCreateRequest) url.Values {
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
			"text": []string{
				"required",
				"min:1",
				"max:30",
			},
			"phone_number": []string{
				"required",
				"min:1",
				"max:30",
			},
			"projectID": []string{
				"required",
				"uuid",
			},
		},
	})
	return v.ValidateStruct()
}
