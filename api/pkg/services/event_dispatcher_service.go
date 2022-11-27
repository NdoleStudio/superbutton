package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/NdoleStudio/superbutton/pkg/listeners"
	"github.com/NdoleStudio/superbutton/pkg/queue"
	"github.com/NdoleStudio/superbutton/pkg/repositories"
	"github.com/NdoleStudio/superbutton/pkg/telemetry"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/palantir/stacktrace"
)

// EventDispatcher dispatches a new event
type EventDispatcher struct {
	logger      telemetry.Logger
	tracer      telemetry.Tracer
	queue       queue.Client
	consumerURL string
	listeners   map[string][]listeners.EventListener
	repository  repositories.EventRepository
}

// NewEventDispatcher creates a new EventDispatcher
func NewEventDispatcher(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	repository repositories.EventRepository,
	queue queue.Client,
	consumerURL string,
) (dispatcher *EventDispatcher) {
	return &EventDispatcher{
		logger:      logger,
		listeners:   make(map[string][]listeners.EventListener),
		tracer:      tracer,
		queue:       queue,
		consumerURL: consumerURL,
		repository:  repository,
	}
}

// Dispatch a new event by adding it to the queue to be processed async
func (dispatcher *EventDispatcher) Dispatch(ctx context.Context, event *cloudevents.Event) error {
	ctx, span := dispatcher.tracer.Start(ctx)
	defer span.End()

	ctxLogger := dispatcher.tracer.CtxLogger(dispatcher.logger, span)

	if err := event.Validate(); err != nil {
		msg := fmt.Sprintf("cannot dispatch event with ID [%s] and type [%s] because it is invalid", event.ID(), event.Type())
		return dispatcher.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	if err := dispatcher.repository.Save(ctx, event); err != nil {
		ctxLogger.Error(stacktrace.Propagate(err, fmt.Sprintf("cannot store [%s] event with id [%s]", event.Type(), event.ID())))
	}

	task, err := dispatcher.createTask(event)
	if err != nil {
		msg := fmt.Sprintf("cannot create push queue task for event with ID [%s] and type [%s]", event.ID(), event.Type())
		return dispatcher.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	taskID, err := dispatcher.queue.Enqueue(ctx, task)
	if err != nil {
		msg := fmt.Sprintf("cannot add event with ID [%s] and type [%s] to producer", event.ID(), event.Type())
		return dispatcher.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	ctxLogger.Info(fmt.Sprintf("push queue task enqueued with ID [%s]", taskID))
	return nil
}

// Subscribe a listener to an event
func (dispatcher *EventDispatcher) Subscribe(eventType string, listener listeners.EventListener) {
	if _, ok := dispatcher.listeners[eventType]; !ok {
		dispatcher.listeners[eventType] = []listeners.EventListener{}
	}

	dispatcher.listeners[eventType] = append(dispatcher.listeners[eventType], listener)
}

// Publish an event to subscribers
func (dispatcher *EventDispatcher) Publish(ctx context.Context, event cloudevents.Event) {
	ctx, span := dispatcher.tracer.Start(ctx)
	defer span.End()

	ctxLogger := dispatcher.tracer.CtxLogger(dispatcher.logger, span)

	subscribers, ok := dispatcher.listeners[event.Type()]
	if !ok {
		ctxLogger.Info(fmt.Sprintf("no listener is configured for event type [%s]", event.Type()))
		return
	}

	var wg sync.WaitGroup
	for _, sub := range subscribers {
		wg.Add(1)
		go func(ctx context.Context, sub listeners.EventListener) {
			if err := sub(ctx, event); err != nil {
				msg := fmt.Sprintf("subscriber [%T] cannot handle event [%s]", sub, event.Type())
				ctxLogger.Error(stacktrace.Propagate(err, msg))
			}
			wg.Done()
		}(ctx, sub)
	}

	wg.Wait()
}

func (dispatcher *EventDispatcher) createTask(event *cloudevents.Event) (*queue.Task, error) {
	eventContent, err := json.Marshal(event)
	if err != nil {
		return nil, stacktrace.Propagate(err, fmt.Sprintf("cannot marshall [%T] with ID [%s]", event, event.ID()))
	}

	return &queue.Task{
		Method: http.MethodPost,
		URL:    dispatcher.consumerURL,
		Body:   eventContent,
	}, nil
}
