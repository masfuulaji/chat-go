package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/masfuulaji/go-chat/internal/app/models"
	"github.com/masfuulaji/go-chat/internal/app/request"
	"github.com/masfuulaji/go-chat/internal/app/services"
)

type UserHandler interface {
    CreateUser(w http.ResponseWriter, r *http.Request)
    GetUser(w http.ResponseWriter, r *http.Request)
    GetUsers(w http.ResponseWriter, r *http.Request)
    UpdateUser(w http.ResponseWriter, r *http.Request)
    DeleteUser(w http.ResponseWriter, r *http.Request)
}

type UserHandlerImpl struct {
    userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandlerImpl {
    return &UserHandlerImpl{userService: userService}
}

func (uh *UserHandlerImpl) CreateUser(w http.ResponseWriter, r *http.Request) {
    var user request.UserRequestInsert
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    err = uh.userService.CreateUser(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("{\"message\": \"User created\"}"))
}

func (uh *UserHandlerImpl) GetUser(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    user, err := uh.userService.GetUser(params["id"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(user)
}

func (uh *UserHandlerImpl) GetUsers(w http.ResponseWriter, r *http.Request) {
    users, err := uh.userService.GetUsers()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(users)
}

func (uh *UserHandlerImpl) UpdateUser(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)

    var user models.User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    err = uh.userService.UpdateUser(params["id"],&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    response := map[string]string{"message": "User updated"}

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}

func (uh *UserHandlerImpl) DeleteUser(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    err := uh.userService.DeleteUser(params["id"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    response := map[string]string{"message": "User deleted"}

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}
