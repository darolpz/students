package model

// Student model info
// @Description student information
// @Description with student_id,first name, last name, age and email
type Student struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
	Email     string `json:"email"`
}
