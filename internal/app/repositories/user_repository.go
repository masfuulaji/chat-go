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

type UserRepository interface {
    InsertUser(user *request.UserRequestInsert) (string, error)
    GetUser(userID string) (*models.User, error)
    GetUsers() ([]*models.User, error)
    UpdateUser(userID string, user *models.User) error
    DeleteUser(userID string) error
    GetUserByName(name string) (*models.User, error)
    CountUsers(query bson.M) (int, error)
}

type UserRepositoryImpl struct {
    db *mongo.Database
}

func NewUserRepository() *UserRepositoryImpl {
    return &UserRepositoryImpl{db: database.DB}
}

func (r *UserRepositoryImpl) InsertUser(user *request.UserRequestInsert) (string, error) {
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()
    res, err := r.db.Collection("users").InsertOne(context.TODO(), user)
    if err != nil {
        return "", err
    }
    insertID := res.InsertedID.(primitive.ObjectID).Hex() 
    return insertID, nil
}

func (r *UserRepositoryImpl) GetUser(userID string) (*models.User, error) {
    var user models.User

    objectID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return &user, err
    }

    err = r.db.Collection("users").FindOne(context.TODO(), bson.M{"_id": objectID, "deleted_at": bson.M{"$exists": false}}).Decode(&user)
    return &user, err
}

func (r *UserRepositoryImpl) GetUsers() ([]*models.User, error) {
    var users []*models.User
    cursor, err := r.db.Collection("users").Find(context.TODO(), bson.M{"deleted_at": bson.M{"$exists": false}})
    if err != nil {
        return users, err
    }

    for cursor.Next(context.TODO()) {
        var user models.User
        err := cursor.Decode(&user)
        if err != nil {
            return users, err
        }
    }

    return users, err
}

func (r *UserRepositoryImpl) UpdateUser(userID string, user *models.User) error {
    objectID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return err
    }

    _, err = r.db.Collection("users").UpdateOne(context.TODO(), bson.M{"_id": objectID}, bson.M{"$set": user})
    return err
}

func (r *UserRepositoryImpl) DeleteUser(userID string) error {
    objectID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return err
    }
    _, err = r.db.Collection("users").UpdateOne(context.TODO(), bson.M{"_id": objectID}, bson.M{"$set": bson.M{"deleted_at": time.Now()}})
    return err
}

func (r *UserRepositoryImpl) GetUserByName(name string) (*models.User, error) {
    var user models.User
    err := r.db.Collection("users").FindOne(context.TODO(), bson.M{"name": name, "deleted_at": bson.M{"$exists": false}}).Decode(&user)
    return &user, err
}

func (r *UserRepositoryImpl) CountUsers(query bson.M) (int, error) {
    count, err := r.db.Collection("users").CountDocuments(context.TODO(), query)
    return int(count), err
}
