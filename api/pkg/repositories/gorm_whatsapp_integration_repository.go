package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/NdoleStudio/superbutton/pkg/entities"
	"github.com/NdoleStudio/superbutton/pkg/telemetry"
	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbgorm"
	"github.com/google/uuid"
	"github.com/palantir/stacktrace"
	"gorm.io/gorm"
)

// gormWhatsappIntegrationRepository is responsible for persisting entities.WhatsappIntegration
type gormWhatsappIntegrationRepository struct {
	logger telemetry.Logger
	tracer telemetry.Tracer
	db     *gorm.DB
}

// NewGormWhatsappIntegrationRepository creates the GORM version of the WhatsappIntegrationRepository
func NewGormWhatsappIntegrationRepository(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	db *gorm.DB,
) WhatsappIntegrationRepository {
	return &gormWhatsappIntegrationRepository{
		logger: logger.WithService(fmt.Sprintf("%T", &gormWhatsappIntegrationRepository{})),
		tracer: tracer,
		db:     db,
	}
}

func (repository *gormWhatsappIntegrationRepository) Store(ctx context.Context, integration *entities.WhatsappIntegration) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	err := crdbgorm.ExecuteTx(ctx, repository.db, nil, func(tx *gorm.DB) error {
		err := tx.WithContext(ctx).Create(integration).Error
		if err != nil {
			return stacktrace.Propagate(err, fmt.Sprintf("cannot store integration with ID [%s]", integration.ID))
		}

		projectIntegration := new(entities.ProjectIntegration)
		tx.WithContext(ctx).
			Where("user_id = ?", integration.UserID).
			Where("project_id = ?", integration.ProjectID).
			Select("order").Order("order desc").First(projectIntegration)

		err = tx.WithContext(ctx).Create(integration.NewProjectIntegration(projectIntegration.Order + 1)).Error
		if err != nil {
			return stacktrace.Propagate(err, fmt.Sprintf("cannot store project integration with ID [%s]", integration.ID))
		}

		return nil
	})
	if err != nil {
		msg := fmt.Sprintf("cannot save whatsapp integration with ID [%s] and project [%s]", integration.ID, integration.ProjectID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

func (repository *gormWhatsappIntegrationRepository) Update(ctx context.Context, integration *entities.WhatsappIntegration) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	err := crdbgorm.ExecuteTx(ctx, repository.db, nil, func(tx *gorm.DB) error {
		err := tx.WithContext(ctx).Save(integration).Error
		if err != nil {
			return stacktrace.Propagate(err, fmt.Sprintf("cannot save integration with ID [%s]", integration.ID))
		}

		err = tx.WithContext(ctx).
			Model(&entities.ProjectIntegration{}).
			Where("integration_id = ?", integration.ID).
			Updates(map[string]interface{}{"updated_at": integration.UpdatedAt, "name": integration.Name}).
			Error
		if err != nil {
			return stacktrace.Propagate(err, fmt.Sprintf("cannot update project integration with integration ID [%s]", integration.ID))
		}
		return nil
	})
	if err != nil {
		msg := fmt.Sprintf("cannot update whatsapp integration with ID [%s] and project [%s]", integration.ID, integration.ProjectID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

func (repository *gormWhatsappIntegrationRepository) Fetch(ctx context.Context, userID entities.UserID, projectID uuid.UUID) ([]*entities.WhatsappIntegration, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	var integrations []*entities.WhatsappIntegration
	err := repository.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("project_id = ?", projectID).
		Find(&integrations).
		Error
	if err != nil {
		msg := fmt.Sprintf("cannot load integrations for user with ID [%s] and project [%s]", userID, projectID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return integrations, nil
}

func (repository *gormWhatsappIntegrationRepository) Delete(ctx context.Context, userID entities.UserID, integrationID uuid.UUID) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	err := crdbgorm.ExecuteTx(ctx, repository.db, nil, func(tx *gorm.DB) error {
		err := tx.WithContext(ctx).
			Where("user_id = ?", userID).
			Where("id = ?", integrationID).
			Delete(&entities.WhatsappIntegration{}).
			Error
		if err != nil {
			return stacktrace.Propagate(err, fmt.Sprintf("cannot delete whatsapp integration with ID [%s]", integrationID))
		}

		err = tx.WithContext(ctx).
			Where("integration_id = ?", integrationID).
			Where("user_id = ?", userID).
			Delete(&entities.ProjectIntegration{}).
			Error
		if err != nil {
			return stacktrace.Propagate(err, fmt.Sprintf("cannot delete project integration with integration ID [%s]", integrationID))
		}

		return nil
	})
	if err != nil {
		msg := fmt.Sprintf("cannot delete integration for user with ID [%s] and integration [%s]", userID, integrationID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

func (repository *gormWhatsappIntegrationRepository) Load(ctx context.Context, userID entities.UserID, integrationID uuid.UUID) (*entities.WhatsappIntegration, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	integration := new(entities.WhatsappIntegration)
	err := repository.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("id = ?", integrationID).
		First(integration).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		msg := fmt.Sprintf("whatsapp integration with ID [%s] for user [%s] does not exist", integrationID, userID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.PropagateWithCode(err, ErrCodeNotFound, msg))
	}

	if err != nil {
		msg := fmt.Sprintf("cannot load integration with ID [%s] and user [%s]", integrationID, userID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return integration, nil
}
