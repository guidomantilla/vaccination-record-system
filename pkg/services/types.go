package services

import (
	"context"

	"github.com/guidomantilla/vaccination-record-system/pkg/models"
)

type TokenManager interface {
	Generate(user *models.User) (*string, error)
	Validate(tokenString string) (*models.User, error)
}

type AuthService interface {
	Signup(ctx context.Context, user *models.User) error
	Login(ctx context.Context, user *models.User) error
	Authorize(ctx context.Context, tokenString string) (*models.User, error)
}
