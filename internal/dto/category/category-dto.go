package category

type CategoryRequest struct {
	TypeCategory string `json:"type" binding:"required"`
}

type CategoryResponse struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
}

// type RegisterUserRequest struct {
// 	Username  string `json:"username" binding:"required"`
// 	Password  string `json:"password" binding:"required"`
// 	Email     string `json:"email" binding:"required,email"`
// 	FirstName string `json:"first_name" binding:"required"`
// 	LastName  string `json:"last_name" binding:"required"`
// }