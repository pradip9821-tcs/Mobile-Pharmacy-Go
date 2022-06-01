package models

type User struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Role        int    `json:"role"`
	Gender      int    `json:"gender"`
	Picture     string `json:"picture"`
	CountryCode string `json:"country_code"`
	Phone       int    `json:"phone"`
	IsTest      int    `json:"is_test"`
	IsActive    int    `json:"is_active"`
	IsVerify    int    `json:"is_verify"`
	IsDelete    int    `json:"is_delete"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
