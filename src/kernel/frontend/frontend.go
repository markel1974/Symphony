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

package frontend

import (
	"github.com/markel1974/c64emu/src/kernel/adaptiveticker"
	"github.com/markel1974/c64emu/src/kernel/frontend/authenticator"
	"github.com/markel1974/c64emu/src/kernel/frontend/ssh"
	"github.com/markel1974/c64emu/src/kernel/frontend/telnet"
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/servers/shell"
)

// IShellServer defines the behavior of a customizable command-line shell interface.
// SetPrompt sets the custom prompt for the shell.
// SetTemplate sets the command template to be used by the shell.
// Start initiates the shell and blocks until execution ends.
// AsyncStart starts the shell execution asynchronously without blocking.
type IShellServer interface {
	SetPrompt(prompt string)
	SetTemplate(template *shell.Command)
	Start()
	AsyncStart()
}

// NewFrontend creates a new shell server instance, either SSH or Telnet, based on the secure parameter.
// It uses the provided authenticator, port, and autosave settings.
// A default simple authenticator is used if none is provided.
func NewFrontend(secure bool, auth interfaces.IAuthenticator, port int, autosave bool) IShellServer {
	var ticker = adaptiveticker.NewAdaptiveTicker()
	if auth == nil {
		auth = authenticator.NewSimpleAuthenticator()
	}
	if secure {
		return ssh.NewServer(ticker, auth, port, autosave)
	} else {
		return telnet.NewServer(ticker, auth, port, autosave)
	}
}
