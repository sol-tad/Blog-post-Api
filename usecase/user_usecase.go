package usecase

import (
	"errors"

	"github.com/sol-tad/Blog-post-Api/domian"
)

type UserUsecase struct {
	UserRepository domian.UserRepository
}

func NewUserUsecase(userRepo domian.UserRepository) *UserUsecase{
	return &UserUsecase{
		UserRepository: userRepo,
	}
}
func (uuc *UserUsecase) Register(user domian.User) error{
	if user.Username=="" ||user.Password==""{
		return errors.New("missing required fileds")
	}
	_,err:=uuc.UserRepository.CreateUser(user)
	return err
}