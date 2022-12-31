package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/NdoleStudio/superbutton/pkg/entities"
	"github.com/NdoleStudio/superbutton/pkg/telemetry"
	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbgorm"
	"github.com/palantir/stacktrace"
	"gorm.io/gorm"
)

// gormUserRepository is responsible for persisting entities.User
type gormUserRepository struct {
	logger telemetry.Logger
	tracer telemetry.Tracer
	db     *gorm.DB
}

// NewGormUserRepository creates the GORM version of the UserRepository
func NewGormUserRepository(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	db *gorm.DB,
) UserRepository {
	return &gormUserRepository{
		logger: logger.WithService(fmt.Sprintf("%T", &gormUserRepository{})),
		tracer: tracer,
		db:     db,
	}
}

func (repository *gormUserRepository) LoadBySubscriptionID(ctx context.Context, subscriptionID string) (*entities.User, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	user := new(entities.User)
	err := repository.db.WithContext(ctx).
		Where("subscription_id = ?", subscriptionID).
		First(user).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		msg := fmt.Sprintf("user with subscriptionID [%s] does not exist", subscriptionID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.PropagateWithCode(err, ErrCodeNotFound, msg))
	}

	if err != nil {
		msg := fmt.Sprintf("cannot load subscription with ID [%s]", subscriptionID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return user, nil
}

func (repository *gormUserRepository) Store(ctx context.Context, user *entities.User) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	if err := repository.db.WithContext(ctx).Create(user).Error; err != nil {
		msg := fmt.Sprintf("cannot save user with ID [%s]", user.ID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

func (repository *gormUserRepository) Update(ctx context.Context, user *entities.User) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	if err := repository.db.WithContext(ctx).Save(user).Error; err != nil {
		msg := fmt.Sprintf("cannot update user with ID [%s]", user.ID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

func (repository *gormUserRepository) Load(ctx context.Context, userID entities.UserID) (*entities.User, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	user := new(entities.User)
	err := repository.db.WithContext(ctx).First(user, userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		msg := fmt.Sprintf("user with ID [%s] does not exist", user.ID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.PropagateWithCode(err, ErrCodeNotFound, msg))
	}

	if err != nil {
		msg := fmt.Sprintf("cannot load user with ID [%s]", userID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return user, nil
}

func (repository *gormUserRepository) LoadOrStore(ctx context.Context, authUser entities.AuthUser) (user *entities.User, created bool, err error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	err = crdbgorm.ExecuteTx(ctx, repository.db, nil, func(tx *gorm.DB) error {
		user = new(entities.User)
		err = tx.First(user, authUser.ID).Error
		if err == nil {
			return nil
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			msg := fmt.Sprintf("cannot check if user exists with ID [%s]", authUser.ID)
			return stacktrace.Propagate(err, msg)
		}

		user = &entities.User{
			ID:               authUser.ID,
			Email:            authUser.Email,
			SubscriptionName: entities.SubscriptionNameFree,
			Name:             authUser.Name,
			CreatedAt:        time.Now().UTC(),
			UpdatedAt:        time.Now().UTC(),
		}
		created = true

		return tx.Where(entities.User{ID: user.ID}).Create(user).Error
	})

	if err != nil {
		msg := fmt.Sprintf("cannot lod or create user from auth user [%+#v]", authUser)
		return user, created, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return user, created, nil
}
