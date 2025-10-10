package repo

import (
	"errors"
	"fmt"

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

// updating for GoalsDeposit
func (gd *GoalsDepositRepo) AutomaticUpdateDepoInsert(amount float64, id uint) error {
	return gd.db.Model(&model.Goal{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"current_amount": gorm.Expr("current_amount + ?", amount),
		}).Error
}

func (gd *GoalsDepositRepo) UpdateGoalsDepo(id uint, goal_id uint, userId uint, newAmount uint, note string) (*model.GoalDeposit, error) {
	var depo model.GoalDeposit

	if err := gd.db.Where("id = ?", id).First(&depo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("deposit not found")
		}
		return nil, err
	}

	oldAmount := depo.Amount
	oldGoalID := depo.GoalID

	if err := gd.db.Model(&depo).Updates(map[string]interface{}{
		"goal_id": goal_id,
		"amount":  newAmount,
		"note":    note,
	}).Error; err != nil {
		return nil, err
	}

	var delta int64 = int64(newAmount) - int64(oldAmount)

	if oldGoalID != 0 && goal_id != oldGoalID {
		if err := gd.updateGoalCurrentAmount(oldGoalID, -int64(oldAmount)); err != nil {
			return nil, err
		}
		if err := gd.updateGoalCurrentAmount(goal_id, int64(newAmount)); err != nil {
			return nil, err
		}
	} else {
		if err := gd.updateGoalCurrentAmount(goal_id, delta); err != nil {
			return nil, err
		}
	}
	return &depo, nil
}

func (gd *GoalsDepositRepo) updateGoalCurrentAmount(goalID uint, delta int64) error {
	return gd.db.Model(&model.Goal{}).
		Where("id = ?", goalID).
		UpdateColumn("current_amount", gorm.Expr("current_amount + ?", delta)).
		Error
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

func (gd *GoalsDepositRepo) DeleteGoalsDepo(id uint, userId uint) (*model.GoalDeposit, error) {
	var depo model.GoalDeposit

	tx := gd.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// 1. Get the deposit record first
	if err := tx.Where("id = ? ", id).First(&depo).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to find deposit: %w", err)
	}

	// 2. Delete the deposit
	if err := tx.Delete(&depo).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to delete deposit: %w", err)
	}

	// 3. Decrease the goal's current_amount & recalculate percentage
	err := tx.Model(&model.Goal{}).Where("user_id = ? AND id = ?", userId, depo.GoalID).
		UpdateColumn("current_amount", gorm.Expr("GREATEST(current_amount - ?, 0)", depo.Amount)).
		Error
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update goal current amount: %w", err)
	}

	// 4. Commit jika semuanya berhasil
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return &depo, nil
}
