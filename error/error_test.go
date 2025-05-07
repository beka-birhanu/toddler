package error_test

import (
	"testing"

	"github.com/beka-birhanu/toddler/error"
	"github.com/beka-birhanu/toddler/status"
)

func TestError_Error(t *testing.T) {
	err := &error.Error{
		PublicStatusCode:  status.BadRequestMissingField,
		ServiceStatusCode: status.BadRequestMissingField,
		PublicMessage:     "Missing required field",
		ServiceMessage:    "Field 'username' is missing in the payload",
		PublicMetaData: map[string]string{
			"field": "username",
		},
		ServiceMetaData: map[string]string{
			"requestId": "abc123",
		},
	}

	expected := "{publicStatus: BadRequest_MissingField (4001), serviceStatus: BadRequest_MissingField (4001), publicMessage: 'Missing required field', serviceMessage: 'Field 'username' is missing in the payload', publicMetaData: {field: 'username'}, serviceMetaData: {requestId: 'abc123'}}"

	actual := err.Error()

	if actual != expected {
		t.Errorf("unexpected error string.\nExpected:\n%s\nGot:\n%s", expected, actual)
	}
}
