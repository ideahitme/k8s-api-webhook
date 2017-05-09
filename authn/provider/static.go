package provider

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/ideahitme/k8s-api-webhook/authn/unversioned"
)

// StaticAuthenticator represents authentication mechanism via hardcoded
// token - user pair passed via files
type StaticAuthenticator map[string]*unversioned.UserInfo

// ReadTokensFile populates StaticAuthneticator object by reading from passed csv file
func ReadTokensFile(filepath string) (StaticAuthenticator, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	tokens := map[string]*unversioned.UserInfo{}
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	for {
		records, err := reader.Read()
		if err == io.EOF {
			break
		}
		if len(records) != 3 {
			return nil, fmt.Errorf(`exactly three columns required in the token file in the format: 
				token, Name, UID`)
		}
		user := &unversioned.UserInfo{
			Name: records[1],
			UID:  records[2],
		}
		tokens[records[0]] = user
	}

	return StaticAuthenticator(tokens), nil
}

// Authenticate looks up a user for the provided token and returns UID
func (tokens StaticAuthenticator) Authenticate(token string) (*unversioned.UserInfo, error) {
	if user, ok := tokens[token]; ok {
		return user, nil
	}
	return nil, nil
}
