package v1beta1

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseToken(t *testing.T) {
	t.Run("valid payload", testValidPayload)
	t.Run("invalid payload", testInvalidPayload)
	t.Run("invalid json", testInvalidJSON)
	t.Run("empty token", testEmptyToken)
}

func testValidPayload(t *testing.T) {
	parser := RequestParser{}
	payload := []byte(`{
		"apiVersion": "authentication.k8s.io/v1beta1",
		"kind": "TokenReview",
		"spec": {
			"token": "my-token"
		}
	}`)

	token, err := parser.ExtractToken(ioutil.NopCloser(bytes.NewReader(payload)))
	assert.Nil(t, err)
	assert.Equal(t, token, "my-token")
}

func testInvalidPayload(t *testing.T) {
	parser := RequestParser{}
	payload := []byte(`{
		"apiVersion": "authentication.k8s.io/v1",
		"kind": "TokenReview",
		"spec": {
			"token": "my-token"
		}
	}`)

	_, err := parser.ExtractToken(ioutil.NopCloser(bytes.NewReader(payload)))
	assert.NotNil(t, err)
}

func testInvalidJSON(t *testing.T) {
	parser := RequestParser{}
	payload := []byte(`{
		"apiVersion": "authentication.k8s.io/v1",
		"kind": "TokenReview"
			"token": "my-token"
		}
	}`)

	_, err := parser.ExtractToken(ioutil.NopCloser(bytes.NewReader(payload)))
	assert.NotNil(t, err)
}

func testEmptyToken(t *testing.T) {
	parser := RequestParser{}
	payload := []byte(`{
		"apiVersion": "authentication.k8s.io/v1beta1",
		"kind": "TokenReview",
		"spec": {
			"token": ""
		}
	}`)

	_, err := parser.ExtractToken(ioutil.NopCloser(bytes.NewReader(payload)))
	assert.NotNil(t, err)
}
