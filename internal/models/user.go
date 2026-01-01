package models

type UserRequest struct {
	Name string `json:"name" validate:"required,min=2"`
	Dob  string `json:"dob" validate:"required,datetime=2006-01-02"`
}
