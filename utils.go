package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const (
	ContextKeyUser ContextKey = "user"
)

type ContextKey string

type User struct {
	Id             int    `json:"id"`
	Username       string `json:"username"`
	PasswordHashed string `json:"passwordHashed"`
	Level          int    `json:"level"`
}

func (u *User) VerifyPassword(pass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHashed), []byte(pass)) == nil
}

type Product struct {
	Id          int     `json:"id"`
	Type        string  `json:"type"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Gender      string  `json:"gender"`
	Color       string  `json:"color"`
	Small       int     `json:"small"`
	Medium      int     `json:"medium"`
	Large       int     `json:"large"`
	ImageUrl    string  `json:"imageUrl"`
	ImageAlt    string  `json:"imageAlt"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func WriteJSON(w http.ResponseWriter, v any, status int) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error(err.Error())
		http.Error(w, "error: internal server error", http.StatusInternalServerError)
	}
}

func CreateJWT(id int) (string, error) {
	claims := &jwt.MapClaims{
		"userId": id,
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
}

func CreateCookie(name string, value string) *http.Cookie {
	cookie := &http.Cookie{}

	cookie.Name = name
	cookie.Value = value
	cookie.Domain = "localhost"
	cookie.Path = "/"
	cookie.SameSite = http.SameSiteStrictMode
	// cookie.Expires = time.Now().Add(24 * time.Hour)
	// cookie.Secure = false
	// cookie.HttpOnly = false

	return cookie
}
