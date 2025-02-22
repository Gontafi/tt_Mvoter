package services

import (
	"context"
	"tt/internal/models"
	"tt/internal/repository"
)

type DynamicTableServiceInterface interface {
	AddRow(ctx context.Context, tableID int64, data map[string]interface{}) (int64, error)
	GetAllRows(ctx context.Context, tableID int64) ([]models.Row, error)
	UpdateTableRow(ctx context.Context, tableID int64, rowID int64, data map[string]interface{}) error
	RemoveTableRow(ctx context.Context, tableID int64, rowID int64) error
	CreateTable(ctx context.Context, tableName string) (int64, error)
	GetTables(ctx context.Context) ([]models.Table, error)
}

type DynamicDataService struct {
	repo *repository.DynamicDataRepository
}

func NewDynamicDataService(repo *repository.DynamicDataRepository) *DynamicDataService {
	return &DynamicDataService{repo: repo}
}

func (service *DynamicDataService) AddRow(ctx context.Context, tableID int64, data map[string]interface{}) (int64, error) {
	return service.repo.CreateRow(ctx, tableID, data)
}

func (service *DynamicDataService) GetAllRows(ctx context.Context, tableID int64) ([]models.Row, error) {
	rows, err := service.repo.GetRows(ctx, tableID)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (service *DynamicDataService) UpdateTableRow(ctx context.Context, tableID int64, rowID int64, data map[string]interface{}) error {
	return service.repo.UpdateRow(ctx, tableID, rowID, data)
}

func (service *DynamicDataService) RemoveTableRow(ctx context.Context, tableID int64, rowID int64) error {
	return service.repo.DeleteRow(ctx, tableID, rowID)
}

func (service *DynamicDataService) CreateTable(ctx context.Context, tableName string) (int64, error) {
	return service.repo.CreateTable(ctx, tableName)
}

func (service *DynamicDataService) GetTables(ctx context.Context) ([]models.Table, error) {
	return service.repo.GetTables(ctx)
}
