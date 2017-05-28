package authorizer

import "testing"
import "github.com/stretchr/testify/assert"

func TestUnauthorizer(t *testing.T) {
	allowed, err := NonResourceUnauthorizer{}.IsAuthorized(nil, nil)
	assert.False(t, allowed, "always false")
	assert.Nil(t, err, "error should be nil")

	allowed, err = ResourceUnauthorizer{}.IsAuthorized(nil, nil)
	assert.False(t, allowed, "always false")
	assert.Nil(t, err, "error should be nil")
}
