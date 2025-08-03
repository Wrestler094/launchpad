package api

import (
	"context"
)

type contextKey string

const userAddressKey contextKey = "user_address"

// addUserToContext adds a user address to the request context
func addUserToContext(ctx context.Context, address string) context.Context {
	return context.WithValue(ctx, userAddressKey, address)
}

// getUserFromContext gets the user address from the request context
func getUserFromContext(ctx context.Context) string {
	if address, ok := ctx.Value(userAddressKey).(string); ok {
		return address
	}
	return ""
}