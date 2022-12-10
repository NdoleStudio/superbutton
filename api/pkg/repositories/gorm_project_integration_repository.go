package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbgorm"

	"github.com/NdoleStudio/superbutton/pkg/entities"
	"github.com/NdoleStudio/superbutton/pkg/telemetry"
	"github.com/google/uuid"
	"github.com/palantir/stacktrace"
	"gorm.io/gorm"
)

// gormProjectIntegrationRepository is responsible for persisting entities.WhatsappIntegration
type gormProjectIntegrationRepository struct {
	logger telemetry.Logger
	tracer telemetry.Tracer
	db     *gorm.DB
}

// NewGormProjectIntegrationRepository creates the GORM version of the ProjectIntegrationRepository
func NewGormProjectIntegrationRepository(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	db *gorm.DB,
) ProjectIntegrationRepository {
	return &gormProjectIntegrationRepository{
		logger: logger.WithService(fmt.Sprintf("%T", &gormProjectIntegrationRepository{})),
		tracer: tracer,
		db:     db,
	}
}

func (repository *gormProjectIntegrationRepository) Fetch(ctx context.Context, userID entities.UserID, projectID uuid.UUID) ([]*entities.ProjectIntegration, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	var integrations []*entities.ProjectIntegration
	err := repository.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("project_id = ?", projectID).
		Order("position asc").
		Find(&integrations).
		Error
	if err != nil {
		msg := fmt.Sprintf("cannot load integrations for user with ID [%s] and project [%s]", userID, projectID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return integrations, nil
}

func (repository *gormProjectIntegrationRepository) UpdatePositions(ctx context.Context, userID entities.UserID, integrationIDs []uuid.UUID) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	updatedAt := time.Now().UTC()
	err := crdbgorm.ExecuteTx(ctx, repository.db, nil, func(tx *gorm.DB) error {
		for index, integrationID := range integrationIDs {
			err := tx.
				Model(&entities.ProjectIntegration{}).
				Where("integration_id = ?", integrationID).
				Where("user_id = ?", userID).
				Updates(map[string]interface{}{"updated_at": updatedAt, "position": index}).
				Error
			if err != nil {
				msg := fmt.Sprintf("cannot update integration [%s] with position [%d] for user [%s]", integrationID, index, userID)
				return stacktrace.Propagate(err, msg)
			}
		}
		return nil
	})
	if err != nil {
		msg := fmt.Sprintf("cannot update integration positions for user [%s] with ID's [%+#v]", userID, integrationIDs)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}
	return nil
}
