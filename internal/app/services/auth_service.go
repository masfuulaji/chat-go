package services

import (
	"errors"

	"github.com/masfuulaji/go-chat/internal/app/models"
	"github.com/masfuulaji/go-chat/internal/app/repositories"
	"github.com/masfuulaji/go-chat/internal/app/request"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(name, password string) (*models.User, error)
	Register(user *request.UserRequestInsert) (*models.User, error)
}

type AuthServiceImpl struct {
	UserRepository repositories.UserRepository
}

func NewAuthService() *AuthServiceImpl {
	return &AuthServiceImpl{UserRepository: repositories.NewUserRepository()}
}

func (s *AuthServiceImpl) Login(name, password string) (*models.User, error) {
	user, err := s.UserRepository.GetUserByName(name)
    if err != nil {
        return nil, err
    }
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return nil, err
    }

    return user, nil
}

func (s *AuthServiceImpl) Register(user *request.UserRequestInsert) (*models.User, error) {
    nUser, err := s.UserRepository.CountUsers(bson.M{"name": user.Name})
    if err != nil {
        return nil, err
    }
    if nUser > 0 {
        return nil, errors.New("user already exists")
    }
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    user.Password = string(hashedPassword)

    id, err := s.UserRepository.InsertUser(user)
    if err != nil {
        return nil, err
    }

    return s.UserRepository.GetUser(id)
}
