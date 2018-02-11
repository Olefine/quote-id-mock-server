package domain

const successMessage = "Counts was sent."

// CountsResponse represents response date will be sent in case of success/error
type CountsResponse struct {
	Result  string `json:"result"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	JSON    string `json:"json,omitempty"`
}

// NewCountsResponseSuccess returns success response with
// Result = "Success"
// Success = true
// Message = const successMessage
// JSON = http.Request.Body(unmodified)
func NewCountsResponseSuccess(payload string) *CountsResponse {
	return &CountsResponse{
		Result: "Success", Success: true,
		Message: successMessage,
		JSON:    payload,
	}
}

// NewCountsResponseError returns error response with in case of any error
// Message will clarify whats was wrong
// JSON field will not be present in response
func NewCountsResponseError(errorMessage string) *CountsResponse {
	return &CountsResponse{
		Result: "Failure", Success: false,
		Message: errorMessage,
	}
}
