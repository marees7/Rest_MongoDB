package auth

import (
	"api/models"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	// "go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
	// "go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
)

func CreateToken(user_id uuid.UUID, email string) (models.LoginResponse, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user_id
	claims["email"] = email
	tokenExpiry, _ := strconv.Atoi("3600000")
	expiry := time.Now().Add(time.Millisecond * time.Duration(tokenExpiry)).Unix() //Token expires after defined time interval
	claims["exp"] = expiry
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	var accessToken string
	var err error
	accessToken, err = token.SignedString([]byte("secretkey"))
	if err != nil {
		return models.LoginResponse{}, err
	}

	loginResponse := models.LoginResponse{}
	loginResponse.AccessToken = accessToken
	loginResponse.ExpiryDate = expiry
	return loginResponse, nil
}

func ExtractTokenID(r *http.Request) (string, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secretkey"), nil
	})
	fmt.Println("token ", token)
	if err != nil {
		return "", errors.New("User not found")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		email := fmt.Sprint(claims["email"])
		fmt.Println("email ", email)
		return email, nil
	}
	return "", nil
}

func ExtractToken(r *http.Request) string {
	keys := r.URL.Query()
	token := keys.Get("token")
	if token != "" {
		return token
	}
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}
