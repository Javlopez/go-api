package models

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid input data"`
}

// ValidationError represents an individual validation error
type ValidationError struct {
	Field   string `json:"field" example:"price"`
	Message string `json:"message" example:"must be greater than 0"`
}

// ValidationErrorResponse represents a response with validation errors
type ValidationErrorResponse struct {
	Errors []ValidationError `json:"errors"`
}
