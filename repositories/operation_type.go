package repositories

import (
	"context"
	"database/sql"

	"github.com/sundayezeilo/pismo/models"
)

type OperationTypeRepository interface {
	GetOpTypeByID(context.Context, int) (*models.OperationType, error)
}

type operationTypeRepository struct {
	db *sql.DB
}

func NewOpTypeTypeRepository(db *sql.DB) OperationTypeRepository {
	return &operationTypeRepository{db}
}

func (r *operationTypeRepository) GetOpTypeByID(ctx context.Context, opTypeID int) (*models.OperationType, error) {
	const query = `SELECT * FROM operation_types WHERE id = $1`

	opType := &models.OperationType{}
	err := r.db.QueryRowContext(ctx, query, opTypeID).Scan(&opType.ID, &opType.Description, &opType.OpType, &opType.ActiveSupport, &opType.CreatedAt, &opType.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return opType, nil
}
