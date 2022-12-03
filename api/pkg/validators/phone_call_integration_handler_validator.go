package validators

import (
	"context"
	"fmt"
	"net/url"

	"github.com/NdoleStudio/superbutton/pkg/requests"
	"github.com/NdoleStudio/superbutton/pkg/telemetry"
	"github.com/thedevsaddam/govalidator"
)

// PhoneCallIntegrationHandlerValidator validates models used in handlers.PhoneCallIntegrationHandler
type PhoneCallIntegrationHandlerValidator struct {
	logger telemetry.Logger
	tracer telemetry.Tracer
}

// NewPhoneCallIntegrationHandlerValidator creates a new handlers.PhoneCallIntegrationHandler validator
func NewPhoneCallIntegrationHandlerValidator(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
) (v *PhoneCallIntegrationHandlerValidator) {
	return &PhoneCallIntegrationHandlerValidator{
		logger: logger.WithService(fmt.Sprintf("%T", v)),
		tracer: tracer,
	}
}

func (validator *PhoneCallIntegrationHandlerValidator) ValidateUpdate(ctx context.Context, request *requests.PhoneCallIntegrationUpdateRequest) url.Values {
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

func (validator *PhoneCallIntegrationHandlerValidator) ValidateCreate(ctx context.Context, request *requests.PhoneCallIntegrationCreateRequest) url.Values {
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
