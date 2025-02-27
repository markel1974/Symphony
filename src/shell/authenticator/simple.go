/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package authenticator

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type SimpleAuthenticator struct {
	username      string
	hash          []byte
	authenticated bool
}

func NewSimpleAuthenticator() *SimpleAuthenticator {
	a := &SimpleAuthenticator{
		username:      "",
		hash:          []byte{},
		authenticated: false,
	}
	return a
}

func (a *SimpleAuthenticator) Setup(username string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %s", err.Error())
	}
	a.username = username
	a.hash = hashedPassword
	return nil
}

func (a *SimpleAuthenticator) Authenticate(user string, pass string) bool {
	if a.username != user {
		a.authenticated = false
	} else {
		a.authenticated = bcrypt.CompareHashAndPassword(a.hash, []byte(pass)) == nil
	}
	return a.authenticated
}

func (a *SimpleAuthenticator) IsAuthenticated() bool {
	return a.authenticated
}
