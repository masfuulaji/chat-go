package handlers

import (
	"encoding/json"
	"net/http"

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

    userJSON, err := json.Marshal(user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write(userJSON)
}

func (ah *AuthHandlerImpl) Register(w http.ResponseWriter, r *http.Request) {
    var data request.UserRequestInsert

    err := json.NewDecoder(r.Body).Decode(&data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    user, err := ah.authService.Register(&data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    userJSON, err := json.Marshal(user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write(userJSON)
}
