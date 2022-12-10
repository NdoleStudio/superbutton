package requests

import (
	"github.com/google/uuid"
)

// ProjectIntegrationsUpdateRequest is the payload for the /projects/{projectID}/integrations update endpoint
type ProjectIntegrationsUpdateRequest struct {
	Order []uuid.UUID `json:"order"`
}
