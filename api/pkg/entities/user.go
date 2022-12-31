package entities

import (
	"time"
)

// UserID is the ID of a user
type UserID string

// SubscriptionName is the name of the subscription
type SubscriptionName string

// SubscriptionNameFree represents a free subscription
const SubscriptionNameFree = SubscriptionName("free")

// SubscriptionNameProMonthly represents a monthly pro subscription
const SubscriptionNameProMonthly = SubscriptionName("pro-monthly")

// SubscriptionNameProYearly represents a yearly pro subscription
const SubscriptionNameProYearly = SubscriptionName("pro-yearly")

// User stores information about a user
type User struct {
	ID                   UserID           `json:"id" gorm:"primaryKey;type:string;" example:"WB7DRDWrJZRGbYrv2CKGkqbzvqdC"`
	Email                string           `json:"email" example:"name@email.com"`
	Name                 string           `json:"name" example:"John Doe"`
	SubscriptionName     SubscriptionName `json:"subscription_name" example:"free"`
	SubscriptionID       string           `json:"subscription_id" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	SubscriptionStatus   string           `json:"subscription_status" example:"on_trial"`
	SubscriptionRenewsAt time.Time        `json:"subscription_renews_at" example:"2022-06-05T14:26:02.302718+03:00"`
	CreatedAt            time.Time        `json:"created_at" example:"2022-06-05T14:26:02.302718+03:00"`
	UpdatedAt            time.Time        `json:"updated_at" example:"2022-06-05T14:26:10.303278+03:00"`
}
