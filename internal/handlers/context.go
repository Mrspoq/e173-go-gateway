package handlers

import (
	"context"
	"github.com/e173-gateway/e173_go_gateway/pkg/models"
)

type contextKey string

const (
	userContextKey contextKey = "user"
)

// WithUser adds a user to the context
func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// GetUserFromContext retrieves the user from the context
func GetUserFromContext(ctx context.Context) *models.User {
	user, ok := ctx.Value(userContextKey).(*models.User)
	if !ok {
		return nil
	}
	return user
}
