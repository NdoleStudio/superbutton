package repositories

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// EventRepository is responsible for persisting cloudevents.Event
type EventRepository interface {
	// Create a new entities.Message
	Create(ctx context.Context, event cloudevents.Event) error

	// Save a new entities.Message
	Save(ctx context.Context, event cloudevents.Event) error

	// FetchAll returns all cloudevents.Event ordered by time in ascending order
	FetchAll(ctx context.Context) (*[]cloudevents.Event, error)
}
