package repositories

import (
	"context"
	"fmt"

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
		Find(&integrations).
		Error
	if err != nil {
		msg := fmt.Sprintf("cannot load integrations for user with ID [%s] and project [%s]", userID, projectID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return integrations, nil
}
