package entities

import (
	"time"
)

// UserID is the ID of a user
type UserID string

// User stores information about a user
type User struct {
	ID        UserID    `json:"id" gorm:"primaryKey;type:string;" example:"WB7DRDWrJZRGbYrv2CKGkqbzvqdC"`
	Email     string    `json:"email" example:"name@email.com"` // gorm:"uniqueIndex"
	Name      string    `json:"name" example:"John Doe"`
	CreatedAt time.Time `json:"created_at" example:"2022-06-05T14:26:02.302718+03:00"`
	UpdatedAt time.Time `json:"updated_at" example:"2022-06-05T14:26:10.303278+03:00"`
}
