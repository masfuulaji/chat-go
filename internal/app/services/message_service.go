package services

import (
	"github.com/masfuulaji/go-chat/internal/app/models"
	"github.com/masfuulaji/go-chat/internal/app/repositories"
	"github.com/masfuulaji/go-chat/internal/app/request"
)

type MessageService interface {
    CreateMessage(message *request.MessageRequestInsert, userId string) error
    GetMessage(messageID string) (*models.Message, error)
    GetMessages(roomID string) ([]*models.Message, error)
    UpdateMessage(messageID string, message *models.Message) error
    DeleteMessage(messageID string) error
}

type MessageServiceImpl struct {
    MessageRepository repositories.MessageRepository
}

func NewMessageService() *MessageServiceImpl {
    return &MessageServiceImpl{MessageRepository: repositories.NewMessageRepository()}
}

func (s *MessageServiceImpl) CreateMessage(message *request.MessageRequestInsert, userId string) error {
    return s.MessageRepository.CreateMessage(message, userId)
}

func (s *MessageServiceImpl) GetMessage(messageID string) (*models.Message, error) {
    return s.MessageRepository.GetMessage(messageID)
}

func (s *MessageServiceImpl) GetMessages(roomID string) ([]*models.Message, error) {
    return s.MessageRepository.GetMessages(roomID)
}

func (s *MessageServiceImpl) UpdateMessage(messageID string, message *models.Message) error {
    return s.MessageRepository.UpdateMessage(messageID, message)
}

func (s *MessageServiceImpl) DeleteMessage(messageID string) error {
    return s.MessageRepository.DeleteMessage(messageID)
}
