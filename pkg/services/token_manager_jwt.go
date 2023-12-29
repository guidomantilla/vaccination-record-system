package services

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	feather_commons_util "github.com/guidomantilla/go-feather-commons/pkg/util"
	feather_security "github.com/guidomantilla/go-feather-security/pkg/security"

	"github.com/guidomantilla/vaccination-record-system/pkg/models"
)

type DefaultClaims struct {
	jwt.RegisteredClaims
	models.User
}

type JwtTokenManagerOption func(tokenManager *JwtTokenManager)

type JwtTokenManager struct {
	issuer        string
	timeout       time.Duration
	signingKey    any
	verifyingKey  any
	signingMethod jwt.SigningMethod
}

func NewJwtTokenManager(options ...JwtTokenManagerOption) *JwtTokenManager {

	tokenManager := &JwtTokenManager{
		issuer:        "",
		timeout:       time.Hour * 24,
		signingKey:    "some_long_signing_key",
		verifyingKey:  "some_long_verifying_key",
		signingMethod: jwt.SigningMethodHS512,
	}

	for _, opt := range options {
		opt(tokenManager)
	}

	return tokenManager
}

func WithIssuer(issuer string) JwtTokenManagerOption {
	return func(tokenManager *JwtTokenManager) {
		tokenManager.issuer = issuer
	}
}

func WithTimeout(timeout time.Duration) JwtTokenManagerOption {
	return func(tokenManager *JwtTokenManager) {
		tokenManager.timeout = timeout
	}
}

func WithSigningMethod(signingMethod jwt.SigningMethod) JwtTokenManagerOption {
	return func(tokenManager *JwtTokenManager) {
		tokenManager.signingMethod = signingMethod
	}
}

func WithSigningKey(signingKey any) JwtTokenManagerOption {
	return func(tokenManager *JwtTokenManager) {
		tokenManager.signingKey = signingKey
	}
}

func WithVerifyingKey(verifyingKey any) JwtTokenManagerOption {
	return func(tokenManager *JwtTokenManager) {
		tokenManager.verifyingKey = verifyingKey
	}
}

func (manager *JwtTokenManager) Generate(user *models.User) (*string, error) {

	claims := &DefaultClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    manager.issuer,
			Subject:   *user.Email,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(manager.timeout)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		User: *user,
	}

	token := jwt.NewWithClaims(manager.signingMethod, claims)

	var err error
	var tokenString string
	if tokenString, err = token.SignedString(manager.signingKey); err != nil {
		return nil, feather_security.ErrTokenGenerationFailed(err)
	}

	return &tokenString, nil
}

func (manager *JwtTokenManager) Validate(tokenString string) (*models.User, error) {

	getKeyFunc := func(token *jwt.Token) (any, error) {
		return manager.verifyingKey, nil
	}

	parserOptions := []jwt.ParserOption{
		jwt.WithIssuer(manager.issuer),
		jwt.WithValidMethods([]string{manager.signingMethod.Alg()}),
	}

	var err error
	var token *jwt.Token
	if token, err = jwt.Parse(tokenString, getKeyFunc, parserOptions...); err != nil {
		return nil, feather_security.ErrTokenValidationFailed(feather_security.ErrTokenFailedParsing, err)
	}

	if !token.Valid {
		return nil, feather_security.ErrTokenValidationFailed(feather_security.ErrTokenInvalid)
	}

	var ok bool
	var mapClaims jwt.MapClaims
	if mapClaims, ok = token.Claims.(jwt.MapClaims); !ok {
		return nil, feather_security.ErrTokenValidationFailed(feather_security.ErrTokenEmptyClaims)
	}

	var value any
	if value, ok = mapClaims["email"]; !ok {
		return nil, feather_security.ErrTokenValidationFailed(errors.New("token email claim is empty"))
	}

	var email string
	if email, ok = value.(string); !ok {
		return nil, feather_security.ErrTokenValidationFailed(errors.New("token email claim is empty"))
	}

	user := &models.User{
		Email: feather_commons_util.ValueToPtr(email),
	}

	return user, nil
}
