package goals

type RequestGoals struct {
	Title          string  `json:"title" binding:"required"`
	Target_amount  float64 `json:"target_amount"  binding:"required"`
	Current_amount float64 `json:"current_amount"       binding:"required"`
	Deadline       string  `json:"deadline"`
	Status         string  `json:"status"`
}

type ResponseGoals struct {
	Id             uint    `json:"id"`
	Title          string  `json:"title"`
	Target_amount  float64 `json:"target_amount"`
	Current_amount float64 `json:"current_amount"`
	Deadline       string  `json:"deadline"`
	Created_at     string  `json:"created_at"`
}
