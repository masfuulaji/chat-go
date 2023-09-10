package request

import "time"

type RoomRequestInsert struct {
    Name string `json:"name"`
    CreatedAt time.Time `json:"created_at" bson:"created_at"`
    UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
