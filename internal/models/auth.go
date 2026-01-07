package models

// SignupRequest represents the request body for user signup
type SignupRequest struct {
	Name     string `json:"name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Dob      string `json:"dob" validate:"required,datetime=2006-01-02"`
}

// SignupResponse represents the response after successful signup
type SignupResponse struct {
	ID        int32  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents the response after successful login
type LoginResponse struct {
	Message string `json:"message"`
	User    struct {
		ID    int32  `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	} `json:"user"`
}

// AuthUser represents the authenticated user in the request context
type AuthUser struct {
	ID   int32  `json:"id"`
	Role string `json:"role"`
}
