package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
	"errors"

	"dionysus/services"

	"github.com/bwmarrin/lit"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
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

type UserIDResponse struct {
	User_ID string `json:"user_id"`
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
		json.NewEncoder(w).Encode(APIResponse{403, "incorrect login"})
		lit.Debug("Password couldn't be found")
		return
	}

	lit.Debug(hash_password + " " + auth.Password)
	err = bcrypt.CompareHashAndPassword([]byte(hash_password), []byte(auth.Password))

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(APIResponse{403, "incorrect login"})
		lit.Debug("Incorrect password supplied: " + auth.Password)
		lit.Error(err.Error())
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

func TokenToID(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Token")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(APIResponse{403, "no token provided"})
		return
	}
	claims := &Claims{}
	b, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(APIResponse{403, err.Error()})
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{400, err.Error()})
		return
	}
	user_service := services.UserService{}
	id, err := user_service.UsernameToID(claims.Username)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(APIResponse{403, err.Error()})
		return
	}
	if !b.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(APIResponse{403, err.Error()})
		return
	}
	json.NewEncoder(w).Encode(APIResponse{200, UserIDResponse{id}})
}

func AdminAccess(r *http.Request) bool {
	token := r.Header.Get("Token")
	if token != "" {
		return IsAdmin(token)
	} else {
		return false
	}

	return true
}

func IsAdmin(token string) bool {
	role, err := tokenToRole(token)
	if err != nil {
		return false
	}
	return role > 0
}

func tokenToID(token string) string {
	claims := &Claims{}
	b, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return ""
	}
	user_service := services.UserService{}
	id, err := user_service.UsernameToID(claims.Username)
	if err != nil {
		return ""
	}
	if !b.Valid {
		return ""
	}
	return id
}

func generateToken(auth Authentication) (string, error) {
	expirationTime := time.Now().Add(60 * time.Minute)
	claims := &Claims{
		Username: auth.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SECRET_KEY))
}

func AuthenticationHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Token")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(APIResponse{403, "no token provided"})
			return
		}
		claims := &Claims{}
		b, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(APIResponse{403, err.Error()})
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(APIResponse{400, err.Error()})
			return
		}
		if !b.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(APIResponse{403, err.Error()})
			return
		}
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

func tokenToRole(token string) (int, error) {
	id := tokenToID(token)
	if id == "" {
		return -1, errors.New("Invalid token")
	}

	s := services.UserService{}
	user, err := s.Get(id)
	if err != nil {
		return -1, err
	}

	return user.Role, nil
}
