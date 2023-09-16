package services

import (
	"errors"

	"github.com/masfuulaji/go-chat/internal/app/models"
	"github.com/masfuulaji/go-chat/internal/app/repositories"
	"github.com/masfuulaji/go-chat/internal/app/request"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
    CreateUser(user *request.UserRequestInsert) error
    GetUser(userID string) (*models.User, error)
    GetUsers() ([]*models.User, error)
    UpdateUser(userID string, user *models.User) error
    DeleteUser(userID string) error
}

type UserServiceImpl struct {
    userRepository repositories.UserRepository
}

func NewUserService() *UserServiceImpl {
    return &UserServiceImpl{userRepository: repositories.NewUserRepository()}
}

func (u *UserServiceImpl) CreateUser(user *request.UserRequestInsert) error {
    nUser, err := u.userRepository.CountUsers(bson.M{"name": user.Name})
    if err != nil {
        return err
    }
    if nUser > 0 {
        return errors.New("user already exists")
    }
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    user.Password = string(hashedPassword)

    _, err = u.userRepository.InsertUser(user)
    if err != nil {
        return err
    }

    return nil
}

func (u *UserServiceImpl) GetUser(userID string) (*models.User, error) {
    user, err := u.userRepository.GetUser(userID)
    if err != nil {
        return nil, err
    }

    return user, nil
}

func (u *UserServiceImpl) GetUsers() ([]*models.User, error) {
    users, err := u.userRepository.GetUsers()
    if err != nil {
        return nil, err
    }

    return users, nil
}

func (u *UserServiceImpl) UpdateUser(userID string, user *models.User) error {
    err := u.userRepository.UpdateUser(userID, user)
    if err != nil {
        return err
    }

    return nil
}

func (u *UserServiceImpl) DeleteUser(userID string) error {
    err := u.userRepository.DeleteUser(userID)
    if err != nil {
        return err
    }

    return nil
}
