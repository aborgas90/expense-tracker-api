package goalsdepo

import "time"

type RequestGoalsDepo struct {
	GoalID uint    `json:"goal_id" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
	Note   string  `json:"note" binding:"required"`
}

type ResponseGoalsDepo struct {
	ID        uint      `json:"id"`
	GoalID    uint      `json:"goal_id"`
	Amount    float64   `json:"amount"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
}
