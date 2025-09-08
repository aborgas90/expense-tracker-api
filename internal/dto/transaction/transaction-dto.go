package transaction

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
