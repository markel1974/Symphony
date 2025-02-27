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

package shell

import (
	"github.com/markel1974/c64emu/src/shell/adaptiveticker"
	"github.com/markel1974/c64emu/src/shell/authenticator"
	"github.com/markel1974/c64emu/src/shell/cli"
	"github.com/markel1974/c64emu/src/shell/interfaces"
	"github.com/markel1974/c64emu/src/shell/ssh"
	"github.com/markel1974/c64emu/src/shell/telnet"
)

type IShellServer interface {
	SetPrompt(prompt string)
	SetTemplate(template *cli.Command)
	Start()
	AsyncStart()
}

func New(secure bool, auth interfaces.IAuthenticator, port int, autosave bool) IShellServer {
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
