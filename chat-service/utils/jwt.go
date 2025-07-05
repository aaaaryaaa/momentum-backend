// chat-service/utils/jwt.go
package utils

import (
	"chat-service/models"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
)

func ValidateJWT(token string) (*models.User, error) {
	token = strings.TrimPrefix(token, "Bearer ")

	req, err := http.NewRequest("GET", os.Getenv("AUTH_SERVICE_URL")+"/me", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, errors.New("auth service verification failed")
	}
	defer resp.Body.Close()

	var user models.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
