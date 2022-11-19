package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Project struct {
	ID        uuid.UUID      `json:"id" gorm:"primaryKey;type:string;" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	UserID    UserID         `json:"user_id" example:"WB7DRDWrJZRGbYrv2CKGkqbzvqdC"`
	URL       string         `json:"url" example:"https://example.com"`
	Name      string         `json:"name" example:"Personal Website"`
	Settings  datatypes.JSON `json:"settings" example:"{}"`
	CreatedAt time.Time      `json:"created_at" example:"2022-06-05T14:26:02.302718+03:00"`
	UpdatedAt time.Time      `json:"updated_at" example:"2022-06-05T14:26:10.303278+03:00"`
}