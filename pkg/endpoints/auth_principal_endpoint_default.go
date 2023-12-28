package endpoints

import (
	"github.com/gin-gonic/gin"
	feather_commons_validation "github.com/guidomantilla/go-feather-commons/pkg/validation"
	feather_security "github.com/guidomantilla/go-feather-security/pkg/security"
	feather_web_rest "github.com/guidomantilla/go-feather-web/pkg/rest"
	"net/http"
	"reflect"
)

type DefaultAuthPrincipalEndpoint struct {
	principalManager feather_security.PrincipalManager
}

func NewDefaultAuthPrincipalEndpoint(principalManager feather_security.PrincipalManager) *DefaultAuthPrincipalEndpoint {
	return &DefaultAuthPrincipalEndpoint{
		principalManager: principalManager,
	}
}

func (endpoint *DefaultAuthPrincipalEndpoint) Signup(ctx *gin.Context) {

	var err error
	var principal *feather_security.Principal
	if err = ctx.ShouldBindJSON(&principal); err != nil {
		ex := feather_web_rest.BadRequestException("error unmarshalling request json to object")
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	if errs := endpoint.validateUpsert(principal); errs != nil {
		ex := feather_web_rest.BadRequestException("error validating the principal", errs...)
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	var exists bool
	var current *feather_security.Principal
	if current, exists = feather_security.GetPrincipalFromContext(ctx); !exists {
		ex := feather_web_rest.NotFoundException("principal not found in context")
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	if reflect.DeepEqual(current.Username, principal.Username) {
		ex := feather_web_rest.BadRequestException("authorized username cannot be the same as the new username")
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	if err = endpoint.principalManager.Create(ctx.Request.Context(), principal); err != nil {
		ex := feather_web_rest.InternalServerErrorException(err.Error())
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	principal.Password, principal.Passphrase = nil, nil
	ctx.JSON(http.StatusCreated, principal)
}

func (endpoint *DefaultAuthPrincipalEndpoint) validateUpsert(principal *feather_security.Principal) []error {

	var errors []error

	if err := feather_commons_validation.ValidateFieldIsRequired("this", "username", principal.Username); err != nil {
		errors = append(errors, err)
	}

	if err := feather_commons_validation.ValidateFieldIsRequired("this", "role", principal.Role); err != nil {
		errors = append(errors, err)
	}

	if err := feather_commons_validation.ValidateFieldIsRequired("this", "password", principal.Password); err != nil {
		errors = append(errors, err)
	}

	if err := feather_commons_validation.ValidateFieldIsRequired("this", "passphrase", principal.Passphrase); err != nil {
		errors = append(errors, err)
	}

	if err := feather_commons_validation.ValidateFieldMustBeUndefined("this", "enabled", principal.Enabled); err != nil {
		errors = append(errors, err)
	}

	if err := feather_commons_validation.ValidateFieldMustBeUndefined("this", "non_locked", principal.NonLocked); err != nil {
		errors = append(errors, err)
	}

	if err := feather_commons_validation.ValidateFieldMustBeUndefined("this", "non_expired", principal.NonExpired); err != nil {
		errors = append(errors, err)
	}

	if err := feather_commons_validation.ValidateFieldMustBeUndefined("this", "password_non_expired", principal.PasswordNonExpired); err != nil {
		errors = append(errors, err)
	}

	if err := feather_commons_validation.ValidateFieldMustBeUndefined("this", "signup_done", principal.SignUpDone); err != nil {
		errors = append(errors, err)
	}

	if err := feather_commons_validation.ValidateStructIsRequired("this", "resources", principal.Resources); err != nil {
		errors = append(errors, err)
		return errors
	}

	if err := feather_commons_validation.ValidateFieldMustBeUndefined("this", "token", principal.Token); err != nil {
		errors = append(errors, err)
	}

	return errors
}
