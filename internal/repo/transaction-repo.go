package repo

import (
	"github.com/aborgas90/expense-tracker-api/internal/model"
	"gorm.io/gorm"

)

type TransactionRepo struct {
	db *gorm.DB
}


func NewTransactionRepo(db *gorm.DB) *TransactionRepo{
	return &TransactionRepo{db: db}
}


func (t *TransactionRepo) GetAllTransactionByUser(userId uint) ([]model.Transaction, error) {
    var txs []model.Transaction
    if err := t.db.
        Preload("Category").
        Where("user_id = ?", userId).
        Order("occurred_at DESC").
        Find(&txs).Error; err != nil {
        return nil, err
    }
    return txs, nil
}

func (t *TransactionRepo) CreateTransactionUser(T *model.Transaction) error {
	return t.db.Create(T).Error
}