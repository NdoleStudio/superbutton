package queue

import (
	"context"
	"fmt"
	"net/http"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"github.com/NdoleStudio/superbutton/pkg/telemetry"
	"github.com/palantir/stacktrace"
)

type googlePushQueue struct {
	logger    telemetry.Logger
	tracer    telemetry.Tracer
	client    *cloudtasks.Client
	queueName string
	authEmail string
}

// NewGooglePushQueue creates a new googlePushQueue
func NewGooglePushQueue(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	client *cloudtasks.Client,
	queueName string,
	authEmail string,
) Client {
	return &googlePushQueue{
		client:    client,
		tracer:    tracer,
		logger:    logger,
		queueName: queueName,
		authEmail: authEmail,
	}
}

// Enqueue a task to the queue
func (queue *googlePushQueue) Enqueue(ctx context.Context, task *Task) (queueID string, err error) {
	ctx, span := queue.tracer.Start(ctx)
	defer span.End()

	ctxLogger := queue.tracer.CtxLogger(queue.logger, span)

	req := &cloudtaskspb.CreateTaskRequest{
		Parent: queue.queueName,
		Task: &cloudtaskspb.Task{
			MessageType: &cloudtaskspb.Task_HttpRequest{
				HttpRequest: &cloudtaskspb.HttpRequest{
					HttpMethod: queue.httpMethodToProtoHTTPMethod(task.Method),
					Url:        task.URL,
					AuthorizationHeader: &cloudtaskspb.HttpRequest_OidcToken{
						OidcToken: &cloudtaskspb.OidcToken{
							ServiceAccountEmail: queue.authEmail,
						},
					},
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
			},
		},
	}

	// Add a payload message if one is present.
	req.Task.GetHttpRequest().Body = task.Body

	queueTask, err := queue.client.CreateTask(ctx, req)
	if err != nil {
		msg := fmt.Sprintf("cannot schedule task %s to URL: %s", string(task.Body), task.URL)
		return queueID, queue.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	ctxLogger.Info(fmt.Sprintf(
		"item added to [%s] queue with id [%s] and schedule [%s]",
		queue.queueName,
		queueTask.Name,
		queueTask.ScheduleTime,
	))

	return queueTask.Name, nil
}

func (queue *googlePushQueue) httpMethodToProtoHTTPMethod(httpMethod string) cloudtaskspb.HttpMethod {
	method, ok := map[string]cloudtaskspb.HttpMethod{
		http.MethodGet:  cloudtaskspb.HttpMethod_GET,
		http.MethodPost: cloudtaskspb.HttpMethod_POST,
		http.MethodPut:  cloudtaskspb.HttpMethod_PUT,
	}[httpMethod]

	if !ok {
		return cloudtaskspb.HttpMethod_POST
	}

	return method
}
