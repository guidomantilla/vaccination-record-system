package services

import (
	"context"
	"reflect"
	"testing"

	"github.com/guidomantilla/vaccination-record-system/pkg/datasource"
	"github.com/guidomantilla/vaccination-record-system/pkg/models"
)

func TestNewDefaultDrugsService(t *testing.T) {
	type args struct {
		transactionHandler datasource.TransactionHandler
	}
	tests := []struct {
		name string
		args args
		want *DefaultDrugsService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDefaultDrugsService(tt.args.transactionHandler); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDefaultDrugsService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultDrugsService_Create(t *testing.T) {
	type args struct {
		ctx  context.Context
		drug *models.Drug
	}
	tests := []struct {
		name    string
		service *DefaultDrugsService
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.service.Create(tt.args.ctx, tt.args.drug); (err != nil) != tt.wantErr {
				t.Errorf("DefaultDrugsService.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultDrugsService_Update(t *testing.T) {
	type args struct {
		ctx  context.Context
		drug *models.Drug
	}
	tests := []struct {
		name    string
		service *DefaultDrugsService
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.service.Update(tt.args.ctx, tt.args.drug); (err != nil) != tt.wantErr {
				t.Errorf("DefaultDrugsService.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultDrugsService_Delete(t *testing.T) {
	type args struct {
		ctx  context.Context
		drug *models.Drug
	}
	tests := []struct {
		name    string
		service *DefaultDrugsService
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.service.Delete(tt.args.ctx, tt.args.drug); (err != nil) != tt.wantErr {
				t.Errorf("DefaultDrugsService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultDrugsService_Find(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		service *DefaultDrugsService
		args    args
		want    []models.Drug
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.service.Find(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("DefaultDrugsService.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultDrugsService.Find() = %v, want %v", got, tt.want)
			}
		})
	}
}
