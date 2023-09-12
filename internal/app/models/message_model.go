package models

import "time"

type Message struct {
    ID   string `json:"id" bson:"_id"`
    RoomID string `json:"room_id" bson:"room_id"`
    SenderID string `json:"sender_id" bson:"sender_id"`
    Content string `json:"content" bson:"content"`
    MessageType int `json:"message_type" bson:"message_type"`
    CreatedAt time.Time `json:"created_at" bson:"created_at"`
    UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
    DeletedAt time.Time `json:"deleted_at" bson:"deleted_at"`
}
