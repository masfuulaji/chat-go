package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
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
    var message request.MessageRequestInsert
    err := json.NewDecoder(r.Body).Decode(&message)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    err = mh.messageService.CreateMessage(&message)
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
