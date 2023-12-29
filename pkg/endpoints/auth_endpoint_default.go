package endpoints

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	feather_commons_validation "github.com/guidomantilla/go-feather-commons/pkg/validation"
	feather_web_rest "github.com/guidomantilla/go-feather-web/pkg/rest"

	"github.com/guidomantilla/vaccination-record-system/pkg/models"
	"github.com/guidomantilla/vaccination-record-system/pkg/services"
)

type DefaultAuthEndpoint struct {
	authService services.AuthService
}

func NewDefaultAuthEndpoint(authService services.AuthService) *DefaultAuthEndpoint {
	return &DefaultAuthEndpoint{
		authService: authService,
	}
}

func (endpoint *DefaultAuthEndpoint) Login(ctx *gin.Context) {

	var err error
	var userLoginIn *models.User
	if err = ctx.ShouldBindJSON(&userLoginIn); err != nil {
		ex := feather_web_rest.BadRequestException("error unmarshalling request json to object")
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	if errs := endpoint.validateLogin(userLoginIn); errs != nil {
		ex := feather_web_rest.BadRequestException("error validating the principal", errs...)
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	if err = endpoint.authService.Login(ctx.Request.Context(), userLoginIn); err != nil {
		ex := feather_web_rest.UnauthorizedException(err.Error())
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	ctx.JSON(http.StatusOK, userLoginIn)
}

func (endpoint *DefaultAuthEndpoint) Signup(ctx *gin.Context) {

	var err error
	var userSigningUp *models.User
	if err = ctx.ShouldBindJSON(&userSigningUp); err != nil {
		ex := feather_web_rest.BadRequestException("error unmarshalling request json to object")
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	if errs := endpoint.validateSignup(userSigningUp); errs != nil {
		ex := feather_web_rest.BadRequestException("error validating the principal", errs...)
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	if err = endpoint.authService.Signup(ctx.Request.Context(), userSigningUp); err != nil {
		ex := feather_web_rest.InternalServerErrorException(err.Error())
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	userSigningUp.Password = nil
	ctx.JSON(http.StatusCreated, userSigningUp)
}

func (endpoint *DefaultAuthEndpoint) Authorize(ctx *gin.Context) {

	header := ctx.Request.Header.Get("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		ex := feather_web_rest.UnauthorizedException("invalid authorization header")
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	splits := strings.Split(header, " ")
	if len(splits) != 2 {
		ex := feather_web_rest.UnauthorizedException("invalid authorization header")
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}
	token := splits[1]

	var err error
	var userToAuthorize *models.User
	//ctxWithResource := context.WithValue(ctx.Request.Context(), ResourceCtxKey{}, strings.Join(resource, " "))
	if userToAuthorize, err = endpoint.authService.Authorize(ctx.Request.Context(), token); err != nil {
		ex := feather_web_rest.UnauthorizedException(err.Error())
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	ctx.Set("principal", userToAuthorize)
	ctx.Next()
}

func (endpoint *DefaultAuthEndpoint) validateLogin(userLoginIn *models.User) []error {

	var errors []error

	if err := feather_commons_validation.ValidateFieldMustBeUndefined("this", "id", userLoginIn.Id); err != nil {
		errors = append(errors, err)
	}
	if err := feather_commons_validation.ValidateFieldMustBeUndefined("this", "name", userLoginIn.Name); err != nil {
		errors = append(errors, err)
	}

	if err := feather_commons_validation.ValidateFieldIsRequired("this", "email", userLoginIn.Email); err != nil {
		errors = append(errors, err)
	}

	if err := feather_commons_validation.ValidateFieldIsRequired("this", "password", userLoginIn.Password); err != nil {
		errors = append(errors, err)
	}

	if err := feather_commons_validation.ValidateFieldMustBeUndefined("this", "token", userLoginIn.Token); err != nil {
		errors = append(errors, err)
	}

	return errors
}

func (endpoint *DefaultAuthEndpoint) validateSignup(userSigningUp *models.User) []error {

	var errors []error

	if err := feather_commons_validation.ValidateFieldMustBeUndefined("this", "id", userSigningUp.Id); err != nil {
		errors = append(errors, err)
	}
	if err := feather_commons_validation.ValidateFieldIsRequired("this", "name", userSigningUp.Name); err != nil {
		errors = append(errors, err)
	}

	if err := feather_commons_validation.ValidateFieldIsRequired("this", "email", userSigningUp.Email); err != nil {
		errors = append(errors, err)
	}

	if err := feather_commons_validation.ValidateFieldIsRequired("this", "password", userSigningUp.Password); err != nil {
		errors = append(errors, err)
	}

	if err := feather_commons_validation.ValidateFieldMustBeUndefined("this", "token", userSigningUp.Token); err != nil {
		errors = append(errors, err)
	}

	return errors
}
