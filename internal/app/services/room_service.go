package services

import (
	"github.com/masfuulaji/go-chat/internal/app/models"
	"github.com/masfuulaji/go-chat/internal/app/repositories"
	"github.com/masfuulaji/go-chat/internal/app/request"
)

type RoomService interface {
	CreateRoom(room *request.RoomRequestInsert) error
	GetRoom(roomID string) (*models.Room, error)
	GetRooms() ([]*models.Room, error)
	UpdateRoom(roomID string, room *models.Room) error
	DeleteRoom(roomID string) error
}

type RoomServiceImpl struct {
	roomRepository *repositories.RoomRepositoryImpl
}

func NewRoomService() *RoomServiceImpl {
	return &RoomServiceImpl{roomRepository: repositories.NewRoomRepository()}
}

func (r *RoomServiceImpl) CreateRoom(room *request.RoomRequestInsert) error {
	return r.roomRepository.CreateRoom(room)
}

func (r *RoomServiceImpl) GetRoom(roomID string) (*models.Room, error) {
	return r.roomRepository.GetRoom(roomID)
}

func (r *RoomServiceImpl) GetRooms() ([]*models.Room, error) {
	return r.roomRepository.GetRooms()
}

func (r *RoomServiceImpl) UpdateRoom(roomID string, room *models.Room) error {
	return r.roomRepository.UpdateRoom(roomID, room)
}

func (r *RoomServiceImpl) DeleteRoom(roomID string) error {
	return r.roomRepository.DeleteRoom(roomID)
}
