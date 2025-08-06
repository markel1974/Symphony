package authenticator

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// SimpleAuthenticator is a struct for managing simple authentication with a username and password hash.
type SimpleAuthenticator struct {
	username      string
	hash          []byte
	authenticated bool
}

// NewSimpleAuthenticator creates a new instance of SimpleAuthenticator with default, uninitialized values.
func NewSimpleAuthenticator() *SimpleAuthenticator {
	a := &SimpleAuthenticator{
		username:      "",
		hash:          []byte{},
		authenticated: false,
	}
	return a
}

// Setup initializes the SimpleAuthenticator with a username and hashed password and returns an error if hashing fails.
func (a *SimpleAuthenticator) Setup(username string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %s", err.Error())
	}
	a.username = username
	a.hash = hashedPassword
	return nil
}

// Authenticate checks if the provided username and password match the stored credentials and sets the authentication status.
func (a *SimpleAuthenticator) Authenticate(user string, pass string) bool {
	if a.username != user {
		a.authenticated = false
	} else {
		a.authenticated = bcrypt.CompareHashAndPassword(a.hash, []byte(pass)) == nil
	}
	return a.authenticated
}

// IsAuthenticated returns the authentication status of the current user.
func (a *SimpleAuthenticator) IsAuthenticated() bool {
	return a.authenticated
}
