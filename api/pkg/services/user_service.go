package services

import (
	"context"
	"fmt"

	"github.com/NdoleStudio/superbutton/pkg/entities"
	"github.com/NdoleStudio/superbutton/pkg/events"
	"github.com/NdoleStudio/superbutton/pkg/repositories"
	"github.com/NdoleStudio/superbutton/pkg/telemetry"
	"github.com/palantir/stacktrace"
)

// User is responsible for managing entities.User
type User struct {
	service
	logger          telemetry.Logger
	tracer          telemetry.Tracer
	repository      repositories.UserRepository
	eventDispatcher *EventDispatcher
}

// NewUserService creates a new UserService
func NewUserService(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	eventDispatcher *EventDispatcher,
	repository repositories.UserRepository,
) (s *User) {
	return &User{
		logger:          logger.WithService(fmt.Sprintf("%T", s)),
		tracer:          tracer,
		eventDispatcher: eventDispatcher,
		repository:      repository,
	}
}

// Get fetches or creates an entities.User
func (service *User) Get(ctx context.Context, source string, authUser entities.AuthUser) (*entities.User, error) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	user, created, err := service.repository.LoadOrStore(ctx, authUser)
	if err != nil {
		msg := fmt.Sprintf("could not get [%T] with from [%+#v]", user, authUser)
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	if created {
		service.dispatchUserCreatedEvent(ctx, source, user)
	}

	return user, nil
}

func (service *User) dispatchUserCreatedEvent(ctx context.Context, source string, user *entities.User) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	ctxLogger := service.tracer.CtxLogger(service.logger, span)

	event, err := service.createEvent(events.UserCreated, source, &events.UserCreatedPayload{
		UserID:        user.ID,
		UserCreatedAt: user.CreatedAt,
		UserName:      user.Name,
		UserEmail:     user.Email,
	})
	if err != nil {
		msg := fmt.Sprintf("cannot created [%s] event for user [%s]", events.UserCreated, user.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return
	}

	if err = service.eventDispatcher.Dispatch(ctx, event); err != nil {
		msg := fmt.Sprintf("cannot dispatch [%s] event for user [%s]", events.UserCreated, user.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return
	}
}
