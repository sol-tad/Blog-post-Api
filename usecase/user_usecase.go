package usecase

import (
	"errors"

	"github.com/sol-tad/Blog-post-Api/domain"
)

type UserUsecase struct {
	UserRepository domain.UserRepository
}

func NewUserUsecase(userRepo domain.UserRepository) *UserUsecase{
	return &UserUsecase{
		UserRepository: userRepo,
	}
}
func (uuc *UserUsecase) Register(user domain.User) error{
	if user.Username=="" ||user.Password==""{
		return errors.New("missing required fileds")
	}
	_,err:=uuc.UserRepository.CreateUser(user)
	return err
}