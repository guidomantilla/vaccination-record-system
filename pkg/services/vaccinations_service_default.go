package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	feather_commons_util "github.com/guidomantilla/go-feather-commons/pkg/util"
	"gorm.io/gorm"

	"github.com/guidomantilla/vaccination-record-system/pkg/datasource"
	"github.com/guidomantilla/vaccination-record-system/pkg/models"
)

type DefaultVaccinationsService struct {
	transactionHandler datasource.TransactionHandler
}

func NewDefaultVaccinationsService(transactionHandler datasource.TransactionHandler) *DefaultVaccinationsService {
	return &DefaultVaccinationsService{
		transactionHandler: transactionHandler,
	}
}

func (service *DefaultVaccinationsService) Create(ctx context.Context, vaccination *models.Vaccination) error {
	return service.transactionHandler.HandleTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {

		var err error
		var dateAtAsDate time.Time
		if dateAtAsDate, err = time.Parse("2006-01-02", *vaccination.DateAsString); err != nil {
			return err
		}

		vaccination.Id = feather_commons_util.ValueToPtr(uuid.New().String())
		vaccination.DateAsDate = &dateAtAsDate

		if err = tx.Find(&vaccination.Drug, vaccination.DrugId).Error; err != nil {
			return err
		}

		if IsVaccinationDoseOk(vaccination) && IsVaccinationAfterAllowedDate(vaccination) {
			return tx.Create(vaccination).Error
		}

		return errors.New("invalid vaccination")
	})
}

func IsVaccinationDoseOk(vaccination *models.Vaccination) bool {
	return vaccination.Drug.Id != nil && *vaccination.Dose >= *vaccination.Drug.MinDose && *vaccination.Dose <= *vaccination.Drug.MaxDose
}

func IsVaccinationAfterAllowedDate(vaccination *models.Vaccination) bool {
	return vaccination.Drug.Id != nil && vaccination.DateAsDate.After(*vaccination.Drug.AvailableAtAsDate)
}

func (service *DefaultVaccinationsService) Update(ctx context.Context, vaccination *models.Vaccination) error {
	return service.transactionHandler.HandleTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {

		var err error
		var dateAtAsDate time.Time
		if dateAtAsDate, err = time.Parse("2006-01-02", *vaccination.DateAsString); err != nil {
			return err
		}

		vaccination.DateAsDate = &dateAtAsDate

		if err = tx.Find(&vaccination.Drug, vaccination.DrugId).Error; err != nil {
			return err
		}

		if IsVaccinationDoseOk(vaccination) && IsVaccinationAfterAllowedDate(vaccination) {
			return tx.Create(vaccination).Error
		}

		return errors.New("invalid vaccination")
	})
}

func (service *DefaultVaccinationsService) Delete(ctx context.Context, vaccination *models.Vaccination) error {
	return service.transactionHandler.HandleTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {
		return tx.Delete(vaccination).Where("name = ?", vaccination.Name).Error
	})
}

func (service *DefaultVaccinationsService) Find(ctx context.Context) ([]models.Vaccination, error) {
	var err error
	var vaccinations []models.Vaccination
	err = service.transactionHandler.HandleTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {
		return tx.Preload("Drug").Order("name asc").Find(&vaccinations).Error
	})

	if err != nil {
		return nil, err
	}

	for i, vaccination := range vaccinations {
		vaccinations[i].DrugId = nil
		vaccinations[i].DateAsString = feather_commons_util.ValueToPtr(vaccination.DateAsDate.Format("2006-01-02"))
		vaccinations[i].Drug.AvailableAtAsString = feather_commons_util.ValueToPtr(vaccinations[i].Drug.AvailableAtAsDate.Format("2006-01-02"))
	}
	return vaccinations, nil
}
