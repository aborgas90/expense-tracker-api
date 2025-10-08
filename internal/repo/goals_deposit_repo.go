package repo

import (
	"github.com/aborgas90/expense-tracker-api/internal/model"
	"gorm.io/gorm"
)

type GoalsDepositRepo struct {
	db *gorm.DB
}

func NewGoalsDepositRepo(db *gorm.DB) *GoalsDepositRepo {
	return &GoalsDepositRepo{db: db}
}

func (gd *GoalsDepositRepo) CreateDeposit(deposit *model.GoalDeposit) error {
	return gd.db.Create(deposit).Error
}

func (gd *GoalsDepositRepo) AutomaticUpdateDepoInsert(amount float64, id uint) error {
	return gd.db.Model(&model.Goal{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"current_amount": gorm.Expr("current_amount + ?", amount),
		}).Error
}

func (gd *GoalsDepositRepo) GetGoalDeposite(id uint) error {
	return gd.db.Preload("goals").First("id", id).Error
}

func (gd *GoalsDepositRepo) GetGoalsDepoIDandUserId(id uint, userId uint) (*model.GoalDeposit, error) {
	var deposit model.GoalDeposit
	err := gd.db.
		Joins("JOIN goals ON goals.id = goal_deposits.goal_id").
		Where("goal_deposits.id = ? AND goals.user_id = ?", id, userId).
		Preload("Goal").
		First(&deposit).Error

	if err != nil {
		return nil, err
	}

	return &deposit, nil
}
