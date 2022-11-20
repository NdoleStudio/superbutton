package requests

import (
	"github.com/NdoleStudio/superbutton/pkg/entities"
	"github.com/NdoleStudio/superbutton/pkg/services"
)

// ProjectCreateRequest is the payload for the /projects/create endpoint
type ProjectCreateRequest struct {
	request
	Name    string `json:"name"`
	Website string `json:"website"`
}

// Sanitize the request by stripping whitespaces
func (request *ProjectCreateRequest) Sanitize() *ProjectCreateRequest {
	request.Name = request.sanitizeString(request.Name)
	request.Website = request.sanitizeString(request.Website)
	return request
}

// ToProjectCreateParams creates services.ProjectCreateParams from ProjectCreateRequest
func (request *ProjectCreateRequest) ToProjectCreateParams(source string, userID entities.UserID) *services.ProjectCreateParams {
	return &services.ProjectCreateParams{
		Name:   request.Name,
		Source: source,
		URL:    request.baseURL(request.Website),
		UserID: userID,
	}
}
