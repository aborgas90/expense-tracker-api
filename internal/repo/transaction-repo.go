package repo

import (
	"time"

	"github.com/aborgas90/expense-tracker-api/internal/model"
	"gorm.io/gorm"
)

type TransactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) *TransactionRepo {
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

func (t *TransactionRepo) GetTransactionById(userId uint, id uint) (*model.Transaction, error) {
	var txs model.Transaction
	if err := t.db.Preload("Category").Where("user_id = ? AND id = ?", userId, id).First(&txs).Error; err != nil {
		return nil, err
	}
	return &txs, nil
}

func (t *TransactionRepo) UpdateTransaction(userId uint, id uint, category_id uint, amount float64, currency string, occured_at time.Time, note string) (*model.Transaction, error) {
	var txs model.Transaction
	if err := t.db.Model(&model.Transaction{}).Where("id = ? AND user_id = ?", id, userId).Updates(map[string]interface{}{"category_id": category_id, "amount": amount, "currency": currency, "occurred_at": occured_at, "note": note}).Error; err != nil {
		return nil, err
	}

	if err := t.db.Preload("Category").Preload("User").
		First(&txs, "id = ? AND user_id = ?", id, userId).Error; err != nil {
		return nil, err
	}

	return &txs, nil
}

func (r *TransactionRepo) DeleteTransaction(userId uint, id uint) (int64, error) {
	res := r.db.Where("id = ? AND user_id = ?", id, userId).Delete(&model.Transaction{})
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}


func (r *TransactionRepo) FindByIDAndUser(id uint, userID uint) (*model.Transaction, error) {
	var txs model.Transaction
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&txs).Error; err != nil {
		return nil, err
	}
	return &txs, nil
}

func (t *TransactionRepo) CreateTransactionUser(T *model.Transaction) error {
	return t.db.Create(T).Error
}
