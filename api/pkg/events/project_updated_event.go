package events

import (
	"time"

	"github.com/google/uuid"

	"github.com/NdoleStudio/superbutton/pkg/entities"
)

// ProjectUpdated is raised when a user is created
const ProjectUpdated = "project.updated"

// ProjectUpdatedPayload stores the data for the ProjectUpdated event
type ProjectUpdatedPayload struct {
	UserID                 entities.UserID `json:"user_id"`
	ProjectID              uuid.UUID       `json:"project_id"`
	ProjectName            string          `json:"project_name"`
	ProjectURL             string          `json:"project_url"`
	ProjectIcon            string          `json:"project_icon"`
	ProjectGreeting        string          `json:"project_greeting"`
	ProjectColor           string          `json:"project_color"`
	ProjectGreetingTimeout uint            `json:"project_greeting_timeout"`
	ProjectUpdatedAt       time.Time       `json:"project_updated_at"`
}
