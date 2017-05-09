package v1beta1

import "testing"

func TestParseToken(t *testing.T) {
	t.Run("valid payload", testValidPayload)
	t.Run("invalid payload", testInvalidPayload)
}

func testValidPayload(t *testing.T) {
	payload := []byte(`{
		"apiVersion": "authentication.k8s.io/v1beta1",
		"kind": "TokenReview",
		"spec": {
			"token": "my-token"
		}
	}`)

	token, err := ParseToken(payload)
	if err != nil {
		t.Fatal(err)
	}
	if token != "my-token" {
		t.Errorf("failed to retrieve the token. expected: my-token, got: %s", token)
	}
}

func testInvalidPayload(t *testing.T) {
	payload := []byte(`{
		"apiVersion": "authentication.k8s.io/v1",
		"kind": "TokenReview",
		"spec": {
			"token": "my-token"
		}
	}`)

	_, err := ParseToken(payload)
	if err == nil {
		t.Fatal("shoud return error")
	}
}
