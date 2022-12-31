package repositories

import (
	"context"

	"github.com/NdoleStudio/superbutton/pkg/entities"
)

// UserRepository loads and persists an entities.User
type UserRepository interface {
	// Store a new entities.User
	Store(ctx context.Context, user *entities.User) error

	// Update a new entities.User
	Update(ctx context.Context, user *entities.User) error

	// Load an entities.User by entities.UserID
	Load(ctx context.Context, userID entities.UserID) (*entities.User, error)

	// LoadBySubscriptionID fetches a user based on the subscriptionID
	LoadBySubscriptionID(ctx context.Context, subscriptionID string) (*entities.User, error)

	// LoadOrStore an entities.User by entities.AuthUser
	LoadOrStore(ctx context.Context, user entities.AuthUser) (*entities.User, bool, error)
}
