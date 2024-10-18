package service

import (
	"github.com/google/uuid"
	"kontest-user-service/model"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func getUser(userID uuid.UUID) (*model.User, error) {
	return nil, nil
}
