package service

import (
	dto "github.com/aborgas90/expense-tracker-api/internal/dto/goals_depo"
	"github.com/aborgas90/expense-tracker-api/internal/model"
	"github.com/aborgas90/expense-tracker-api/internal/repo"
)

type GoalsDepoService struct {
	repo *repo.GoalsDepositRepo
}

func NewGoalsDepoService(r *repo.GoalsDepositRepo) *GoalsDepoService {
	return &GoalsDepoService{repo: r}
}

func (s *GoalsDepoService) CreateDepoServ(r *dto.RequestGoalsDepo) (*model.GoalDeposit, error) {
	deposit := &model.GoalDeposit{
		GoalID: r.GoalID,
		Amount: r.Amount,
		Note:   r.Note,
	}

	if err := s.repo.CreateDeposit(deposit); err != nil {
		return nil, err
	}

	// 2️⃣ Update goal progress
	err := s.repo.AutomaticUpdateDepoInsert(r.Amount, r.GoalID)

	if err != nil {
		return nil, err
	}

	return deposit, nil
}

func (s *GoalsDepoService) GetDepoByID(userID, depoID uint) (*model.GoalDeposit, error) {
	deposit, err := s.repo.GetGoalsDepoIDandUserId(depoID, userID)
	if err != nil {
		return nil, err
	}
	return deposit, nil
}
