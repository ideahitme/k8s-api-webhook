package authenticator

import "testing"
import "github.com/stretchr/testify/assert"

func TestNoop(t *testing.T) {
	user, err := Noop{}.Authenticate("1234")
	assert.Nil(t, user, "Noop authenticator returns empty user")
	assert.NoError(t, err, "Noop authenticate should not return any error")
}
