package services

import (
	"context"
	"errors"

	feather_security "github.com/guidomantilla/go-feather-security/pkg/security"
	"gorm.io/gorm"

	"github.com/guidomantilla/vaccination-record-system/pkg/datasource"
	"github.com/guidomantilla/vaccination-record-system/pkg/models"
)

type DefaultAuthService struct {
	transactionHandler datasource.TransactionHandler
	tokenManager       TokenManager
	passwordManager    feather_security.PasswordManager
}

func NewDefaultAuthService(transactionHandler datasource.TransactionHandler, tokenManager TokenManager, passwordManager feather_security.PasswordManager) *DefaultAuthService {
	return &DefaultAuthService{
		transactionHandler: transactionHandler,
		tokenManager:       tokenManager,
		passwordManager:    passwordManager,
	}
}

func (service *DefaultAuthService) Signup(ctx context.Context, user *models.User) error {
	return service.transactionHandler.HandleTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {

		var err error
		var dbUser models.User

		err = tx.Where("email = ?", user.Email).First(&dbUser).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if dbUser.Email != nil {
			return feather_security.ErrAccountExistingUsername
		}

		encodedPassword, err := service.passwordManager.Encode(*user.Password)
		if err != nil {
			return err
		}

		dbUser.Name, dbUser.Email, dbUser.Password = user.Name, user.Email, encodedPassword
		tx.Save(dbUser)

		return nil
	})
}

func (service *DefaultAuthService) Login(ctx context.Context, user *models.User) error {
	return service.transactionHandler.HandleTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {

		var err error
		var dbUser models.User

		if err = tx.Where("email = ?", user.Email).First(&dbUser).Error; err != nil {
			return feather_security.ErrAuthenticationFailed(err)
		}

		var needsUpgrade *bool
		if needsUpgrade, err = service.passwordManager.UpgradeEncoding(*(dbUser.Password)); err != nil || *(needsUpgrade) {
			return feather_security.ErrAuthenticationFailed(feather_security.ErrAccountExpiredPassword)
		}

		var matches *bool
		if matches, err = service.passwordManager.Matches(*(dbUser.Password), *user.Password); err != nil || !*(matches) {
			return feather_security.ErrAuthenticationFailed(feather_security.ErrAccountInvalidPassword)
		}

		user.Password = nil
		if user.Token, err = service.tokenManager.Generate(user); err != nil {
			return feather_security.ErrAuthenticationFailed(err)
		}

		return nil
	})
}

func (service *DefaultAuthService) Authorize(ctx context.Context, tokenString string) (*models.User, error) {
	var err error
	var user *models.User
	err = service.transactionHandler.HandleTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {

		if user, err = service.tokenManager.Validate(tokenString); err != nil {
			return err
		}

		dbUser := &models.User{}
		if err = tx.Where("email = ?", user.Email).First(dbUser).Error; err != nil {
			return feather_security.ErrAuthorizationFailed(err)
		}

		dbUser.Password, dbUser.Token = nil, &tokenString
		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}
