package service

import (
	"github.com/google/uuid"
	"kontest-user-service/database"
	error2 "kontest-user-service/error"
	"kontest-user-service/model"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (us *UserService) GetUser(uid uuid.UUID) (*model.User, error) {
	user, err := database.FindUserByID(uid)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, &error2.UserNotFoundError{}
	}

	return user, nil
}
