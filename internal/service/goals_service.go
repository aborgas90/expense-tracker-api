package service

import (
	"fmt"
	"math"
	"time"

	dto "github.com/aborgas90/expense-tracker-api/internal/dto/goals"
	"github.com/aborgas90/expense-tracker-api/internal/model"
	"github.com/aborgas90/expense-tracker-api/internal/repo"
)

type GoalsServices struct {
	repo *repo.GoalRepo
}

func NewGoalsService(r *repo.GoalRepo) *GoalsServices {
	return &GoalsServices{repo: r}
}

func (g *GoalsServices) CreateGoalsTarget(userId uint, req *dto.RequestGoals) (*dto.ResponseGoals, error) {
	var parsedTime time.Time

	// Kalau dikirim deadline, parse RFC3339-nya
	if req.Deadline != "" {
		t, err := time.Parse(time.RFC3339, req.Deadline)
		if err != nil {
			return nil, fmt.Errorf("invalid deadline (expect RFC3339): %w", err)
		}
		parsedTime = t
	}

	newGoals := &model.Goal{
		UserID:        userId,
		Title:         req.Title,
		TargetAmount:  req.Target_amount,
		CurrentAmount: req.Current_amount,
		Deadline:      parsedTime, // bisa 0001-01-01 kalau kosong
	}

	if err := g.repo.CreateGoalsTarget(newGoals); err != nil {
		return nil, err
	}

	resp := &dto.ResponseGoals{
		Id:             newGoals.ID,
		Title:          newGoals.Title,
		Target_amount:  newGoals.TargetAmount,
		Current_amount: newGoals.CurrentAmount,
		Created_at:     newGoals.CreatedAt.Format(time.RFC3339),
	}

	if !newGoals.Deadline.IsZero() {
		resp.Deadline = newGoals.Deadline.Format(time.RFC3339)
	}

	return resp, nil
}

func (g *GoalsServices) GetGoalDataByIdUser(userId uint) ([]dto.ResponseGoals, error) {
	data, err := g.repo.GetGoalDataByIdUser(userId)
	if err != nil {
		return nil, err
	}

	var res []dto.ResponseGoals
	for _, goal := range data {
		// --- Handle nullable deadline
		var deadlineStr string
		if val, ok := goal["deadline"]; ok && val != nil {
			if t, ok := val.(time.Time); ok {
				deadlineStr = t.Format(time.RFC3339)
			}
		} else {
			deadlineStr = "" // atau "N/A"
		}

		// --- Handle created_at (optional juga biar gak panic)
		var createdAtStr string
		if val, ok := goal["created_at"]; ok && val != nil {
			if t, ok := val.(time.Time); ok {
				createdAtStr = t.Format(time.RFC3339)
			}
		}

		res = append(res, dto.ResponseGoals{
			Id:             uint(goal["id"].(uint64)),
			Title:          goal["title"].(string),
			Target_amount:  goal["target_amount"].(float64),
			Current_amount: goal["current_amount"].(float64),
			Deadline:       deadlineStr,
			Percentage:     math.Round(goal["percentage"].(float64)*100) / 100,
			Created_at:     createdAtStr,
		})
	}

	return res, nil
}

func (s *GoalsServices) UpdateGoalsDataById(userId uint, id uint, req *dto.RequestGoals) (*model.Goal, error) {
	// Ambil data goal lama dari database
	existingGoal, err := s.repo.GetGoalById(userId, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch goal: %w", err)
	}

	if existingGoal == nil {
		return nil, fmt.Errorf("goal not found")
	}

	// Update hanya field yang dikirim user
	if req.Title != "" {
		existingGoal.Title = req.Title
	}
	if req.Target_amount > 0 {
		existingGoal.TargetAmount = req.Target_amount
	}
	if req.Current_amount >= 0 {
		existingGoal.CurrentAmount = req.Current_amount
	}
	if req.Status != "" {
		existingGoal.Status = req.Status
	}
	if req.Deadline != "" {
		parsedDeadline, err := time.Parse(time.RFC3339, req.Deadline)
		if err != nil {
			return nil, fmt.Errorf("invalid deadline format")
		}
		existingGoal.Deadline = parsedDeadline
	}

	existingGoal.UpdatedAt = time.Now()

	// Kirim ke repo untuk disimpan
	if err := s.repo.UpdateGoalsDataById(existingGoal); err != nil {
		return nil, fmt.Errorf("failed to update goal: %w", err)
	}

	return existingGoal, nil
}

func (s *GoalsServices) DeleteGoalsDataById(userId uint, id uint) (*model.Goal, error) {
	goal, err := s.repo.GetGoalById(userId, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch goal: %w", err)
	}

	if goal == nil {
		return nil, fmt.Errorf("goal not found")
	}

	if err := s.repo.DeleteGoalDataById(userId, id); err != nil {
		return nil, fmt.Errorf("failed to delete goal: %w", err)
	}

	return goal, nil
}

func (s *GoalsServices) GetGoalsById(userId uint, id uint) (*dto.ResponseGoals, error) {
	goal, err := s.repo.GetGoalById(userId, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch goal: %w", err)
	}

	if goal == nil {
		return nil, fmt.Errorf("goal not found")
	}

	// mapping ke DTO
	res := &dto.ResponseGoals{
		Id:             goal.ID,
		Title:          goal.Title,
		Target_amount:  goal.TargetAmount,
		Current_amount: goal.CurrentAmount,
		Deadline:       goal.Deadline.Format(time.RFC3339),
		Created_at:     goal.CreatedAt.Format(time.RFC3339),
	}

	return res, nil
}
