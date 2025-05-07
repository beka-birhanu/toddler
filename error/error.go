// Package error provides custom error handling with detailed status codes,
// public messages, and metadata. The errors are designed to give clear
// feedback to both the user and the service, helping with precise error
// classification and debugging.
//
// It uses custom application-specific status codes that extend standard
// HTTP status semantics, providing more granular codes for clearer error
// handling, especially when integrating with an API or other services.
package error

import (
	"fmt"

	"github.com/beka-birhanu/toddler/status"
)

type ErrorTypes string

type Error struct {
	PublicStatusCode  status.StatusCode
	ServiceStatusCode status.StatusCode
	PublicMessage     string
	ServiceMessage    string
	PublicMetaData    map[string]string
	ServiceMetaData   map[string]string
}

// Error implements the error interface.
// Return a formatted string with status code, messages, and metadata
func (e *Error) Error() string {
	return fmt.Sprintf(
		"{publicStatus: %s (%d), serviceStatus: %s (%d), publicMessage: '%s', serviceMessage: '%s', publicMetaData: %s, serviceMetaData: %s}",
		status.GetErrorName(e.PublicStatusCode),
		e.PublicStatusCode,
		status.GetErrorName(e.ServiceStatusCode),
		e.ServiceStatusCode,
		e.PublicMessage,
		e.ServiceMessage,
		formatMetaData(e.PublicMetaData),
		formatMetaData(e.ServiceMetaData),
	)
}

// Helper function to format metadata as a string
func formatMetaData(metaData map[string]string) string {
	if len(metaData) == 0 {
		return "{}"
	}

	formatted := "{"
	for key, value := range metaData {
		formatted += key + ": '" + value + "', "
	}
	formatted = formatted[:len(formatted)-2] + "}"
	return formatted
}

func (e *Error) NeutralizeOverDetailedStatus() {
	e.PublicStatusCode = status.SuppressOverDetail(e.PublicStatusCode)
}
