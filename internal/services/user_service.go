package services

import (
	"bookmarker/internal/models"
	"bookmarker/internal/repositories"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{UserRepo: userRepo}
}

func (s *UserService) GetUserByUsername(username string) (models.User, error) {
	return s.UserRepo.GetUserByUsername(username)
}

func (s *UserService) GetUserByID(id int64) (models.User, error) {
	return s.UserRepo.GetUserByID(id)
}

func (s *UserService) CreateUser(username, password string) (models.User, error) {
	if username == "" || password == "" {
		return models.User{}, errors.New("username and password required")
	}
	_, err := s.UserRepo.GetUserByUsername(username)
	if err == nil {
		return models.User{}, errors.New("username already exists")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}
	return s.UserRepo.CreateUser(username, string(hash))
}
