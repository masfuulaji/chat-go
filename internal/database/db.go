package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func InitMongoDB() (*mongo.Database, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var credentials = options.Credential{
        Username: "root",
        Password: "root",
    }

    connectionOprtions := options.Client().ApplyURI("mongodb://mongodb:27017").SetAuth(credentials)
    client, err := mongo.Connect(ctx, connectionOprtions)
    if err != nil {
        return nil, err
    }

    DB = client.Database("chatter")

    return client.Database("chatter"), nil
}
