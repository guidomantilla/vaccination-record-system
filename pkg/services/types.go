package services

import (
	"context"

	"github.com/guidomantilla/vaccination-record-system/pkg/models"
)

var (
	_ TokenManager        = (*JwtTokenManager)(nil)
	_ AuthService         = (*DefaultAuthService)(nil)
	_ DrugsService        = (*DefaultDrugsService)(nil)
	_ VaccinationsService = (*DefaultVaccinationsService)(nil)
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

type DrugsService interface {
	Create(ctx context.Context, drug *models.Drug) error
	Update(ctx context.Context, drug *models.Drug) error
	Delete(ctx context.Context, drug *models.Drug) error
	Find(ctx context.Context) ([]models.Drug, error)
}

type VaccinationsService interface {
	Create(ctx context.Context, vaccination *models.Vaccination) error
	Update(ctx context.Context, vaccination *models.Vaccination) error
	Delete(ctx context.Context, vaccination *models.Vaccination) error
	Find(ctx context.Context) ([]models.Vaccination, error)
}
