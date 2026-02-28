package models

type User struct {
	ID         int    `json:"id" db:"id"`
	Username   string `json:"username" db:"username"`
	Email      string `json:"email" db:"email"`
	Password   string `json:"password" db:"password"`
	FirstName  string `json:"first_name" db:"first_name"`
	MiddleName string `json:"middle_name" db:"middle_name"`
	LastName   string `json:"last_name" db:"last_name"`
	IsActive   bool   `json:"is_active" db:"is_active"`
	IsVerified bool   `json:"is_verified" db:"is_verified"`
}
