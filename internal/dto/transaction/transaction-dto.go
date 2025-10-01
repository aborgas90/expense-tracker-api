package transaction

import "time"

// package transaction
type RequestTransaction struct {
	CategoryId uint    `json:"categoryId" binding:"required"`
	OccuredAt  string  `json:"occuredAt"  binding:"required"`
	Note       string  `json:"note"       binding:"required"`
	Amount     float64 `json:"amount"     binding:"required"`
	Currency   string  `json:"currency"   binding:"omitempty,len=3"`
}

type ResponseTransaction struct {
	ID           uint    `json:"id"`
	UserId       uint    `json:"userId"`
	CategoryId   uint    `json:"categoryId"`
	OccuredAt    string  `json:"occuredAt"`
	Note         string  `json:"note"`
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
	CreatedAt    string  `json:"createdAt"`
	CategoryName string  `json:"categoryName,omitempty"` // opsional, kalau preload Category
}

type UpdateTransactionRequest struct {
	CategoryId uint    `json:"categoryId" binding:"required"`
	OccuredAt  string  `json:"occuredAt"  binding:"required"`
	Note       string  `json:"note"       binding:"required"`
	Amount     float64 `json:"amount"     binding:"required"`
	Currency   string  `json:"currency"   binding:"omitempty,len=3"`
}

type LastTransactionDTO struct {
	Id       uint      `gorm:"column:id" json:"id"`
	Date     time.Time `gorm:"column:occurred_at" json:"date"`
	Category string    `gorm:"column:category" json:"category"`
	Amount   float64   `gorm:"column:amount" json:"amount"`
	Currency string    `gorm:"column:currency" json:"currency"`
	Note     string    `gorm:"column:note" json:"note"`
}
