package user

import (
	"clean-architecture/model/dto/userDto"
	"clean-architecture/model/entity"
)

type UserRepository interface {
	CreateUser(user *userDto.CreateUserRequest) error
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByID(id string) (*entity.User, error)
	GetUsers(page, limit int, email, fullName string) ([]*entity.User, int, error)
	UpdateUser(user *userDto.UpdateUserRequest) error
	DeleteUser(id string) error
}

type UserUseCase interface {
	CreateUser(user *userDto.CreateUserRequest) error
	GetUserByEmail(email string) (*entity.User, error)
	GetUsers(page, limit int, email, fullName string) ([]*entity.User, int, error)
	GetUserByID(id string) (*entity.User, error)
	UpdateUser(user *userDto.UpdateUserRequest) error
	DeleteUser(id string) error
	ComparePasswords(hashed string, plain []byte) bool
	HashPassword(password string) (string, error)
	IsValidPassword(password string) bool
}
