package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/ideahitme/k8s-api-webhook/authn/unversioned"
	"github.com/ideahitme/k8s-api-webhook/internal/testutils"
	"github.com/stretchr/testify/assert"
)

var _ Authenticator = Aggregator{}

func TestNewAggregator(t *testing.T) {
	a := NewAggregator()
	assert.Nil(t, a.providers)

	a = NewAggregator(Aggregator{}, Aggregator{})
	assert.Len(t, a.providers, 2)
}

func TestAuthenticate(t *testing.T) {
	usersA := [][]string{
		{
			"1", "kubelet", "kubelet",
		},
		{
			"2", "controller-manager", "controller-manager", "admin", "owner",
		},
	}

	usersB := [][]string{
		{
			"3", "foo", "bar",
		},
		{
			"4", "baz", "qux",
		},
	}

	tmpA := testutils.GenerateTestData(usersA)
	tmpB := testutils.GenerateTestData(usersB)
	defer os.Remove(tmpA.Name())
	defer os.Remove(tmpB.Name())

	authnA, _ := NewStaticAuthenticator(tmpA.Name())
	authnB, _ := NewStaticAuthenticator(tmpB.Name())

	aggr := NewAggregator(authnA, authnB)

	user1, err := aggr.Authenticate("1")
	assert.Nil(t, err)
	assert.Equal(t, &unversioned.UserInfo{
		Name: usersA[0][1],
		UID:  usersA[0][2],
	}, user1)

	user2, err := aggr.Authenticate("3")
	assert.Nil(t, err)
	assert.Equal(t, &unversioned.UserInfo{
		Name: usersB[0][1],
		UID:  usersB[0][2],
	}, user2)

	userEmpty, err := aggr.Authenticate("5")
	assert.Nil(t, err)
	assert.Nil(t, userEmpty)

	faggr := NewAggregator(ErrAuthenticator{})
	_, err = faggr.Authenticate("123")
	assert.NotNil(t, err)
}

type ErrAuthenticator struct {
}

func (ma ErrAuthenticator) Authenticate(string) (*unversioned.UserInfo, error) {
	return nil, fmt.Errorf("failed to authenticate")
}
