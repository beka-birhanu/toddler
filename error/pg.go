package error

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/beka-birhanu/toddler/status"
	"github.com/lib/pq"
)

const (
	postgresErrUniqueViolation  = "23505"
	postgresErrForeignKey       = "23503"
	postgresErrNotNullViolation = "23502"
	postgresErrCheckViolation   = "23514"
)

// FromDBError maps database-level errors into structured application errors.
func FromDBError(err error, entityName string) *Error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return &Error{
			PublicStatusCode:  status.NotFoundResource,
			ServiceStatusCode: status.NotFoundResource,
			PublicMessage:     fmt.Sprintf("%s not found", entityName),
			PublicMetaData: map[string]string{
				"error_type":   "Data not found",
				"resourceName": entityName,
			},
			ServiceMessage: fmt.Sprintf("No record found for %s: %s", entityName, err),
			ServiceMetaData: map[string]string{
				"error_type":   "Data not found",
				"resourceName": entityName,
				"raw_error":    err.Error(),
			},
		}
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case postgresErrUniqueViolation:
			return &Error{
				PublicStatusCode:  status.ConflictDuplicateData,
				ServiceStatusCode: status.ConflictDuplicateData,
				PublicMessage:     fmt.Sprintf("A %s with the same value already exists", entityName),
				PublicMetaData: map[string]string{
					"error_type":   "Data duplication",
					"resourceName": entityName,
				},
				ServiceMessage: fmt.Sprintf("Unique constraint violation on %s: %s", entityName, pqErr.Message),
				ServiceMetaData: map[string]string{
					"pgcode":         string(pqErr.Code),
					"constraint":     pqErr.Constraint,
					"error_type":     "Data duplication",
					"resourceName":   entityName,
					"error_message":  pqErr.Message,
					"error_severity": pqErr.Severity,
					"raw_error":      pqErr.Error(),
				},
			}
		case postgresErrForeignKey:
			return &Error{
				PublicStatusCode:  status.BadRequest,
				ServiceStatusCode: status.BadRequest,
				PublicMessage:     fmt.Sprintf("%s has invalid reference to related data", entityName),
				PublicMetaData: map[string]string{
					"error_type":   "Foreign key violation",
					"resourceName": entityName,
				},
				ServiceMessage: fmt.Sprintf("Foreign key constraint failed on %s: %s", entityName, pqErr.Message),
				ServiceMetaData: map[string]string{
					"pgcode":         string(pqErr.Code),
					"constraint":     pqErr.Constraint,
					"error_type":     "Foreign key violation",
					"resourceName":   entityName,
					"error_message":  pqErr.Message,
					"error_severity": pqErr.Severity,
					"raw_error":      pqErr.Error(),
				},
			}
		case postgresErrNotNullViolation:
			return &Error{
				PublicStatusCode:  status.BadRequest,
				ServiceStatusCode: status.BadRequest,
				PublicMessage:     fmt.Sprintf("%s is missing required fields", entityName),
				PublicMetaData: map[string]string{
					"error_type":   "Missing field",
					"resourceName": entityName,
				},
				ServiceMessage: fmt.Sprintf("NOT NULL constraint failed on %s: %s", entityName, pqErr.Message),
				ServiceMetaData: map[string]string{
					"pgcode":         string(pqErr.Code),
					"column":         pqErr.Column,
					"error_type":     "Missing field",
					"resourceName":   entityName,
					"error_message":  pqErr.Message,
					"error_severity": pqErr.Severity,
					"raw_error":      pqErr.Error(),
				},
			}
		case postgresErrCheckViolation:
			return &Error{
				PublicStatusCode:  status.BadRequest,
				ServiceStatusCode: status.BadRequest,
				PublicMessage:     fmt.Sprintf("%s failed validation rules", entityName),
				PublicMetaData: map[string]string{
					"error_type":   "Constraint check failed",
					"resourceName": entityName,
				},
				ServiceMessage: fmt.Sprintf("CHECK constraint violation on %s: %s", entityName, pqErr.Message),
				ServiceMetaData: map[string]string{
					"pgcode":         string(pqErr.Code),
					"constraint":     pqErr.Constraint,
					"error_type":     "Constraint check failed",
					"resourceName":   entityName,
					"error_message":  pqErr.Message,
					"error_severity": pqErr.Severity,
					"raw_error":      pqErr.Error(),
				},
			}
		default:
			// Unhandled DB errors — treat as server errors
			return &Error{
				PublicStatusCode:  status.ServerError,
				ServiceStatusCode: status.ServerErrorDatabase,
				PublicMessage:     "A server error occurred. Please try again later.",
				PublicMetaData: map[string]string{
					"error_type":   "Internal database error",
					"resourceName": entityName,
				},
				ServiceMessage: fmt.Sprintf("Unhandled PostgreSQL error for %s: %s", entityName, pqErr.Message),
				ServiceMetaData: map[string]string{
					"pgcode":         string(pqErr.Code),
					"resourceName":   entityName,
					"error_message":  pqErr.Message,
					"error_severity": pqErr.Severity,
					"raw_error":      pqErr.Error(),
				},
			}
		}
	}

	// Fallback: truly unknown error — treat as internal error
	return &Error{
		PublicStatusCode:  status.ServerError,
		ServiceStatusCode: status.ServerErrorDatabase,
		PublicMessage:     "A server error occurred. Please try again later.",
		PublicMetaData: map[string]string{
			"error_type":   "Unknown server error",
			"resourceName": entityName,
		},
		ServiceMessage: fmt.Sprintf("Unexpected DB error for %s: %s", entityName, err),
		ServiceMetaData: map[string]string{
			"error_type":   "Unknown database error",
			"resourceName": entityName,
			"raw_error":    err.Error(),
		},
	}
}
