package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/NdoleStudio/superbutton/pkg/entities"
	"github.com/NdoleStudio/superbutton/pkg/telemetry"
	"github.com/google/uuid"
	"github.com/palantir/stacktrace"
	"gorm.io/gorm"
)

// gormProjectRepository is responsible for persisting entities.Project
type gormProjectRepository struct {
	logger telemetry.Logger
	tracer telemetry.Tracer
	db     *gorm.DB
}

// NewGormProjectRepository creates the GORM version of the ProjectRepository
func NewGormProjectRepository(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	db *gorm.DB,
) ProjectRepository {
	return &gormProjectRepository{
		logger: logger.WithService(fmt.Sprintf("%T", &gormProjectRepository{})),
		tracer: tracer,
		db:     db,
	}
}

func (repository *gormProjectRepository) Store(ctx context.Context, project *entities.Project) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	if err := repository.db.WithContext(ctx).Create(project).Error; err != nil {
		msg := fmt.Sprintf("cannot save project with ID [%s]", project.ID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

func (repository *gormProjectRepository) Update(ctx context.Context, project *entities.Project) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	if err := repository.db.WithContext(ctx).Save(project).Error; err != nil {
		msg := fmt.Sprintf("cannot update project with ID [%s]", project.ID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

func (repository *gormProjectRepository) Fetch(ctx context.Context, userID entities.UserID) ([]*entities.Project, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	var projects []*entities.Project
	err := repository.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&projects).
		Error
	if err != nil {
		msg := fmt.Sprintf("cannot load projects for user with ID [%s]", userID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return projects, nil
}

func (repository *gormProjectRepository) Load(ctx context.Context, userID entities.UserID, projectID uuid.UUID) (*entities.Project, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	project := new(entities.Project)
	err := repository.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("id = ?", projectID).
		First(project).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		msg := fmt.Sprintf("project with ID [%s] does not exist", projectID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.PropagateWithCode(err, ErrCodeNotFound, msg))
	}

	if err != nil {
		msg := fmt.Sprintf("cannot load project with ID [%s]", projectID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return project, nil
}
