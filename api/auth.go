package api

import (
	"time"
	"encoding/json"
	"net/http"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
	"dionysus/services"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/bwmarrin/lit"
)

var SECRET_KEY string

type Authentication struct {
	Username string
	Password string
}

type Claims struct {
	Username string
	jwt.StandardClaims
}

type TokenResponse struct {
	Token string `json:"token"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var auth Authentication
	err := json.NewDecoder(r.Body).Decode(&auth)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		return
	}

	s := services.AuthService{}
	hash_password, err := s.GetPasswordHash(auth.Username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(APIResponse{403, "Incorrect login"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash_password), []byte(auth.Password))
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(APIResponse{403, "Incorrect login"})
		return
	}

	tokenString, err := generateToken(auth)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{500, err.Error()})
		return
	}
	json.NewEncoder(w).Encode(APIResponse{200, TokenResponse{tokenString}})
	return
}

func generateToken(auth Authentication) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: auth.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SECRET_KEY))
}

func authenticationHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(fmt.Sprint(r.Body))
		next.ServeHTTP(w, r)
	})
}

func passwordToHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hash), err
}

func SetSecretKey() {
	SECRET_KEY = os.Getenv("DIONYSUS_KEY")
	if SECRET_KEY == "" {
		lit.Warn("You are running dionysus with no secret key.")
	}
}