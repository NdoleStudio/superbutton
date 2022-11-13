package responses

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// UserResponse is the payload containing entities.User
type UserResponse struct {
	response
	Data cloudevents.Event `json:"data"`
}
