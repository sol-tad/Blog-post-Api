package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/sol-tad/Blog-post-Api/delivery/controllers"
	"github.com/sol-tad/Blog-post-Api/domain"
	mocks "github.com/sol-tad/Blog-post-Api/mock"
)

func setupRouter(mockUsecase *mocks.MockUserUsecase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	ctrl := controllers.NewUserController(mockUsecase)
	router := gin.Default()
	router.Use(func(c *gin.Context) {
        c.Set("id", "uid9")
        c.Next()
    })
	router.Use(func(c *gin.Context) {
		c.Set("id", "admin1")
		c.Next()
	})
	// register routes as in your app
	router.POST("/register", ctrl.Register)
	router.POST("/verify", ctrl.VerifyOTP)
	router.POST("/login", ctrl.Login)
	router.POST("/refresh", ctrl.RefreshTokenController)

	// middleware to set user id for routes that need it
	// router.Use(func(c *gin.Context) {
	// 	// do nothing by default, some tests will set via header
	// 	c.Next()
	// })

	router.POST("/logout", func(c *gin.Context) { ctrl.Logout(c) })

	router.POST("/send-reset", ctrl.SendResetOTP)
	router.POST("/reset-password", ctrl.ResetPassword)

	router.PUT("/promote/:id", ctrl.PromoteUser)
	router.PUT("/demote/:id", ctrl.DemoteUser)
	router.PUT("/profile", ctrl.UpdateProfile)

	return router
}

func TestRegister_Success(t *testing.T) {
	mockU := &mocks.MockUserUsecase{}
	// Expect Register to be called with any context and a domain.User
	mockU.On("Register", mock.Anything, mock.MatchedBy(func(u interface{}) bool {
		_, ok := u.(domain.User)
		return ok
	})).Return(nil)

	router := setupRouter(mockU)

	user := domain.User{Username: "u1", Email: "st@gmail.com", Password: "pass123"}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "OTP sent to your email")
	mockU.AssertExpectations(t)
}

func TestRegister_BadRequest(t *testing.T) {
	mockU := &mocks.MockUserUsecase{}
	router := setupRouter(mockU)
	
	mockU.On("Register", mock.Anything, mock.MatchedBy(func(u interface{}) bool {
		_, ok := u.(domain.User)
		return ok
	})).Return(nil)

	// invalid JSON/missing fields
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(`{"email":"bad"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestVerifyOTP_Success(t *testing.T) {
	mockU := &mocks.MockUserUsecase{}
	mockU.On("VerifyOTP", mock.Anything, "a@test.com", "123456").Return(nil)

	router := setupRouter(mockU)

	payload := map[string]string{"email": "a@test.com", "otp": "123456"}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/verify", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	mockU.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	mockU := &mocks.MockUserUsecase{}
	mockU.On("Login", mock.Anything, "user1", "pass1").Return("access_tok", "refresh_tok", nil)

	router := setupRouter(mockU)
	payload := map[string]string{"username": "user1", "password": "pass1"}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "access_tok", resp["access_token"])
	assert.Equal(t, "refresh_tok", resp["refresh_token"])
	mockU.AssertExpectations(t)
}

func TestLogin_Unauthorized(t *testing.T) {
	mockU := &mocks.MockUserUsecase{}
	mockU.On("Login", mock.Anything, "user1", "wrong").Return("", "", errors.New("invalid"))

	router := setupRouter(mockU)
	payload := map[string]string{"username": "user1", "password": "wrong"}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockU.AssertExpectations(t)
}

func TestRefreshToken_Success(t *testing.T) {
	mockU := &mocks.MockUserUsecase{}
	mockU.On("RefreshToken", mock.Anything, "refresh_tok").Return("new_access", nil)

	router := setupRouter(mockU)
	payload := map[string]string{"refresh_token": "refresh_tok"}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	mockU.AssertExpectations(t)
}

func TestLogout_Success(t *testing.T) {
	mockU := &mocks.MockUserUsecase{}
	// this route uses c.GetString("id"), so we set up a small middleware to set the id
	router := setupRouter(mockU)
	router.Use(func(c *gin.Context) {
		c.Set("id", "uid123")
		c.Next()
	})
	mockU.On("Logout", mock.Anything, "uid9").Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockU.AssertExpectations(t)
}

func TestSendResetOTP_Success(t *testing.T) {
	mockU := &mocks.MockUserUsecase{}
	mockU.On("SendResetOTP", mock.Anything, "r@test.com").Return(nil)

	router := setupRouter(mockU)
	payload := map[string]string{"email": "r@test.com"}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/send-reset", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockU.AssertExpectations(t)
}

func TestResetPassword_Success(t *testing.T) {
	mockU := &mocks.MockUserUsecase{}
	mockU.On("ResetPassword", mock.Anything, "r@test.com", "111111", "newpass").Return(nil)

	router := setupRouter(mockU)
	payload := map[string]string{"email": "r@test.com", "otp": "111111", "new_password": "newpass"}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/reset-password", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockU.AssertExpectations(t)
}

func TestPromote_Demote_Success(t *testing.T) {
	mockU := &mocks.MockUserUsecase{}
	router := setupRouter(mockU)

	// middleware to set admin id
	

	// Promote
	mockU.On("PromoteUser", mock.Anything, "admin1", "target1").Return(nil)
	req := httptest.NewRequest(http.MethodPut, "/promote/target1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Demote
	mockU.On("DemoteUser", mock.Anything, "admin1", "target1").Return(nil)
	req2 := httptest.NewRequest(http.MethodPut, "/demote/target1", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	mockU.AssertExpectations(t)
}

func TestUpdateProfile_Success(t *testing.T) {
    mockU := &mocks.MockUserUsecase{}
    router := setupRouter(mockU)
    
    // Middleware that properly sets ID
    router.Use(func(c *gin.Context) {
        c.Set("id", "uid9") // Set as string
        c.Next()
    })

    updated := domain.User{Username: "name", Email: "e@x.com"}
    
    // Mock expects the actual call signature
    mockU.On("UpdateProfile", mock.AnythingOfType("*gin.Context"), "uid9", updated).
        Return(updated, nil)

    b, _ := json.Marshal(updated)
    req := httptest.NewRequest(http.MethodPut, "/profile", bytes.NewBuffer(b))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    var got domain.User
    _ = json.Unmarshal(w.Body.Bytes(), &got)
    assert.Equal(t, "name", got.Username)
    assert.Equal(t, "e@x.com", got.Email)
    mockU.AssertExpectations(t)
}

