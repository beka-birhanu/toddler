// Package status defines custom application-specific status codes.
//
// These codes extend standard HTTP status semantics with more granular,
// 4-digit codes for clearer error handling. Each group of codes is
// aligned with HTTP categories (e.g., 4000s for Bad Request, 5000s for Server Error).
//
// Example categories:
//   - 4000–4009: Bad Request (client input errors)
//   - 4010–4019: Unauthorized (auth failures)
//   - 4030–4039: Forbidden (access control)
//   - 4040–4049: Not Found (missing resources)
//   - 5000–5009: Server Errors (internal failures)
//
// Each status code has a short constant name for code clarity and
// can be mapped to user-friendly messages or used in API responses.
package status

import "fmt"

// StatusCode defines custom application-specific status codes.
type StatusCode int

// BadRequest-related errors (4000 - 4009)
const (
	BadRequest                StatusCode = 4000 + iota // Generic bad request
	BadRequestMissingField                             // Required field missing
	BadRequestTypeMismatch                             // Type mismatch
	BadRequestFieldConstraint                          // Field constraint failed
	BadRequestInvalidFormat                            // Invalid format
	BadRequestOutOfRange                               // Value out of range
	BadRequestInvalidValue                             // Invalid value
	BadRequestEnumViolation                            // Enum value not allowed
)

// Unauthorized-related errors (4010 - 4019)
const (
	Unauthorized                  StatusCode = 4010 + iota // Generic unauthorized
	UnauthorizedInvalidCredential                          // Invalid credentials
	UnauthorizedTokenRequired                              // Token required
	UnauthorizedInvalidToken                               // Invalid token
)

// Forbidden-related errors (4030 - 4039)
const (
	Forbidden                   StatusCode = 4030 + iota // Generic forbidden
	ForbiddenNotEnoughPrivilege                          // Insufficient privileges
	ForbiddenOnlyOwners                                  // Allowed for resource owners only
)

// NotFound-related errors (4040 - 4049)
const (
	NotFound         StatusCode = 4040 + iota // Generic not found
	NotFoundResource                          // Resource not found
)

// Conflict-realted errors(4040 - 4049)
const (
	Conflict              StatusCode = 4090 + iota // Generic conflict
	ConflictDuplicateData                          // Conflict Duplicate Data
)

// Server-related errors (5000 - 5009)
const (
	ServerError                     StatusCode = 5000 + iota // Generic server error
	ServerErrorDatabase                                      // Database error
	ServerErrorServiceCommunication                          // Service communication failed
)

// A map to associate StatusCode with error names.
var statusCodeMap = map[StatusCode]string{
	BadRequest:                      "BadRequest",
	BadRequestMissingField:          "BadRequest_MissingField",
	BadRequestTypeMismatch:          "BadRequest_TypeMismatch",
	BadRequestFieldConstraint:       "BadRequest_FieldConstraint",
	BadRequestInvalidFormat:         "BadRequest_InvalidFormat",
	BadRequestOutOfRange:            "BadRequest_OutOfRange",
	BadRequestInvalidValue:          "BadRequest_InvalidValue",
	BadRequestEnumViolation:         "BadRequest_EnumViolation",
	Unauthorized:                    "Unauthorized",
	UnauthorizedInvalidCredential:   "Unauthorized_InvalidCredential",
	UnauthorizedTokenRequired:       "Unauthorized_TokenRequired",
	UnauthorizedInvalidToken:        "Unauthorized_InvalidToken",
	Forbidden:                       "Forbidden",
	ForbiddenNotEnoughPrivilege:     "Forbidden_NotEnoughPrivilege",
	ForbiddenOnlyOwners:             "Forbidden_OnlyOwners",
	NotFound:                        "NotFound",
	NotFoundResource:                "NotFound_Resource",
	Conflict:                        "Conflict",
	ConflictDuplicateData:           "Conflict_DuplicateData",
	ServerError:                     "ServerError",
	ServerErrorDatabase:             "ServerError_Database",
	ServerErrorServiceCommunication: "ServerError_ServiceCommunication",
}

// GetErrorName takes a StatusCode and returns the corresponding error name as a string.
func GetErrorName(code StatusCode) string {
	if name, exists := statusCodeMap[code]; exists {
		return name
	}
	return fmt.Sprintf("UnknownStatusCode-%d", code)
}

// suppressMap maps over-detailed status codes to generalized public-safe ones.
var suppressMap = map[StatusCode]StatusCode{
	BadRequestOutOfRange:            BadRequest,
	BadRequestInvalidValue:          BadRequest,
	BadRequestEnumViolation:         BadRequest,
	ForbiddenOnlyOwners:             Forbidden,
	ServerErrorDatabase:             ServerError,
	ServerErrorServiceCommunication: ServerError,
}

// SuppressOverDetail returns a neutralized version of the given StatusCode.
// If a mapping is not found, it returns the original code.
func SuppressOverDetail(code StatusCode) StatusCode {
	if suppressed, ok := suppressMap[code]; ok {
		return suppressed
	}
	return code
}
