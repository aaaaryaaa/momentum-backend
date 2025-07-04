// context_key.go
package middleware

import "context"

type contextKey string

const userKey contextKey = "user"

func contextWithUser(ctx context.Context, user *UserData) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func GetUserFromContext(ctx context.Context) *UserData {
	user, _ := ctx.Value(userKey).(*UserData)
	return user
}
