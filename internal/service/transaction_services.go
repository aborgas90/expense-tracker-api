package service

import (
	"fmt"
	"strings"
	"time"

	dto "github.com/aborgas90/expense-tracker-api/internal/dto/transaction"
	"github.com/aborgas90/expense-tracker-api/internal/model"
	"github.com/aborgas90/expense-tracker-api/internal/repo"
)

type TransactionService struct {
	repo *repo.TransactionRepo
}

func NewTransactionService(r *repo.TransactionRepo) *TransactionService {
	return &TransactionService{repo: r}
}

func (s *TransactionService) GetTransactionByUser(userId uint) ([]dto.ResponseTransaction, error) {
	txs, err := s.repo.GetAllTransactionByUser(userId)
	if err != nil {
		return nil, err
	}

	res := make([]dto.ResponseTransaction, 0, len(txs))
	for _, t := range txs {
		var catID uint
		var catName string
		if t.CategoryID != nil {
			catID = *t.CategoryID
		}
		if t.Category != nil {
			catName = t.Category.Type
		}

		res = append(res, dto.ResponseTransaction{
			ID:           t.ID,
			UserId:       t.UserID,
			CategoryId:   catID,
			CategoryName: catName,
			Amount:       t.Amount,
			Currency:     t.Currency,
			OccuredAt:    t.OccurredAt.Format(time.RFC3339),
			Note:         t.Note,
			CreatedAt:    t.CreatedAt.Format(time.RFC3339),
		})
	}
	return res, nil
}

func (s *TransactionService) CreateTransactionUser(userId uint, req *dto.RequestTransaction) (*dto.ResponseTransaction, error) {
	parsedTime, err := time.Parse(time.RFC3339, req.OccuredAt)
	if err != nil {
		return nil, fmt.Errorf("invalid occuredAt (expect RFC3339): %w", err)
	}

	currency := strings.ToUpper(strings.TrimSpace(req.Currency))
	if currency == "" {
		currency = "IDR"
	}

	newTransaction := &model.Transaction{
		UserID:     userId,
		CategoryID: UintPtr(req.CategoryId),
		Note:       req.Note,
		Amount:     req.Amount,
		Currency:   req.Currency,
		OccurredAt: parsedTime,
	}

	if err := s.repo.CreateTransactionUser(newTransaction); err != nil {
		return nil, err
	}

	return &dto.ResponseTransaction{
		ID:         newTransaction.ID,
		UserId:     newTransaction.UserID,
		CategoryId: *UintPtr(*newTransaction.CategoryID),
		OccuredAt:  newTransaction.OccurredAt.Format(time.RFC3339),
		Note:       newTransaction.Note,
		Amount:     newTransaction.Amount,
		Currency:   newTransaction.Currency,
		CreatedAt:  newTransaction.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *TransactionService) GetTransactionById(userId uint, id uint) (*dto.ResponseTransaction, error) {
	txs, err := s.repo.GetTransactionById(userId, id) // now return model.Transaction, not []model.Transaction
	if err != nil {
		return nil, err
	}

	var catID uint
	var catName string
	if txs.CategoryID != nil {
		catID = *txs.CategoryID
	}
	if txs.Category != nil {
		catName = txs.Category.Type
	}

	return &dto.ResponseTransaction{
		ID:           txs.ID,
		UserId:       txs.UserID,
		CategoryId:   catID,
		CategoryName: catName,
		Amount:       txs.Amount,
		Currency:     txs.Currency,
		OccuredAt:    txs.OccurredAt.Format(time.RFC3339),
		Note:         txs.Note,
		CreatedAt:    txs.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *TransactionService) UpdateTransactionUser(
	userId uint,
	id uint,
	categoryId uint,
	amount float64,
	currency string,
	occuredAt string,
	note string,
) (*dto.ResponseTransaction, error) {

	parsedTime, err := time.Parse(time.RFC3339, occuredAt)
	if err != nil {
		return nil, fmt.Errorf("invalid occuredAt (expect RFC3339): %w", err)
	}

	txs, err := s.repo.UpdateTransaction(
		userId,
		id,
		categoryId,
		amount,
		currency,
		parsedTime,
		note,
	)
	if err != nil {
		return nil, err
	}

	var catID uint
	var catName string
	if txs.CategoryID != nil {
		catID = *txs.CategoryID
	}
	if txs.Category != nil {
		catName = txs.Category.Type
	}

	return &dto.ResponseTransaction{
		ID:           txs.ID,
		UserId:       txs.UserID,
		CategoryId:   catID,
		CategoryName: catName,
		Amount:       txs.Amount,
		Currency:     txs.Currency,
		OccuredAt:    txs.OccurredAt.Format(time.RFC3339),
		Note:         txs.Note,
		CreatedAt:    txs.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *TransactionService) DeleteTransaction(userId uint, id uint) (int64, error) {
	rows, err := s.repo.DeleteTransaction(userId, id)
	if err != nil {
		return 0, err
	}
	return rows, nil
}

func UintPtr(v uint) *uint {
	return &v
}

func (s *TransactionService) SummaryTransaction(userId uint, month int, year int) (*dto.SummaryTransaction, error) {
	res, err := s.repo.SummaryTransaction(userId, month, year)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return &dto.SummaryTransaction{Income: 0, Expense: 0}, nil
	}

	return &dto.SummaryTransaction{
		Income:  res.Income,
		Expense: res.Expense,
		Balance: res.Balance,
	}, nil
}

func (s *TransactionService) CheckSurplusDeficitTransaction(userId uint) ([]dto.SummaryTransaction, error) {
	res, err := s.repo.CheckSurplusTransaction(userId)
	if err != nil {
		return nil, err
	}

	var summaries []dto.SummaryTransaction
	for _, r := range res {
		summaries = append(summaries, dto.SummaryTransaction{
			Year:    r.Year,
			Month:   r.Month,
			Income:  r.Income,
			Expense: r.Expense,
			Balance: r.Balance,
			Status:  r.Status,
		})
	}

	return summaries, nil
}

func (s *TransactionService) Last7Transaction(userId uint) ([]dto.LastTransactionDTO, error) {
	return s.repo.Last7Transaction(userId)
}
