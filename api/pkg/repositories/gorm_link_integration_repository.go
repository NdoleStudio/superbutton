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

// gormLinkIntegrationRepository is responsible for persisting entities.LinkIntegration
type gormLinkIntegrationRepository struct {
	gormIntegrationRepository
}

// NewGormLinkIntegrationRepository creates the GORM version of the LinkIntegrationRepository
func NewGormLinkIntegrationRepository(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	db *gorm.DB,
) LinkIntegrationRepository {
	return &gormLinkIntegrationRepository{
		gormIntegrationRepository{
			integrationType: entities.IntegrationTypeLink,
			logger:          logger.WithService(fmt.Sprintf("%T", &gormLinkIntegrationRepository{})),
			tracer:          tracer,
			db:              db,
		},
	}
}

func (repository *gormLinkIntegrationRepository) FetchMultiple(ctx context.Context, userID entities.UserID, IDs []uuid.UUID) ([]*entities.LinkIntegration, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	var integrations []*entities.LinkIntegration
	err := repository.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("id IN ?", IDs).
		Find(&integrations).
		Error
	if err != nil {
		msg := fmt.Sprintf("cannot load [%s] integrations for user with ID [%s] and projects [%+#v]", repository.integrationType, userID, IDs)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return integrations, nil
}

func (repository *gormLinkIntegrationRepository) Store(ctx context.Context, integration *entities.LinkIntegration) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	err := crdbgorm.ExecuteTx(ctx, repository.db, nil, func(tx *gorm.DB) error {
		err := tx.Create(integration).Error
		if err != nil {
			return stacktrace.Propagate(err, fmt.Sprintf("cannot store [%s] integration with ID [%s]", repository.integrationType, integration.ID))
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
		msg := fmt.Sprintf("cannot save [%s] integration with ID [%s] and project [%s]", repository.integrationType, integration.ID, integration.ProjectID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

func (repository *gormLinkIntegrationRepository) Update(ctx context.Context, integration *entities.LinkIntegration) error {
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
		msg := fmt.Sprintf("cannot [%s] integration with ID [%s] and project [%s]", repository.integrationType, integration.ID, integration.ProjectID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

func (repository *gormLinkIntegrationRepository) Fetch(ctx context.Context, userID entities.UserID, projectID uuid.UUID) ([]*entities.LinkIntegration, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	var integrations []*entities.LinkIntegration
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

func (repository *gormLinkIntegrationRepository) Delete(ctx context.Context, userID entities.UserID, integrationID uuid.UUID) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	err := crdbgorm.ExecuteTx(ctx, repository.db, nil, func(tx *gorm.DB) error {
		err := tx.
			Where("user_id = ?", userID).
			Where("id = ?", integrationID).
			Delete(&entities.LinkIntegration{}).
			Error
		if err != nil {
			return stacktrace.Propagate(err, fmt.Sprintf("cannot delete [%s] integration with ID [%s]", repository.integrationType, integrationID))
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

func (repository *gormLinkIntegrationRepository) Load(ctx context.Context, userID entities.UserID, integrationID uuid.UUID) (*entities.LinkIntegration, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	integration := new(entities.LinkIntegration)
	err := repository.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("id = ?", integrationID).
		First(integration).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		msg := fmt.Sprintf("[%s] integration with ID [%s] for user [%s] does not exist", repository.integrationType, integrationID, userID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.PropagateWithCode(err, ErrCodeNotFound, msg))
	}

	if err != nil {
		msg := fmt.Sprintf("cannot load [%s] integration with ID [%s] and user [%s]", repository.integrationType, integrationID, userID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return integration, nil
}
