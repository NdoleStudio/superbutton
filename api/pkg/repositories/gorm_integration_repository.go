package repositories

import (
	"errors"
	"fmt"

	"github.com/NdoleStudio/superbutton/pkg/entities"
	"github.com/NdoleStudio/superbutton/pkg/telemetry"
	"github.com/google/uuid"
	"github.com/palantir/stacktrace"
	"gorm.io/gorm"
)

// gormIntegrationRepository is responsible for persisting integrations
type gormIntegrationRepository struct {
	logger telemetry.Logger
	tracer telemetry.Tracer
	db     *gorm.DB
}

func (repository *gormContentIntegrationRepository) getPosition(tx *gorm.DB, integration *entities.Integration) (uint, error) {
	projectIntegration := new(entities.ProjectIntegration)
	err := tx.Where("user_id = ?", integration.UserID).
		Where("project_id = ?", integration.ProjectID).
		Select("position").Order("position desc").First(projectIntegration).
		Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		msg := fmt.Sprintf("cannot fetch last inegration for project [%s] and user [%s]", integration.ProjectID, integration.UserID)
		return 0, stacktrace.Propagate(err, msg)
	}
	return projectIntegration.Position + 1, nil
}

func (repository *gormContentIntegrationRepository) updateProjectIntegration(tx *gorm.DB, integration *entities.Integration) error {
	err := tx.
		Model(&entities.ProjectIntegration{}).
		Where("integration_id = ?", integration.IntegrationID).
		Updates(map[string]interface{}{"updated_at": integration.UpdatedAt, "name": integration.Name}).
		Error
	if err != nil {
		return stacktrace.Propagate(err, fmt.Sprintf("cannot update project integration with integration ID [%s]", integration.IntegrationID))
	}
	return nil
}

func (repository *gormContentIntegrationRepository) deleteProjectIntegration(tx *gorm.DB, userID entities.UserID, integrationID uuid.UUID) error {
	err := tx.
		Where("integration_id = ?", integrationID).
		Where("user_id = ?", userID).
		Delete(&entities.ProjectIntegration{}).
		Error
	if err != nil {
		return stacktrace.Propagate(err, fmt.Sprintf("cannot delete project integration with integration ID [%s]", integrationID))
	}
	return nil
}
