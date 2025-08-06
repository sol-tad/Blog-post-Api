package usecase

import (
	"context"
	"encoding/json"
	"io"

	"github.com/sol-tad/Blog-post-Api/config"
	"github.com/sol-tad/Blog-post-Api/domain"
	"github.com/sol-tad/Blog-post-Api/infrastructure"
)

type OAuthUsecase struct {
	userRepo domain.UserRepository
}

	
func NewOAuthUsecase(urepo domain.UserRepository) *OAuthUsecase{
	return &OAuthUsecase{userRepo: urepo}
}


//callback function
func (u *OAuthUsecase) HandleGoogleCallback(code string) (*domain.User, string, string, error){
	token,err:=config.GoogleOAuthConfig.Exchange(context.Background(),code)

	if err != nil {
		return nil, "","", err
	}

	client :=config.GoogleOAuthConfig.Client(context.Background(),token)
	response,err:=client.Get("https://www.googleapis.com/oauth2/v2/userinfo")

	if err != nil {
		return nil, "","", err
	}
	defer response.Body.Close()
	body,_:=io.ReadAll(response.Body)

	var googleUser struct{
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.Unmarshal(body, &googleUser); err != nil {
		return nil, "","", err
	}

	existingUser,_:=u.userRepo.FindByGoogleID(context.Background(),googleUser.ID)
	if existingUser==nil{
		existingUser=&domain.User{
			GoogleID: googleUser.ID,
			Email:    googleUser.Email,
			FullName: googleUser.Name,
			Picture:  googleUser.Picture,
		}
	}

// Generate tokens
	accessToken, err := infrastructure.GenerateAccessToken(existingUser.ID.Hex(),existingUser.Role)
	if err != nil {
		return nil, "","", err
	}

	refreshToken, err := infrastructure.GenerateRefreshToken(existingUser.ID.Hex())
	if err != nil {
		return nil, "","", err
	}

	existingUser.RefreshToken = refreshToken
	if err := u.userRepo.Save(context.Background(),existingUser); err != nil {
		return nil, "", "", err
	}

	return existingUser, accessToken, refreshToken, nil
	
}