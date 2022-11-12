package events

import (
	"time"

	"github.com/NdoleStudio/superbutton/pkg/entities"
)

// UserCreated is raised when a user is created
const UserCreated = "user.created"

// UserCreatedPayload stores the data for the user created event
type UserCreatedPayload struct {
	ID        entities.UserID `json:"id"`
	CreatedAt time.Time       `json:"created_at"`
	Name      *string         `json:"name"`
	Email     string          `json:"email"`
}
