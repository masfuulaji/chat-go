package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/masfuulaji/go-chat/internal/app/models"
	"github.com/masfuulaji/go-chat/internal/app/request"
	"github.com/masfuulaji/go-chat/internal/app/services"
)

type RoomHandler interface {
    CreateRoom(w http.ResponseWriter, r *http.Request)
    GetRoom(w http.ResponseWriter, r *http.Request)
    GetRooms(w http.ResponseWriter, r *http.Request)
    UpdateRoom(w http.ResponseWriter, r *http.Request)
    DeleteRoom(w http.ResponseWriter, r *http.Request)
}

type RoomHandlerImpl struct {
    roomService services.RoomService
}

func NewRoomHandler(roomService services.RoomService) *RoomHandlerImpl {
    return &RoomHandlerImpl{roomService: roomService}
}

func (rm *RoomHandlerImpl) CreateRoom(w http.ResponseWriter, r *http.Request) {
    var room request.RoomRequestInsert
    err := json.NewDecoder(r.Body).Decode(&room)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    err = rm.roomService.CreateRoom(&room)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("{\"message\": \"Room created\"}"))
}

func (rm *RoomHandlerImpl) GetRoom(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    room, err := rm.roomService.GetRoom(params["id"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(room)
}

func (rm *RoomHandlerImpl) GetRooms(w http.ResponseWriter, r *http.Request) {
    rooms, err := rm.roomService.GetRooms()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(rooms)
}

func (rm *RoomHandlerImpl) UpdateRoom(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var room models.Room
    err := json.NewDecoder(r.Body).Decode(&room)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    err = rm.roomService.UpdateRoom(params["id"], &room)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("{\"message\": \"Room updated\", \"id\": \"" + params["id"] + "\"}"))
}

func (rm *RoomHandlerImpl) DeleteRoom(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    err := rm.roomService.DeleteRoom(params["id"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("{\"message\": \"Room deleted\"}"))
}
