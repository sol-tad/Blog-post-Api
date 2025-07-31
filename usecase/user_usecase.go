package usecase

import (
	"context"
	"errors"

	"github.com/sol-tad/Blog-post-Api/domian"
	"github.com/sol-tad/Blog-post-Api/infrastructure"
)

type UserUsecase struct {
	UserRepository domian.UserRepository
}

func NewUserUsecase(userRepo domian.UserRepository) *UserUsecase{
	return &UserUsecase{
		UserRepository: userRepo,
	}
}
func (uuc *UserUsecase) Register(ctx context.Context,user domian.User) error{
	if user.Username=="" ||user.Password==""{
		return errors.New("missing required fileds")
	}

	hashedPassword,_:=infrastructure.HashPassword(user.Password)
	user.Password=hashedPassword

	_, err := uuc.UserRepository.Register(ctx, user)
	return err
}


func (uuc *UserUsecase) Login(ctx context.Context, username, password string )(accessToken string, refreshToken string, err error){
	user,err:=uuc.UserRepository.Login(ctx,username)
	 if err != nil {
        return "", "", errors.New("invalid username or password")
    }

	if!infrastructure.CheckPassword(password,user.Password){
		return "", "", errors.New("invalid username or password")
	}

	accessToken,err=infrastructure.GenerateAccessToken(user.ID.Hex(),user.Role)
	    if err != nil {
        return "", "", err
    }
	refreshToken,err=infrastructure.GenerateRefreshToken(user.ID.Hex())
	    if err != nil {
        return "", "", err
    }

	//store the refresh token in database
	err=uuc.UserRepository.SaveRefreshToken(ctx, user.ID.Hex(), refreshToken)
	if err != nil {
        return "", "", err
    }

    return accessToken, refreshToken, nil

}


func (uuc *UserUsecase) RefreshToken(ctx context.Context, refreshToken string)(string, error){

    userID, err := infrastructure.VerifyRefreshToken(refreshToken)
    if err != nil {
        return "", err
    }

	ok, err := uuc.UserRepository.VerifyRefreshToken(ctx, userID, refreshToken)
    if err != nil || !ok {
        return "", errors.New("invalid or expired refresh token")
    }

	    // Generate new access token
    user, err := uuc.UserRepository.FindByID(ctx, userID)
    if err != nil {
        return "", err
    }

    return infrastructure.GenerateAccessToken(user.ID.Hex(), user.Role)

}

func (uuc *UserUsecase) Logout(ctx context.Context, userID string) error {
    return uuc.UserRepository.DeleteRefreshToken(ctx, userID)
}