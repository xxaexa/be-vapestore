package userUseCase

import (
	"clean-architecture/model/dto/userDto"
	"clean-architecture/model/entity"
	"clean-architecture/src/user"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

type UserUC struct {
	userRepo user.UserRepository
}

func NewUserUseCase(userRepo user.UserRepository) user.UserUseCase {
	return &UserUC{userRepo}
}

func (useCase *UserUC) CreateUser(user *userDto.CreateUserRequest) error {
	return useCase.userRepo.CreateUser(user)
}

func (useCase *UserUC) GetUserByEmail(email string) (*entity.User, error) {
	return useCase.userRepo.GetUserByEmail(email)
}

func (useCase *UserUC) ComparePasswords(hashed string, plain []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), plain)
	return err == nil
}

func (useCase *UserUC) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (useCase *UserUC) IsValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}

func (useCase *UserUC) GetUsers(page, limit int, email, fullName string) ([]*entity.User, int, error) {
	return useCase.userRepo.GetUsers(page, limit, email, fullName)
}

func (useCase *UserUC) GetUserByID(id string) (*entity.User, error) {
	return useCase.userRepo.GetUserByID(id)
}

func (useCase *UserUC) UpdateUser(user *userDto.UpdateUserRequest) error {
	return useCase.userRepo.UpdateUser(user)
}

func (useCase *UserUC) DeleteUser(id string) error {
	return useCase.userRepo.DeleteUser(id)
}
