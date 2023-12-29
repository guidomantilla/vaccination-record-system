package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	feather_commons_util "github.com/guidomantilla/go-feather-commons/pkg/util"
	"gorm.io/gorm"

	"github.com/guidomantilla/vaccination-record-system/pkg/datasource"
	"github.com/guidomantilla/vaccination-record-system/pkg/models"
)

type DefaultDrugsService struct {
	transactionHandler datasource.TransactionHandler
}

func NewDefaultDrugsService(transactionHandler datasource.TransactionHandler) *DefaultDrugsService {
	return &DefaultDrugsService{
		transactionHandler: transactionHandler,
	}
}

func (service *DefaultDrugsService) Create(ctx context.Context, drug *models.Drug) error {
	return service.transactionHandler.HandleTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {

		var err error
		var availableAtAsDate time.Time
		if availableAtAsDate, err = time.Parse("2006-01-02", *drug.AvailableAtAsString); err != nil {
			return err
		}

		drug.Id = feather_commons_util.ValueToPtr(uuid.New().String())
		drug.AvailableAtAsDate = &availableAtAsDate
		return tx.Create(drug).Error
	})
}

func (service *DefaultDrugsService) Update(ctx context.Context, drug *models.Drug) error {
	return service.transactionHandler.HandleTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {

		var err error
		var availableAtAsDate time.Time
		if availableAtAsDate, err = time.Parse("2006-01-02", *drug.AvailableAtAsString); err != nil {
			return err
		}

		drug.AvailableAtAsDate = &availableAtAsDate
		return tx.Save(drug).Error
	})
}

func (service *DefaultDrugsService) Delete(ctx context.Context, drug *models.Drug) error {
	return service.transactionHandler.HandleTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {
		return tx.Delete(drug).Where("name = ?", drug.Name).Error
	})
}

func (service *DefaultDrugsService) Find(ctx context.Context) ([]models.Drug, error) {

	var err error
	var drugs []models.Drug
	err = service.transactionHandler.HandleTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {
		return tx.Order("name asc").Find(&drugs).Error
	})

	if err != nil {
		return nil, err
	}

	for i, drug := range drugs {
		drugs[i].AvailableAtAsString = feather_commons_util.ValueToPtr(drug.AvailableAtAsDate.Format("2006-01-02"))
	}
	return drugs, nil
}
