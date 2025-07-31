package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sol-tad/Blog-post-Api/domian"
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

func (uc *UserController) Register(c *gin.Context){
	var user domian.User

	if err:=c.ShouldBindJSON(&user); err !=nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}

	if err:=uc.UserUsecase.Register(user);err!=nil{
		 c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
	}
	 c.JSON(http.StatusOK, gin.H{"message": "user registered successfully"})
}

func (uc *UserController) Login(c *gin.Context){
	
}