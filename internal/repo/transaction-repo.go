package repo

import (
	"time"

	"github.com/aborgas90/expense-tracker-api/internal/model"
	dto "github.com/aborgas90/expense-tracker-api/internal/dto/transaction"
	
	"gorm.io/gorm"
)

type TransactionRepo struct {
	db *gorm.DB
}

type SummaryTransaction struct {
	Year    float64
	Month   float64
	Income  float64
	Expense float64
	Balance float64
	Status  string
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

func (r *TransactionRepo) SummaryTransaction(userId uint, month int, year int) (*SummaryTransaction, error) {
	var summary SummaryTransaction

	err := r.db.
		Table("transactions t").
		Select(`
			COALESCE(SUM(CASE WHEN c.transaction_type_id = 1 THEN t.amount END), 0) AS income,
			COALESCE(SUM(CASE WHEN c.transaction_type_id = 2 THEN t.amount END), 0) AS expense,
			COALESCE(SUM(CASE WHEN c.transaction_type_id = 1 THEN t.amount END), 0) -
			COALESCE(SUM(CASE WHEN c.transaction_type_id = 2 THEN t.amount END), 0) AS balance
		`).
		Joins("JOIN categories c ON t.category_id = c.id").
		Where("t.user_id = ? AND MONTH(t.occurred_at) = ? AND YEAR(t.occurred_at) = ?", userId, month, year).
		Scan(&summary).Error

	if err != nil {
		return nil, err
	}
	return &summary, nil
}

func (r *TransactionRepo) CheckSurplusTransaction(userId uint) ([]SummaryTransaction, error) {
	var summary []SummaryTransaction

	err := r.db.Table("transactions t").Select(`YEAR(t.occurred_at) AS year,
    MONTH(t.occurred_at) AS month,
    COALESCE(SUM(CASE WHEN c.transaction_type_id = 1 THEN t.amount END), 0) AS income,
    COALESCE(SUM(CASE WHEN c.transaction_type_id = 2 THEN t.amount END), 0) AS expense,
    COALESCE(SUM(CASE WHEN c.transaction_type_id = 1 THEN t.amount END), 0) -
    COALESCE(SUM(CASE WHEN c.transaction_type_id = 2 THEN t.amount END), 0) AS balance,
    CASE 
        WHEN COALESCE(SUM(CASE WHEN c.transaction_type_id = 1 THEN t.amount END), 0) -
             COALESCE(SUM(CASE WHEN c.transaction_type_id = 2 THEN t.amount END), 0) >= 0 
        THEN 'Surplus'
        ELSE 'Deficit'
    END AS Status`).Joins("JOIN categories c ON t.category_id = c.id").Where("t.user_id = ?", userId).Group("YEAR(t.occurred_at), MONTH(t.occurred_at)").Order("year, month").Scan(&summary).Error

	if err != nil {
		return nil, err
	}
	return summary, nil
}


func (r *TransactionRepo) Last7Transaction(userId uint) ([]dto.LastTransactionDTO, error) {
	var res []dto.LastTransactionDTO

	err := r.db.Table("transactions t").
		Select("t.id as id, t.amount, t.currency, t.note,t.occurred_at, c.type as category, t.amount").
		Joins("JOIN categories c ON t.category_id = c.id").
		Where("t.user_id = ?", userId).
		Order("t.occurred_at DESC").
		Limit(7).
		Scan(&res).Error

	if err != nil {
		return nil, err
	}

	return res, nil
}
