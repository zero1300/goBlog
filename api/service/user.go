package service

import (
	"blog/api/repository"
	"blog/models"
)

//UserService UserService struct
type UserService struct {
	repo repository.UserRepository
}

//NewUserService : get injected user repo
func NewUserService(repo repository.UserRepository) UserService {
	return UserService{
		repo: repo,
	}
}

// FindAllUser -> get all user info
func (u UserService) FindAllUser() (*[]models.User, int64, error) {
	return u.repo.FindAll()
}

//Save -> saves users entity
func (u UserService) CreateUser(user models.UserRegister) error {
	return u.repo.CreateUser(user)
}

//Login -> Gets validated user
func (u UserService) LoginUser(user models.UserLogin) (*models.User, error) {
	return u.repo.LoginUser(user)
}

//Delete -> delete user enttity
func (u UserService) DeleteUser(id int64) error {
	var user models.User
	user.ID = id
	return u.repo.DeleteUser(user)
}
