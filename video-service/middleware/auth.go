package middleware

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

type UserData struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func Protected(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" || !strings.HasPrefix(token, "Bearer ") {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")

		// Call auth service to get user data
		req, _ := http.NewRequest("GET", os.Getenv("AUTH_SERVICE_URL")+"/me", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != 200 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var user UserData
		json.NewDecoder(resp.Body).Decode(&user)
		defer resp.Body.Close()

		// Set user in context
		ctx := r.Context()
		ctx = contextWithUser(ctx, &user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
