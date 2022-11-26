package entities

import (
	"time"

	"github.com/google/uuid"
)

// Project is a superbutton project belonging to a user
type Project struct {
	ID                     uuid.UUID `json:"id" gorm:"primaryKey;type:string;" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	UserID                 UserID    `json:"user_id" example:"WB7DRDWrJZRGbYrv2CKGkqbzvqdC"`
	URL                    string    `json:"url" example:"https://example.com"`
	Name                   string    `json:"name" example:"Joe's Store"`
	Icon                   string    `json:"icon" example:"https://cdn.superbutton.app/chat-icon.svg"`
	Greeting               string    `json:"greeting" example:"Need some help?"`
	GreetingTimeoutSeconds uint      `json:"greeting_timeout_seconds" example:"0"`
	Color                  string    `json:"color" example:"#283593"`
	CreatedAt              time.Time `json:"created_at" example:"2022-06-05T14:26:02.302718+03:00"`
	UpdatedAt              time.Time `json:"updated_at" example:"2022-06-05T14:26:10.303278+03:00"`
}
