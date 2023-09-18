package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/masfuulaji/go-chat/internal/app/models"
	"github.com/masfuulaji/go-chat/internal/app/request"
	"github.com/masfuulaji/go-chat/internal/app/services"
)

type MessageHandler interface {
    CreateMessage(w http.ResponseWriter, r *http.Request)
    GetMessages(w http.ResponseWriter, r *http.Request)
    GetMessage(w http.ResponseWriter, r *http.Request)
    DeleteMessage(w http.ResponseWriter, r *http.Request)
    UpdateMessage(w http.ResponseWriter, r *http.Request)
}

type MessageHandlerImpl struct {
    messageService services.MessageService
}

func NewMessageHandler(messageService services.MessageService) *MessageHandlerImpl {
    return &MessageHandlerImpl{messageService: messageService}
}

func (mh *MessageHandlerImpl) CreateMessage(w http.ResponseWriter, r *http.Request) {
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

    var message request.MessageRequestInsert
    err = json.NewDecoder(r.Body).Decode(&message)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    userID := claim["id"].(string)

    err = mh.messageService.CreateMessage(&message, userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("{\"message\": \"Message created\"}"))
}

func (mh *MessageHandlerImpl) GetMessages(w http.ResponseWriter, r *http.Request) {
    room := r.URL.Query().Get("room")
    messages, err := mh.messageService.GetMessages(room)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(messages)
}

func (mh *MessageHandlerImpl) GetMessage(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    message, err := mh.messageService.GetMessage(params["id"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(message)
}

func (mh *MessageHandlerImpl) UpdateMessage(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var message models.Message
    err := json.NewDecoder(r.Body).Decode(&message)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    err = mh.messageService.UpdateMessage(params["id"], &message)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func (mh *MessageHandlerImpl) DeleteMessage(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    err := mh.messageService.DeleteMessage(params["id"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
    w.Write([]byte("{\"message\": \"Message deleted\"}"))
}
