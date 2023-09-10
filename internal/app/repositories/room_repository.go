package repositories

import (
	"context"
	"time"

	"github.com/masfuulaji/go-chat/internal/app/models"
	"github.com/masfuulaji/go-chat/internal/app/request"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomRepository interface {
    CreateRoom(room *request.RoomRequestInsert) error
    GetRoom(roomID string) (*models.Room, error)
    GetRooms() ([]*models.Room, error)
    UpdateRoom(roomID string,room *models.Room) error
    DeleteRoom(roomID string) error
}

type RoomRepositoryImpl struct {
    db *mongo.Database
}

func NewRoomRepository(db *mongo.Database) *RoomRepositoryImpl {
    return &RoomRepositoryImpl{db: db}
}

func (r *RoomRepositoryImpl) CreateRoom(room *request.RoomRequestInsert) error {
    room.CreatedAt = time.Now()
    room.UpdatedAt = time.Now()
    _, err := r.db.Collection("rooms").InsertOne(context.TODO(), room)
    return err
}

func (r *RoomRepositoryImpl) GetRoom(roomID string) (*models.Room, error) {
    var room models.Room
    
    objectID, err := primitive.ObjectIDFromHex(roomID)
    if err != nil {
        return &room, err
    }
    err = r.db.Collection("rooms").FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&room)
    return &room, err
}

func (r *RoomRepositoryImpl) GetRooms() ([]*models.Room, error) {
    var rooms []*models.Room
    cursor, err := r.db.Collection("rooms").Find(context.TODO(), bson.M{})
    if err != nil {
        return rooms, err
    }

    for cursor.Next(context.TODO()) {
        var room models.Room
        err := cursor.Decode(&room)
        if err != nil {
            return rooms, err
        }
        rooms = append(rooms, &room)
    }
    return rooms, err
}

func (r *RoomRepositoryImpl) UpdateRoom(roomID string,room *models.Room) error {
    objectID, err := primitive.ObjectIDFromHex(roomID)
    if err != nil {
        return err
    }
    _, err = r.db.Collection("rooms").UpdateOne(context.TODO(), bson.M{"_id": objectID}, bson.M{"$set": room})
    return err
}

func (r *RoomRepositoryImpl) DeleteRoom(roomID string) error {
    objectID, err := primitive.ObjectIDFromHex(roomID)
    if err != nil {
        return err
    }
    _, err = r.db.Collection("rooms").UpdateOne(context.TODO(), bson.M{"_id": objectID}, bson.M{"$set": bson.M{"deleted_at": time.Now()}})
    return err
}
