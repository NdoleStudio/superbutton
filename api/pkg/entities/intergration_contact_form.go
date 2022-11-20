package entities

import (
	"time"

	"github.com/google/uuid"
)

// IntegrationContactForm are contact form integration settings
type IntegrationContactForm struct {
	ID              uuid.UUID `json:"id" gorm:"primaryKey;type:string;" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	UserID          UserID    `json:"user_id" example:"WB7DRDWrJZRGbYrv2CKGkqbzvqdC"`
	ProjectID       uuid.UUID `json:"project_id" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	NameText        string    `json:"name_text" example:"Name"`
	NameEnabled     bool      `json:"name_enabled" example:"true"`
	EmailText       string    `json:"email_text" example:"Email"`
	EmailEnabled    string    `json:"email_enabled" example:"true"`
	PhoneNumberText string    `json:"phone_number_text" example:"Phone Number"`
	PhoneEnabled    string    `json:"phone_enabled" example:"true"`
	MessageText     string    `json:"message_text" example:"Message"`
	MessageEnabled  string    `json:"message_enabled" example:"true"`
	Icon            string    `json:"icon" example:"https://cdn.superbutton.app/contact-form-icon.svg"`
	CreatedAt       time.Time `json:"created_at" example:"2022-06-05T14:26:02.302718+03:00"`
	UpdatedAt       time.Time `json:"updated_at" example:"2022-06-05T14:26:10.303278+03:00"`
}
