package entities

import (
	"time"

	"github.com/google/uuid"
)

// IntegrationWhatsapp contains whatsapp integration settings
type IntegrationWhatsapp struct {
	ID          uuid.UUID `json:"id" gorm:"primaryKey;type:string;" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	UserID      UserID    `json:"user_id" example:"WB7DRDWrJZRGbYrv2CKGkqbzvqdC"`
	ProjectID   uuid.UUID `json:"project_id" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	Enabled     bool      `json:"enabled" example:"true"`
	Name        string    `json:"name" example:"FAQ"`
	Text        string    `json:"text" example:"Contact us on WhatsApp"`
	PhoneNumber string    `json:"phone_number" example:"+18005550199"`
	Icon        string    `json:"icon" example:"https://cdn.superbutton.app/whatsapp-icon.svg"`
	CreatedAt   time.Time `json:"created_at" example:"2022-06-05T14:26:02.302718+03:00"`
	UpdatedAt   time.Time `json:"updated_at" example:"2022-06-05T14:26:10.303278+03:00"`
}
