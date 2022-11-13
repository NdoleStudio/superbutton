package events

import (
	"time"

	"github.com/NdoleStudio/superbutton/pkg/entities"
)

// UserCreated is raised when a user is created
const UserCreated = "user.created"

// UserCreatedPayload stores the data for the user created event
type UserCreatedPayload struct {
	UserID        entities.UserID `json:"user_id"`
	UserCreatedAt time.Time       `json:"user_created_at"`
	UserName      string          `json:"user_name"`
	UserEmail     string          `json:"user_email"`
}
