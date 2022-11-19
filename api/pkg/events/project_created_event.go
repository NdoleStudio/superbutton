package events

import (
	"time"

	"github.com/google/uuid"

	"github.com/NdoleStudio/superbutton/pkg/entities"
)

// ProjectCreated is raised when a user is created
const ProjectCreated = "project.created"

// ProjectCreatedPayload stores the data for the ProjectCreated event
type ProjectCreatedPayload struct {
	UserID           entities.UserID `json:"user_id"`
	ProjectID        uuid.UUID       `json:"project_id"`
	ProjectName      string          `json:"project_name"`
	ProjectURL       string          `json:"project_url"`
	ProjectCreatedAt time.Time       `json:"project_created_at"`
}
