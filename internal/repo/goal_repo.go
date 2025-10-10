package repo

import (
	"errors"
	"time"

	"github.com/aborgas90/expense-tracker-api/internal/model"
	"gorm.io/gorm"
)

type GoalRepo struct {
	db *gorm.DB
}

func NewGoalRepo(db *gorm.DB) *GoalRepo {
	return &GoalRepo{db: db}
}

func (g *GoalRepo) CreateGoalsTarget(c *model.Goal) error {
	return g.db.Create(c).Error
}

func (g *GoalRepo) GetGoalDataByIdUser(userId uint) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := g.db.Table("goals").
		Select("id, title, target_amount, current_amount, (current_amount / target_amount * 100) AS percentage, status, deadline, created_at, updated_at").
		Where("user_id = ?", userId).
		Find(&results).Error

	if err != nil {
		return nil, err
	}
	return results, nil
}

func (g *GoalRepo) UpdateGoalsDataById(goal *model.Goal) error {
	updateData := map[string]interface{}{
		"title":          goal.Title,
		"target_amount":  goal.TargetAmount,
		"current_amount": goal.CurrentAmount,
		"status":         goal.Status,
		"updated_at":     time.Now(),
	}

	if !goal.Deadline.IsZero() {
		updateData["deadline"] = goal.Deadline
	}

	if err := g.db.Model(&model.Goal{}).
		Where("id = ? AND user_id = ?", goal.ID, goal.UserID).
		Updates(updateData).Error; err != nil {
		return err
	}

	return nil
}

func (g *GoalRepo) DeleteGoalDataById(userId uint, id uint) error {
	if err := g.db.Where("id = ? AND user_id = ?", id, userId).Delete(&model.Goal{}).Error; err != nil {
		return err
	}
	return nil
}

func (g *GoalRepo) GetGoalById(userId uint, id uint) (*model.Goal, error) {
	var goal model.Goal
	if err := g.db.Where("id = ? AND user_id = ?", id, userId).First(&goal).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &goal, nil
}
