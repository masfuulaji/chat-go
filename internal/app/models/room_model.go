package models

import "time"

type Room struct {
    ID   string `json:"id" bson:"_id"`
    Name string `json:"name"`
    CreatedAt time.Time `json:"created_at" bson:"created_at"`
    UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
    DeletedAt time.Time `json:"deleted_at" bson:"deleted_at"`
}
