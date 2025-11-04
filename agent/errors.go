package agent

import (
	"strings"
)

// isTimeoutError checks if an error is a timeout error by examining the error message
// This helps identify timeout errors from WaitForJobCompletion and other timeout scenarios
func isTimeoutError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	// Check for common timeout error patterns
	return strings.Contains(errStr, "timed out") ||
		strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "context deadline exceeded") ||
		strings.Contains(errStr, "deadline exceeded")
}

