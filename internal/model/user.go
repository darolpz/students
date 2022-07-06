package model

// User model info
// @Description user information
// @Description with user_id, name, email, password and role
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
