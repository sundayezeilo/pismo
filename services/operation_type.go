package services

import (
	"context"

	"github.com/sundayezeilo/pismo/models"
	"github.com/sundayezeilo/pismo/repositories"
)

type OperationTypeService interface {
	GetOpTypeByID(context.Context, int) (*models.OperationType, error)
}

type opTypeService struct {
	repo repositories.OperationTypeRepository
}

func NewOpTypeService(repo repositories.OperationTypeRepository) OperationTypeService {
	return &opTypeService{repo}
}

func (srv *opTypeService) GetOpTypeByID(ctx context.Context, oTypeID int) (*models.OperationType, error) {
	oType, err := srv.repo.GetOpTypeByID(ctx, oTypeID)

	if err != nil {
		return nil, err
	}

	return oType, nil
}
