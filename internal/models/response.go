package models

type UserResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Dob  string `json:"dob"`
}

type UserWithAgeResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Dob  string `json:"dob"`
	Age  int    `json:"age"`
}

// ErrorDetail contains detailed error information
type ErrorDetail struct {
	Message   string `json:"message"`
	Code      string `json:"code"`
	RequestID string `json:"request_id,omitempty"`
}

// ErrorResponse is the standardized error response format
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// Legacy simple error response (deprecated, for backward compatibility)
type SimpleErrorResponse struct {
	Error string `json:"error"`
}

type PaginationMeta struct {
	Total       int64 `json:"total"`
	Page        int   `json:"page"`
	Limit       int   `json:"limit"`
	TotalPages  int   `json:"total_pages"`
	HasNext     bool  `json:"has_next"`
	HasPrevious bool  `json:"has_previous"`
}

type PaginatedUsersResponse struct {
	Data       []UserWithAgeResponse `json:"data"`
	Pagination PaginationMeta        `json:"pagination"`
}
