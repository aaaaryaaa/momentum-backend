package middleware

import (
	"context"
	"net/http"
	"post-service/models"
	"strings"

	"github.com/golang-jwt/jwt"
)

func Protected(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("funnilol"), nil
		})
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userIDFloat := claims["user_id"].(float64)
		ctx := context.WithValue(r.Context(), "user", &models.User{ID: int(userIDFloat)})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(ctx context.Context) *models.User {
	user, _ := ctx.Value("user").(*models.User)
	return user
}
