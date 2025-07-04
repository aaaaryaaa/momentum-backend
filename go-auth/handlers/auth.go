package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"go-auth/db"
	"go-auth/models"
	"go-auth/utils"

	"github.com/dgrijalva/jwt-go"
)

// func Signup(w http.ResponseWriter, r *http.Request) {
// 	var user models.User
// 	json.NewDecoder(r.Body).Decode(&user)

//		hashed, _ := utils.HashPassword(user.Password)
//		err := db.DB.QueryRow("INSERT INTO users(email, password) VALUES($1, $2) RETURNING id", user.Email, hashed).Scan(&user.ID)
//		if err != nil {
//			http.Error(w, "User creation failed", http.StatusInternalServerError)
//			return
//		}
//		json.NewEncoder(w).Encode(user)
//	}
func Signup(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	hashed, _ := utils.HashPassword(user.Password)
	err := db.DB.QueryRow(
		"INSERT INTO users(name, email, password) VALUES($1, $2, $3) RETURNING id",
		user.Name, user.Email, hashed,
	).Scan(&user.ID)

	if err != nil {
		http.Error(w, "User creation failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds models.User
	json.NewDecoder(r.Body).Decode(&creds)

	var user models.User
	err := db.DB.QueryRow("SELECT id, password FROM users WHERE email=$1", creds.Email).Scan(&user.ID, &user.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if !utils.CheckPasswordHash(creds.Password, user.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenStr, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	w.Write([]byte(tokenStr))
}

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID int `json:"user_id"`
	}
	json.NewDecoder(r.Body).Decode(&input)

	var user models.User
	err := db.DB.QueryRow(
		"SELECT name, email FROM users WHERE id = $1",
		input.UserID,
	).Scan(&user.Name, &user.Email)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func GetMe(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" || !strings.HasPrefix(tokenStr, "Bearer ") {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	userID, ok := claims["user_id"].(float64) // JWT stores numeric values as float64
	if !ok {
		http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
		return
	}

	var user models.User
	err = db.DB.QueryRow("SELECT name, email FROM users WHERE id = $1", int(userID)).
		Scan(&user.Name, &user.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	user.ID = int(userID)
	json.NewEncoder(w).Encode(user)
}
