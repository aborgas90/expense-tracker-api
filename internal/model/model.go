package model

import "time"

type User struct {
	Id         uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Username   string     `gorm:"unique;not null" json:"username"`
	Password   string     `gorm:"not null" `
	Email      string     `gorm:"unique;not null" json:"email"`
	FirstName  string     `gorm:"not null" json:"first_name"`
	LastName   string     `gorm:"not null" json:"last_name"`
	Categories []Category `gorm:"foreignKey:UserID" json:"categories"`
}

type Category struct {
	ID                uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID            uint            `json:"user_id"`
	User              User            `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Type              string          `gorm:"not null" json:"type"`
	CreatedAt         time.Time       `gorm:"autoCreateTime" json:"created_at"`
	TransactionTypeID uint            `json:"transaction_type_id"`
	TransactionType   TransactionType `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"transaction_type"`
}

type Transaction struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     uint      `gorm:"not null;index" json:"user_id"`
	User       User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	CategoryID *uint     `json:"category_id"`
	Category   *Category `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"category"`
	Amount     float64   `gorm:"not null" json:"amount"`
	Currency   string    `gorm:"size:3;not null;default:IDR" json:"currency"`
	OccurredAt time.Time `gorm:"index;not null" json:"occurred_at"`
	Note       string    `json:"note"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type TransactionType struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"not null" json:"name"` // contoh: Income, Expense
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type Goal struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        uint      `gorm:"not null;index" json:"user_id"`
	User          User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Title         string    `gorm:"not null" json:"title"`
	TargetAmount  float64   `gorm:"not null" json:"target_amount"`
	CurrentAmount float64   `gorm:"not null;default:0" json:"current_amount"`
	Deadline      time.Time `gorm:"default:null" json:"deadline"`
	Status        string    `gorm:"default:'in-progress'" json:"status"` // in-progress, completed, canceled
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type GoalDeposit struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	GoalID    uint      `gorm:"not null;index" json:"goal_id"`
	Goal      Goal      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"goal"`
	Amount    float64   `gorm:"not null" json:"amount"`
	Note      string    `json:"note"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
