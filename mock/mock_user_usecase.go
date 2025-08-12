package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/sol-tad/Blog-post-Api/domain"
)

// MockUserUsecase implements usecase.UserUsecaseInterface for tests.
type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) Register(ctx context.Context, user domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserUsecase) VerifyOTP(ctx context.Context, email, otp string) error {
	args := m.Called(ctx, email, otp)
	return args.Error(0)
}

func (m *MockUserUsecase) Login(ctx context.Context, username, password string) (string, string, error) {
	args := m.Called(ctx, username, password)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockUserUsecase) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	args := m.Called(ctx, refreshToken)
	return args.String(0), args.Error(1)
}

func (m *MockUserUsecase) Logout(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserUsecase) SendResetOTP(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockUserUsecase) ResetPassword(ctx context.Context, email, otp, newPassword string) error {
	args := m.Called(ctx, email, otp, newPassword)
	return args.Error(0)
}

func (m *MockUserUsecase) PromoteUser(ctx context.Context, adminID string, targetUserID string) error {
	args := m.Called(ctx, adminID, targetUserID)
	return args.Error(0)
}

func (m *MockUserUsecase) DemoteUser(ctx context.Context, adminID string, targetUserID string) error {
	args := m.Called(ctx, adminID, targetUserID)
	return args.Error(0)
}

func (m *MockUserUsecase) UpdateProfile(ctx context.Context, userID string, updated domain.User) (domain.User, error) {
	args := m.Called(ctx, userID, updated)
	u := domain.User{}
	if tmp := args.Get(0); tmp != nil {
		u = tmp.(domain.User)
	}
	return u, args.Error(1)
}
