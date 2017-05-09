package unversioned

// UserInfo common concept for storing and extracting user information
type UserInfo struct {
	UID    string
	Name   string
	Groups []string
}
