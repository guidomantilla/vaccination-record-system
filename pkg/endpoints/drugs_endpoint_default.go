package endpoints

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	feather_commons_validation "github.com/guidomantilla/go-feather-commons/pkg/validation"
	feather_web_rest "github.com/guidomantilla/go-feather-web/pkg/rest"

	"github.com/guidomantilla/vaccination-record-system/pkg/models"
	"github.com/guidomantilla/vaccination-record-system/pkg/services"
)

type DefaultDrugsEndpoint struct {
	drugsService services.DrugsService
}

func NewDefaultDrugsEndpoint(drugsService services.DrugsService) *DefaultDrugsEndpoint {
	return &DefaultDrugsEndpoint{
		drugsService: drugsService,
	}
}

func (endpoint *DefaultDrugsEndpoint) Create(ctx *gin.Context) {

	var err error
	var drugToSave *models.Drug
	if err = ctx.ShouldBindJSON(&drugToSave); err != nil {
		ex := feather_web_rest.BadRequestException("error unmarshalling request json to object")
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	if errs := endpoint.validateCreate(drugToSave); errs != nil {
		ex := feather_web_rest.BadRequestException("error validating the object", errs...)
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	if err = endpoint.drugsService.Create(ctx.Request.Context(), drugToSave); err != nil {
		ex := feather_web_rest.UnauthorizedException(err.Error())
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	ctx.JSON(http.StatusCreated, drugToSave)
}

func (endpoint *DefaultDrugsEndpoint) Update(ctx *gin.Context) {
	var err error

	id := ctx.Params.ByName("id")
	if id == "" {
		ex := feather_web_rest.BadRequestException("object id not defined in url path")
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	var drugToSave *models.Drug
	if err = ctx.ShouldBindJSON(&drugToSave); err != nil {
		ex := feather_web_rest.BadRequestException("error unmarshalling request json to object")
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	if errs := endpoint.validateCreate(drugToSave); errs != nil {
		ex := feather_web_rest.BadRequestException("error validating the object", errs...)
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	drugToSave.Id = &id

	if err = endpoint.drugsService.Update(ctx.Request.Context(), drugToSave); err != nil {
		ex := feather_web_rest.UnauthorizedException(err.Error())
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	ctx.JSON(http.StatusOK, drugToSave)
}

func (endpoint *DefaultDrugsEndpoint) Delete(ctx *gin.Context) {
	var err error

	var body []byte
	if body, err = io.ReadAll(ctx.Request.Body); err != nil {
		ex := feather_web_rest.BadRequestException("error reading body")
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	if len(body) != 0 {
		ex := feather_web_rest.BadRequestException("body is not empty")
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	id := ctx.Params.ByName("id")
	if id == "" {
		ex := feather_web_rest.BadRequestException("object id not defined in url path")
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	drugToSave := &models.Drug{
		Id: &id,
	}

	if err = endpoint.drugsService.Delete(ctx.Request.Context(), drugToSave); err != nil {
		ex := feather_web_rest.UnauthorizedException(err.Error())
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	ctx.JSON(http.StatusOK, drugToSave)
}

func (endpoint *DefaultDrugsEndpoint) Find(ctx *gin.Context) {
	var err error

	var body []byte
	if body, err = io.ReadAll(ctx.Request.Body); err != nil {
		ex := feather_web_rest.BadRequestException("error reading body")
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	if len(body) != 0 {
		ex := feather_web_rest.BadRequestException("body is not empty")
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	var drugs []models.Drug
	if drugs, err = endpoint.drugsService.Find(ctx.Request.Context()); err != nil {
		ex := feather_web_rest.UnauthorizedException(err.Error())
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	ctx.JSON(http.StatusOK, drugs)
}

//

func (endpoint *DefaultDrugsEndpoint) validateCreate(drugToSave *models.Drug) []error {

	var errors []error

	if err := feather_commons_validation.ValidateFieldMustBeUndefined("this", "id", drugToSave.Id); err != nil {
		errors = append(errors, err)
	}
	if err := feather_commons_validation.ValidateFieldIsRequired("this", "name", drugToSave.Name); err != nil {
		errors = append(errors, err)
	}

	if err := feather_commons_validation.ValidateFieldIsRequired("this", "email", drugToSave.Approved); err != nil {
		errors = append(errors, err)
	}

	if err := feather_commons_validation.ValidateFieldIsRequired("this", "min_dose", drugToSave.MinDose); err != nil {
		errors = append(errors, err)
	}

	if err := feather_commons_validation.ValidateFieldIsRequired("this", "max_dose", drugToSave.MaxDose); err != nil {
		errors = append(errors, err)
	}

	if err := feather_commons_validation.ValidateFieldIsRequired("this", "available_at", drugToSave.AvailableAtAsString); err != nil {
		errors = append(errors, err)
	}

	return errors
}
