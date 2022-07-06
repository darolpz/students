package model

// Authentication model info
// @Description Authentication information
// @Description with email and password
type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
