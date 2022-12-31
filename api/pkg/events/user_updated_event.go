package events

import (
	"time"

	"github.com/NdoleStudio/superbutton/pkg/entities"
)

// UserUpdated is raised when a user is created
const UserUpdated = "user.updated"

// UserUpdatedPayload stores the data for the UserUpdated event
type UserUpdatedPayload struct {
	UserID           entities.UserID           `json:"user_id"`
	SubscriptionName entities.SubscriptionName `json:"subscription_name"`
	UserUpdatedAt    time.Time                 `json:"user_updated_at"`
	UserEmail        string                    `json:"user_email"`
}
