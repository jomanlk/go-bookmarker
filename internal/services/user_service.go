package services

import (
	"bookmarker/internal/models"
	"bookmarker/internal/repositories"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(username, password string) (models.User, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(username, password string) (models.User, error) {
	if username == "" || password == "" {
		return models.User{}, errors.New("username and password required")
	}
	_, err := s.repo.GetUserByUsername(username)
	if err == nil {
		return models.User{}, errors.New("username already exists")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}
	return s.repo.CreateUser(username, string(hash))
}
