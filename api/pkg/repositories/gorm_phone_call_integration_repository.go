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

// gormPhoneCallIntegrationRepository is responsible for persisting entities.PhoneCallIntegration
type gormPhoneCallIntegrationRepository struct {
	gormIntegrationRepository
}

// NewGormPhoneCallIntegrationRepository creates the GORM version of the PhoneCallIntegrationRepository
func NewGormPhoneCallIntegrationRepository(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	db *gorm.DB,
) PhoneCallIntegrationRepository {
	return &gormPhoneCallIntegrationRepository{
		gormIntegrationRepository{
			logger: logger.WithService(fmt.Sprintf("%T", &gormPhoneCallIntegrationRepository{})),
			tracer: tracer,
			db:     db,
		},
	}
}

func (repository *gormPhoneCallIntegrationRepository) FetchMultiple(ctx context.Context, userID entities.UserID, IDs []uuid.UUID) ([]*entities.PhoneCallIntegration, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	var integrations []*entities.PhoneCallIntegration
	err := repository.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("id IN ?", IDs).
		Find(&integrations).
		Error
	if err != nil {
		msg := fmt.Sprintf("cannot load phone call integrations for user with ID [%s] and projects [%+#v]", userID, IDs)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return integrations, nil
}

func (repository *gormPhoneCallIntegrationRepository) Store(ctx context.Context, integration *entities.PhoneCallIntegration) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	err := crdbgorm.ExecuteTx(ctx, repository.db, nil, func(tx *gorm.DB) error {
		err := tx.Create(integration).Error
		if err != nil {
			return stacktrace.Propagate(err, fmt.Sprintf("cannot store [phone call] integration with ID [%s]", integration.ID))
		}

		position, err := repository.getPosition(tx, integration.Integration())
		if err != nil {
			msg := fmt.Sprintf("cannot fetch last integration for project [%s] and user [%s]", integration.ProjectID, integration.UserID)
			return stacktrace.Propagate(err, msg)
		}

		err = tx.Create(integration.NewProjectIntegration(position)).Error
		if err != nil {
			return stacktrace.Propagate(err, fmt.Sprintf("cannot store project integration with ID [%s]", integration.ID))
		}

		return nil
	})
	if err != nil {
		msg := fmt.Sprintf("cannot save content integration with ID [%s] and project [%s]", integration.ID, integration.ProjectID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

func (repository *gormPhoneCallIntegrationRepository) Update(ctx context.Context, integration *entities.PhoneCallIntegration) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	err := crdbgorm.ExecuteTx(ctx, repository.db, nil, func(tx *gorm.DB) error {
		if err := tx.Save(integration).Error; err != nil {
			return stacktrace.Propagate(err, fmt.Sprintf("cannot save integration with ID [%s]", integration.ID))
		}
		if err := repository.updateProjectIntegration(tx, integration.Integration()); err != nil {
			return stacktrace.Propagate(err, fmt.Sprintf("cannot update project integration for integrtion [%s]", integration.ID))
		}
		return nil
	})
	if err != nil {
		msg := fmt.Sprintf("cannot content integration with ID [%s] and project [%s]", integration.ID, integration.ProjectID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

func (repository *gormPhoneCallIntegrationRepository) Fetch(ctx context.Context, userID entities.UserID, projectID uuid.UUID) ([]*entities.PhoneCallIntegration, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	var integrations []*entities.PhoneCallIntegration
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

func (repository *gormPhoneCallIntegrationRepository) Delete(ctx context.Context, userID entities.UserID, integrationID uuid.UUID) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	err := crdbgorm.ExecuteTx(ctx, repository.db, nil, func(tx *gorm.DB) error {
		err := tx.
			Where("user_id = ?", userID).
			Where("id = ?", integrationID).
			Delete(&entities.ContentIntegration{}).
			Error
		if err != nil {
			return stacktrace.Propagate(err, fmt.Sprintf("cannot delete content integration with ID [%s]", integrationID))
		}

		err = repository.deleteProjectIntegration(tx, userID, integrationID)
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

func (repository *gormPhoneCallIntegrationRepository) Load(ctx context.Context, userID entities.UserID, integrationID uuid.UUID) (*entities.PhoneCallIntegration, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	integration := new(entities.PhoneCallIntegration)
	err := repository.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("id = ?", integrationID).
		First(integration).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		msg := fmt.Sprintf("content integration with ID [%s] for user [%s] does not exist", integrationID, userID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.PropagateWithCode(err, ErrCodeNotFound, msg))
	}

	if err != nil {
		msg := fmt.Sprintf("cannot load content integration with ID [%s] and user [%s]", integrationID, userID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return integration, nil
}
