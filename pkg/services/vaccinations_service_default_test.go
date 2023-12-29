package services

import (
	"context"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	feather_commons_util "github.com/guidomantilla/go-feather-commons/pkg/util"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/guidomantilla/vaccination-record-system/pkg/datasource"
	"github.com/guidomantilla/vaccination-record-system/pkg/models"
)

func TestDefaultVaccinationsService_Create(t *testing.T) {

	type args struct {
		ctx         context.Context
		vaccination *models.Vaccination
	}
	tests := []struct {
		name    string
		service *DefaultVaccinationsService
		args    args
		wantErr bool
	}{
		{
			name: "Create vaccination",
			service: func() *DefaultVaccinationsService {

				db, mock, _ := sqlmock.New()
				dialector := mysql.New(mysql.Config{
					Conn:                      db,
					DriverName:                "mock",
					SkipInitializeWithVersion: true,
				})

				gdb, _ := gorm.Open(dialector, &gorm.Config{})

				mock.ExpectBegin()
				mock.ExpectQuery("SELECT (.+) FROM `drugs` WHERE `drugs`.`id` = ?").WithArgs("some drug id").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "approved", "min_dose", "max_dose", "available_at"}).
						AddRow("some_drug_id", "some_drug_name", true, 1, 2, time.Now()))
				mock.ExpectExec("INSERT INTO").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				transactionHandler := datasource.NewTransactionHandler(gdb)
				return NewDefaultVaccinationsService(transactionHandler)
			}(),
			args: args{
				ctx: context.TODO(),
				vaccination: &models.Vaccination{
					Name:         feather_commons_util.ValueToPtr("some vaccine name"),
					DrugId:       feather_commons_util.ValueToPtr("some drug id"),
					Dose:         feather_commons_util.ValueToPtr(int64(1)),
					DateAsString: feather_commons_util.ValueToPtr("2024-02-01"),
				},
			},
		},
		{
			name: "Create vaccination with invalid dose",
			service: func() *DefaultVaccinationsService {

				db, mock, _ := sqlmock.New()
				dialector := mysql.New(mysql.Config{
					Conn:                      db,
					DriverName:                "mock",
					SkipInitializeWithVersion: true,
				})

				gdb, _ := gorm.Open(dialector, &gorm.Config{})

				mock.ExpectBegin()
				mock.ExpectQuery("SELECT (.+) FROM `drugs` WHERE `drugs`.`id` = ?").WithArgs("some drug id").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "approved", "min_dose", "max_dose", "available_at"}).
						AddRow("some_drug_id", "some_drug_name", true, 1, 2, time.Now()))
				mock.ExpectExec("INSERT INTO").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				transactionHandler := datasource.NewTransactionHandler(gdb)
				return NewDefaultVaccinationsService(transactionHandler)
			}(),
			args: args{
				ctx: context.TODO(),
				vaccination: &models.Vaccination{
					Name:         feather_commons_util.ValueToPtr("some vaccine name"),
					DrugId:       feather_commons_util.ValueToPtr("some drug id"),
					Dose:         feather_commons_util.ValueToPtr(int64(1)),
					DateAsString: feather_commons_util.ValueToPtr("2021-02-01"),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.service.Create(tt.args.ctx, tt.args.vaccination); (err != nil) != tt.wantErr {
				t.Errorf("DefaultVaccinationsService.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
