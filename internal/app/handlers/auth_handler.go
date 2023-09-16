package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/masfuulaji/go-chat/internal/app/request"
	"github.com/masfuulaji/go-chat/internal/app/services"
)

type AuthHandler interface {
    Login(w http.ResponseWriter, r *http.Request)
    Register(w http.ResponseWriter, r *http.Request)
}

type AuthHandlerImpl struct {
    authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandlerImpl {
    return &AuthHandlerImpl{authService: authService}
}

func (ah *AuthHandlerImpl) Login(w http.ResponseWriter, r *http.Request) {
    user, err := ah.authService.Login(r.FormValue("name"), r.FormValue("password"))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "id": user.ID,
        "name": user.Name,
        "exp": time.Now().Add(time.Hour * 24).Unix(),
    })

    tokenString, err := token.SignedString([]byte("secret"))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    response := map[string]string{
        "token": tokenString,
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
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
