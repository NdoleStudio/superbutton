package requests

import (
	"strings"

	"github.com/NdoleStudio/superbutton/pkg/entities"
	"github.com/NdoleStudio/superbutton/pkg/services"
	"github.com/google/uuid"
)

// ProjectUpdateRequest is the payload for the /projects/create endpoint
type ProjectUpdateRequest struct {
	request
	ProjectID       string `json:"project_id" swaggerignore:"true"`
	Source          string `json:"source" swaggerignore:"true"`
	Name            string `json:"name"`
	Website         string `json:"website"`
	Icon            string `json:"icon"`
	Greeting        string `json:"greeting"`
	GreetingTimeout uint   `json:"greeting_timeout"`
	Color           string `json:"color"`
}

// Sanitize the request by stripping whitespaces
func (request *ProjectUpdateRequest) Sanitize() *ProjectUpdateRequest {
	request.Name = request.sanitizeString(request.Name)
	request.Website = request.sanitizeString(request.Website)
	request.Icon = request.sanitizeString(request.Icon)
	request.Greeting = request.sanitizeString(request.Greeting)

	request.Color = strings.ToUpper(request.sanitizeString(request.Color))
	if request.Color == "" {
		request.Color = "#283593"
	}

	return request
}

// ToProjectUpdatePrams creates services.ProjectUpdateParams from ProjectUpdateRequest
func (request *ProjectUpdateRequest) ToProjectUpdatePrams(source string, userID entities.UserID) *services.ProjectUpdateParams {
	return &services.ProjectUpdateParams{
		Name:                   request.Name,
		URL:                    request.Website,
		Icon:                   request.Icon,
		Greeting:               request.Greeting,
		ProjectID:              uuid.MustParse(request.ProjectID),
		UserID:                 userID,
		Source:                 source,
		GreetingTimeoutSeconds: request.GreetingTimeout,
		Color:                  request.Color,
	}
}
