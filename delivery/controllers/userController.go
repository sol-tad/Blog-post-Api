package controllers

import (
	"context"
	"log"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sol-tad/Blog-post-Api/domain"
	"github.com/sol-tad/Blog-post-Api/usecase"
)

type UserController struct {
	UserUsecase *usecase.UserUsecase
}

func NewUserController(userUsecase *usecase.UserUsecase) *UserController {
	return &UserController{
		UserUsecase: userUsecase,
	}
}
func (uc *UserController) Register(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uc.UserUsecase.Register(context.Background(), user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent to your email. Please verify."})
}

func (uc *UserController) VerifyOTP(c *gin.Context) {
	var input struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uc.UserUsecase.VerifyOTP(context.Background(), input.Email, input.OTP); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account verified successfully!"})
}

func (uc *UserController) Login(c *gin.Context){
	var req struct{
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err:=c.ShouldBindJSON(&req); err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}
	accessToken,refreshToken,err:=uc.UserUsecase.Login(context.Background(), req.Username, req.Password)

	    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "access_token":  accessToken,
        "refresh_token": refreshToken,
    })
}


func (uc UserController) RefreshTokenController(c *gin.Context){
	var req struct{
		RefreshToken  string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
	// log.Println("Refresh Token----------->:", req.RefreshToken)

	token,err:=uc.UserUsecase.RefreshToken(context.Background(),req.RefreshToken)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token*****"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"access_token": token})
}

func (uc UserController) Logout(c *gin.Context) {
		userID := c.GetString("id")
		log.Println("id============:", userID)

		err := uc.UserUsecase.Logout(context.Background(), userID)
			if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "logout failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})

	}
