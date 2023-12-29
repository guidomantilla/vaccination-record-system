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

type DefaultVaccinationsEndpoint struct {
	vaccinationsService services.VaccinationsService
}

func NewDefaultVaccinationsEndpoint(vaccinationsService services.VaccinationsService) *DefaultVaccinationsEndpoint {
	return &DefaultVaccinationsEndpoint{
		vaccinationsService: vaccinationsService,
	}
}

func (endpoint *DefaultVaccinationsEndpoint) Create(ctx *gin.Context) {
	var err error
	var vaccinationToSave *models.Vaccination
	if err = ctx.ShouldBindJSON(&vaccinationToSave); err != nil {
		ex := feather_web_rest.BadRequestException("error unmarshalling request json to object")
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	if errs := endpoint.validate(vaccinationToSave); errs != nil {
		ex := feather_web_rest.BadRequestException("error validating the object", errs...)
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	if err = endpoint.vaccinationsService.Create(ctx.Request.Context(), vaccinationToSave); err != nil {
		ex := feather_web_rest.UnauthorizedException(err.Error())
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	ctx.JSON(http.StatusCreated, vaccinationToSave)
}

func (endpoint *DefaultVaccinationsEndpoint) Update(ctx *gin.Context) {
	var err error

	id := ctx.Params.ByName("id")
	if id == "" {
		ex := feather_web_rest.BadRequestException("object id not defined in url path")
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	var vaccinationToSave *models.Vaccination
	if err = ctx.ShouldBindJSON(&vaccinationToSave); err != nil {
		ex := feather_web_rest.BadRequestException("error unmarshalling request json to object")
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	if errs := endpoint.validate(vaccinationToSave); errs != nil {
		ex := feather_web_rest.BadRequestException("error validating the object", errs...)
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	vaccinationToSave.Id = &id

	if err = endpoint.vaccinationsService.Update(ctx.Request.Context(), vaccinationToSave); err != nil {
		ex := feather_web_rest.UnauthorizedException(err.Error())
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	ctx.JSON(http.StatusOK, vaccinationToSave)
}

func (endpoint *DefaultVaccinationsEndpoint) Delete(ctx *gin.Context) {
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

	vaccinationToDelete := &models.Vaccination{
		Id: &id,
	}

	if err = endpoint.vaccinationsService.Delete(ctx.Request.Context(), vaccinationToDelete); err != nil {
		ex := feather_web_rest.UnauthorizedException(err.Error())
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	ctx.JSON(http.StatusOK, vaccinationToDelete)
}

func (endpoint *DefaultVaccinationsEndpoint) Find(ctx *gin.Context) {
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

	var vaccinations []models.Vaccination
	if vaccinations, err = endpoint.vaccinationsService.Find(ctx.Request.Context()); err != nil {
		ex := feather_web_rest.UnauthorizedException(err.Error())
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	ctx.JSON(http.StatusOK, vaccinations)
}

//

func (endpoint *DefaultVaccinationsEndpoint) validate(vaccinationToSave *models.Vaccination) []error {

	var errors []error

	if err := feather_commons_validation.ValidateFieldMustBeUndefined("this", "id", vaccinationToSave.Id); err != nil {
		errors = append(errors, err)
	}
	if err := feather_commons_validation.ValidateFieldIsRequired("this", "name", vaccinationToSave.Name); err != nil {
		errors = append(errors, err)
	}

	if err := feather_commons_validation.ValidateFieldIsRequired("this", "drug_id", vaccinationToSave.DrugId); err != nil {
		errors = append(errors, err)
	}

	if err := feather_commons_validation.ValidateStructMustBeUndefined("this", "drug", vaccinationToSave.Drug); err != nil {
		errors = append(errors, err)
	}

	if err := feather_commons_validation.ValidateFieldIsRequired("this", "dose", vaccinationToSave.Dose); err != nil {
		errors = append(errors, err)
	}

	if err := feather_commons_validation.ValidateFieldIsRequired("this", "date", vaccinationToSave.DateAsString); err != nil {
		errors = append(errors, err)
	}

	return errors
}
