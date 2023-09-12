package repositories

import (
	"context"
	"time"

	"github.com/masfuulaji/go-chat/internal/app/models"
	"github.com/masfuulaji/go-chat/internal/app/request"
	"github.com/masfuulaji/go-chat/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageRepository interface {
    CreateMessage(message *request.MessageRequestInsert) error 
    GetMessage(messageID string) (*models.Message, error)
    GetMessages(roomID string) ([]*models.Message, error)
    UpdateMessage(messageID string, message *models.Message) error
    DeleteMessage(messageID string) error
}

type MessageRepositoryImpl struct {
    db *mongo.Database
}

func NewMessageRepository() *MessageRepositoryImpl {
    return &MessageRepositoryImpl{db: database.DB}
}

func (r *MessageRepositoryImpl) CreateMessage(message *request.MessageRequestInsert) error {
    message.CreatedAt = time.Now()
    message.UpdatedAt = time.Now()
    _, err := r.db.Collection("messages").InsertOne(context.TODO(), message)
    return err
}

func (r *MessageRepositoryImpl) GetMessage(messageID string) (*models.Message, error) {
    var message models.Message
    objectID, err := primitive.ObjectIDFromHex(messageID)
    if err != nil {
        return &message, err
    }
    err = r.db.Collection("messages").FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&message)
    return &message, err
}

func (r *MessageRepositoryImpl) GetMessages(roomID string) ([]*models.Message, error) {
    var messages []*models.Message
    cursor, err := r.db.Collection("messages").Find(context.TODO(), bson.M{})
    if err != nil {
        return messages, err
    }

    for cursor.Next(context.TODO()) {
        var message models.Message
        err := cursor.Decode(&message)
        if err != nil {
            return messages, err
        }
    }
    return messages, err
}

func (r *MessageRepositoryImpl) UpdateMessage(messageID string, message *models.Message) error {
    message.UpdatedAt = time.Now()
    objectID, err := primitive.ObjectIDFromHex(messageID)
    if err != nil {
        return err
    }
    _, err = r.db.Collection("messages").UpdateOne(context.TODO(), bson.M{"_id": objectID}, bson.M{"$set": message})
    return err
}

func (r *MessageRepositoryImpl) DeleteMessage(messageID string) error {
    objectID, err := primitive.ObjectIDFromHex(messageID)
    if err != nil {
        return err
    }
    _, err = r.db.Collection("messages").DeleteOne(context.TODO(), bson.M{"_id": objectID})
    return err
}
