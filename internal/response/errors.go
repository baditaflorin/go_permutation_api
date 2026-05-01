package response

import "fmt"

// Error codes — stable strings clients can match against.
const (
	CodeTooManyElements  = "TOO_MANY_ELEMENTS"
	CodeEmptyElement     = "EMPTY_ELEMENT"
	CodeMissingParameter = "MISSING_PARAMETER"
	CodeInvalidJSON      = "INVALID_JSON"
	CodeMethodNotAllowed = "METHOD_NOT_ALLOWED"
	CodeRateLimitExceeded = "RATE_LIMIT_EXCEEDED"
	CodeRequestTimeout   = "REQUEST_TIMEOUT"
	CodeInternalError    = "INTERNAL_ERROR"
)

const docsBase = "https://github.com/baditaflorin/go_permutation_api/blob/main/docs/API.md"

// APIError represents a machine-readable, human-friendly error.
type APIError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Suggestion string `json:"suggestion,omitempty"`
	DocsURL    string `json:"docs_url,omitempty"`
}

func (e *APIError) Error() string { return fmt.Sprintf("%s: %s", e.Code, e.Message) }

func ErrTooManyElements(got, max int) *APIError {
	return &APIError{
		Code:       CodeTooManyElements,
		Message:    fmt.Sprintf("You provided %d elements, but the maximum allowed is %d.", got, max),
		Suggestion: fmt.Sprintf("Reduce your input to %d or fewer elements, or increase MAX_ELEMENTS in your server configuration.", max),
		DocsURL:    docsBase + "#error-handling",
	}
}

func ErrEmptyElement(index int) *APIError {
	return &APIError{
		Code:       CodeEmptyElement,
		Message:    fmt.Sprintf("Element at index %d is empty or contains only whitespace.", index),
		Suggestion: "Remove blank elements from your input array.",
		DocsURL:    docsBase + "#error-handling",
	}
}

func ErrMissingParameter(param string) *APIError {
	return &APIError{
		Code:       CodeMissingParameter,
		Message:    fmt.Sprintf("Required query parameter '%s' is missing.", param),
		Suggestion: fmt.Sprintf("Add ?%s=a,b,c to your request URL.", param),
		DocsURL:    docsBase + "#get-generate-permutations-get",
	}
}

func ErrInvalidJSON(detail string) *APIError {
	return &APIError{
		Code:       CodeInvalidJSON,
		Message:    "The request body contains invalid JSON: " + detail,
		Suggestion: `Ensure the body is a valid JSON array, e.g. ["a","b","c"].`,
		DocsURL:    docsBase + "#post-generate-permutations-post",
	}
}

func ErrMethodNotAllowed(method string) *APIError {
	return &APIError{
		Code:       CodeMethodNotAllowed,
		Message:    fmt.Sprintf("HTTP method '%s' is not supported on this endpoint.", method),
		Suggestion: "Use GET or POST.",
		DocsURL:    docsBase,
	}
}

func ErrRateLimitExceeded() *APIError {
	return &APIError{
		Code:       CodeRateLimitExceeded,
		Message:    "You have exceeded the request rate limit.",
		Suggestion: "Wait a moment and retry. Consider batching requests or contacting us for a higher limit.",
		DocsURL:    docsBase + "#rate-limiting",
	}
}

func ErrRequestTimeout() *APIError {
	return &APIError{
		Code:       CodeRequestTimeout,
		Message:    "The request took too long to process.",
		Suggestion: "Try fewer elements. Requests with more than 10 elements may exceed the timeout.",
		DocsURL:    docsBase + "#performance-considerations",
	}
}

func ErrInternal() *APIError {
	return &APIError{
		Code:       CodeInternalError,
		Message:    "An unexpected error occurred on the server.",
		Suggestion: "Please try again. If the problem persists, open an issue on GitHub.",
		DocsURL:    "https://github.com/baditaflorin/go_permutation_api/issues",
	}
}
