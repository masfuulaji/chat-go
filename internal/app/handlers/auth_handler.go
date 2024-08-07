package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/securecookie"
	"github.com/masfuulaji/go-chat/internal/app/request"
	"github.com/masfuulaji/go-chat/internal/app/services"
)

type AuthHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	CheckAuth(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

type AuthHandlerImpl struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandlerImpl {
	return &AuthHandlerImpl{authService: authService}
}

func (ah *AuthHandlerImpl) Login(w http.ResponseWriter, r *http.Request) {
	var data request.UserRequestInsert

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := ah.authService.Login(data.Name, data.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   user.ID,
		"name": user.Name,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encoded, err := securecookie.New([]byte("secret"), nil).Encode("jwt", tokenString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    encoded,
		Path:     "/",
		HttpOnly: true,
		// Expires: time.Now().Add(time.Hour * 24),
		Secure: false,
	}

	http.SetCookie(w, cookie)

	respose := map[string]string{
		"message": "User logged in",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respose)
}

func (ah *AuthHandlerImpl) Register(w http.ResponseWriter, r *http.Request) {
	var data request.UserRequestInsert

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = ah.authService.Register(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respose := map[string]string{
		"message": "User registered",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respose)
}

func (ah *AuthHandlerImpl) CheckAuth(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var tokenString string

	err = securecookie.New([]byte("secret"), nil).Decode("jwt", cookie.Value, &tokenString)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("secret"), nil
	})
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	claim, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if time.Now().Unix() > int64(claim["exp"].(float64)) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	response := map[string]interface{}{
		"message": "Authorized",
		"status":  true,
		"id":      claim["id"].(string),
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (ah *AuthHandlerImpl) Logout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		Expires:  time.Now(),
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)

	respose := map[string]string{
		"message": "User logged out",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respose)
}
