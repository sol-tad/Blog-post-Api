package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sol-tad/Blog-post-Api/config"
	"github.com/sol-tad/Blog-post-Api/usecase"
)





type OAuthController struct {
	OAuthUsecase *usecase.OAuthUsecase
}

func NewOAuthController(oauthUsecase *usecase.OAuthUsecase) *OAuthController{
	return &OAuthController{
		OAuthUsecase:oauthUsecase ,
	}
}

func (oauc *OAuthController) Login(c *gin.Context){
	url:=config.GoogleOAuthConfig.AuthCodeURL("state_value")
	c.Redirect(http.StatusTemporaryRedirect,url)
}

func (oauc *OAuthController) Callback(c *gin.Context){
	code:=c.Query("code")
	log.Println("**************----",code,"-----------------------****")
	user,accessToken,refreshToken,err:=oauc.OAuthUsecase.HandleGoogleCallback(code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OAuth failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OAuth login successful",
		"user":    user,
		"access_token":accessToken,
		"refresh_token":refreshToken,
	})

}