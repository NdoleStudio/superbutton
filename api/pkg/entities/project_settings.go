package entities

import "github.com/google/uuid"

// ProjectSettings fetches all settings for a project
type ProjectSettings struct {
	Project      *Project                      `json:"project"`
	Integrations []*ProjectSettingsIntegration `json:"integrations"`
}

// ProjectSettingsIntegration represents a project integration
type ProjectSettingsIntegration struct {
	Type     IntegrationType `json:"type"`
	ID       uuid.UUID       `json:"id"`
	Settings any             `json:"settings"`
}
