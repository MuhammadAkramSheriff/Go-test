package models

type UserResponse struct {
	ID   string `json:"id" example:"123"`
	Name string `json:"name" example:"John Doe"`
}
