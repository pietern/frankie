package errors

import (
	"fmt"
	"strings"
)

// errorHints maps error messages or types to helpful hints
var errorHints = map[string]string{
	"authentication required":    "Run 'frankie login' to authenticate",
	"not logged in":              "Run 'frankie login' to authenticate",
	"invalid email or password":  "Check your credentials and try again",
	"smart trading not enabled":  "Enable in Frank Energie app",
	"connection refused":         "Check your internet connection",
	"no such host":               "Check your internet connection",
	"timeout":                    "Check your internet connection",
	"network is unreachable":     "Check your internet connection",
	"server error":               "Frank Energie API may be temporarily unavailable",
	"500":                        "Frank Energie API may be temporarily unavailable",
	"502":                        "Frank Energie API may be temporarily unavailable",
	"503":                        "Frank Energie API may be temporarily unavailable",
}

// Format formats an error with a helpful hint if available
func Format(err error) string {
	if err == nil {
		return ""
	}

	msg := err.Error()
	hint := getHint(msg)

	if hint != "" {
		return fmt.Sprintf("%s\n\nHint: %s", msg, hint)
	}

	return msg
}

// getHint finds a hint for the given error message
func getHint(msg string) string {
	lowerMsg := strings.ToLower(msg)

	for pattern, hint := range errorHints {
		if strings.Contains(lowerMsg, strings.ToLower(pattern)) {
			return hint
		}
	}

	return ""
}
