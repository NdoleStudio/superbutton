package events

import (
	"time"

	"github.com/google/uuid"

	"github.com/NdoleStudio/superbutton/pkg/entities"
)

// ProjectDeleted is raised when a user is created
const ProjectDeleted = "project.deleted"

// ProjectDeletedPayload stores the data for the ProjectDeleted event
type ProjectDeletedPayload struct {
	UserID           entities.UserID `json:"user_id"`
	ProjectID        uuid.UUID       `json:"project_id"`
	ProjectDeletedAt time.Time       `json:"project_deleted_at"`
}
