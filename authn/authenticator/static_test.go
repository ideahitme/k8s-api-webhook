package authenticator

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ideahitme/k8s-api-webhook/authn/unversioned"
	"github.com/ideahitme/k8s-api-webhook/internal/testutils"
)

var _ Authenticator = Static{}

func TestNewStatic(t *testing.T) {
	_, err := NewStatic("")
	assert.NotNil(t, err)

	for _, ti := range []struct {
		title          string
		lines          [][]string
		expectError    bool
		expectedOutput Static
	}{
		{
			title: "invalid csv",
			lines: [][]string{
				{
					"123456", "kubelet",
				},
				{
					"123456", "controller-manager", "controller-manager",
				},
			},
			expectError:    true,
			expectedOutput: map[string]*unversioned.UserInfo{},
		},
		{
			title: "valid csv",
			lines: [][]string{
				{
					"1234567kubelet", "kubelet", "kubelet",
				},
				{
					"1234567manager", "controller-manager", "controller-manager", "admin", "owner",
				},
			},
			expectError: false,
			expectedOutput: map[string]*unversioned.UserInfo{
				"1234567kubelet": {"kubelet", "kubelet", nil},
				"1234567manager": {"controller-manager", "controller-manager", []string{"admin", "owner"}},
			},
		},
	} {
		t.Run(ti.title, func(t *testing.T) {
			tmpFile := testutils.GenerateTestData(ti.lines)
			defer os.Remove(tmpFile.Name())
			authn, err := NewStatic(tmpFile.Name())
			if ti.expectError {
				assert.NotNil(t, err)
			}
			if !ti.expectError {
				assert.Nil(t, err)
				assert.Equal(t, ti.expectedOutput, authn)
			}
		})
	}
}

func TestStaticAuthenticate(t *testing.T) {
	for _, ti := range []struct {
		title        string
		lines        [][]string
		token        string
		expectedUser *unversioned.UserInfo
	}{
		{
			title: "no user",
			lines: [][]string{
				{
					"1234531226", "kubelet",
				},
				{
					"12345689", "controller-manager", "controller-manager",
				},
			},
			token:        "14123123",
			expectedUser: nil,
		},
		{
			title: "valid user",
			lines: [][]string{
				{
					"1234567kubelet", "kubelet", "kubelet",
				},
				{
					"1234567mvnager", "controller-manager", "controller-manager", "admin", "owner",
				},
			},
			token:        "1234567mvnager",
			expectedUser: &unversioned.UserInfo{"controller-manager", "controller-manager", []string{"admin", "owner"}},
		},
	} {
		t.Run(ti.title, func(t *testing.T) {
			tmpFile := testutils.GenerateTestData(ti.lines)
			defer os.Remove(tmpFile.Name())
			authn, _ := NewStatic(tmpFile.Name())
			user, err := authn.Authenticate(ti.token)
			assert.Nil(t, err)
			assert.Equal(t, ti.expectedUser, user)
		})
	}
}
