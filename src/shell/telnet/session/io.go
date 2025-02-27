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

package session

import (
	"net"
	"time"
)

// RFC 854: http://tools.ietf.org/html/rfc854, http://support.microsoft.com/kb/231866

type Telnet struct {
	conn net.Conn
	p    *processor
}

func NewTelnet(conn net.Conn) *Telnet {
	t := &Telnet{
		conn: conn,
		p:    newProcessor(),
	}
	return t
}

func (t *Telnet) Write(p []byte) (int, error) {
	return t.conn.Write(p)
}

func (t *Telnet) Read(p []byte) (int, error) {
	for {
		var err error
		var n int
		buf := make([]byte, 1024)

		n, err = t.conn.Read(buf)
		t.p.addBytes(buf[:n])
		if err != nil {
			return 0, err
		}

		n, err = t.p.Read(p)
		if n > 0 {
			return n, err
		}
	}
}

func (t *Telnet) Data(code IOCode) []byte {
	return t.p.subData[code]
}

func (t *Telnet) SetListenFunc(listenFunc func(IOCode, []byte)) {
	t.p.listenFunc = listenFunc
}

func (t *Telnet) Close() error {
	return t.conn.Close()
}

func (t *Telnet) LocalAddr() net.Addr {
	return t.conn.LocalAddr()
}

func (t *Telnet) RemoteAddr() net.Addr {
	return t.conn.RemoteAddr()
}

func (t *Telnet) SetDeadline(dl time.Time) error {
	return t.conn.SetDeadline(dl)
}

func (t *Telnet) SetReadDeadline(dl time.Time) error {
	return t.conn.SetReadDeadline(dl)
}

func (t *Telnet) SetWriteDeadline(dl time.Time) error {
	return t.conn.SetWriteDeadline(dl)
}

func (t *Telnet) WillSga() {
	t.SendCommand(WILL, SGA)
}

func (t *Telnet) WillEcho() {
	t.SendCommand(WILL, ECHO)
}

func (t *Telnet) WontEcho() {
	t.SendCommand(WONT, ECHO)
}

func (t *Telnet) DoWindowSize() {
	t.SendCommand(DO, WS)
}

func (t *Telnet) DoTerminalType() {
	// See http://tools.ietf.org/html/rfc884
	t.SendCommand(DO, TT, IAC, SB, TT, 1, IAC, SE) // 1 = SEND
}

func (t *Telnet) SendCommand(codes ...IOCode) {
	_, _ = t.conn.Write(t.BuildCommand(codes...))
}

func (t *Telnet) BuildCommand(codes ...IOCode) []byte {
	command := make([]byte, len(codes)+1)
	command[0] = codeToByte[IAC]

	for i, code := range codes {
		command[i+1] = codeToByte[code]
	}

	return command
}
