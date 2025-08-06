package session

import (
	"net"
	"time"
)

// RFC 854: http://tools.ietf.org/html/rfc854, http://support.microsoft.com/kb/231866

// Telnet provides a high-level abstraction for interacting with TELNET protocol connections over a network.
// It encapsulates a net.Conn and a processor for handling TELNET-specific communication and control signals.
type Telnet struct {
	conn net.Conn
	p    *processor
}

// NewTelnet initializes and returns a new Telnet instance using the provided network connection.
func NewTelnet(conn net.Conn) *Telnet {
	t := &Telnet{
		conn: conn,
		p:    newProcessor(),
	}
	return t
}

// Write sends the provided data through the Telnet connection and returns the number of bytes written and any error encountered.
func (t *Telnet) Write(p []byte) (int, error) {
	return t.conn.Write(p)
}

// Read reads data into the provided byte slice and returns the number of bytes read and any error encountered.
// Data is refined by processing internal buffers before being passed to the caller.
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

// Data retrieves the data associated with the given IOCode from the processor's subData map.
func (t *Telnet) Data(code IOCode) []byte {
	return t.p.subData[code]
}

// SetListenFunc assigns a callback function to handle specific IOCode and associated data during session processing.
func (t *Telnet) SetListenFunc(listenFunc func(IOCode, []byte)) {
	t.p.listenFunc = listenFunc
}

// Close terminates the underlying network connection and releases associated resources.
func (t *Telnet) Close() error {
	return t.conn.Close()
}

// LocalAddr returns the local
func (t *Telnet) LocalAddr() net.Addr {
	return t.conn.LocalAddr()
}

// RemoteAddr returns the remote network address of the underlying connection.
func (t *Telnet) RemoteAddr() net.Addr {
	return t.conn.RemoteAddr()
}

// SetDeadline sets the read and write deadlines for the Telnet connection.
func (t *Telnet) SetDeadline(dl time.Time) error {
	return t.conn.SetDeadline(dl)
}

// SetReadDeadline sets the read deadline for the underlying connection to the specified time.
// If the deadline is reached, read operations fail with a timeout error.
// Passing a zero value disables the deadline.
func (t *Telnet) SetReadDeadline(dl time.Time) error {
	return t.conn.SetReadDeadline(dl)
}

// SetWriteDeadline sets the deadline for future Write calls. A deadline is an absolute time after which writes fail.
func (t *Telnet) SetWriteDeadline(dl time.Time) error {
	return t.conn.SetWriteDeadline(dl)
}

// WillSga sends a Telnet WILL command with the SGA (Suppress-go-ahead) option to the remote server.
func (t *Telnet) WillSga() {
	t.SendCommand(WILL, SGA)
}

// WillEcho sends a WILL ECHO command to indicate the local system's desire to begin or continue echoing characters.
func (t *Telnet) WillEcho() {
	t.SendCommand(WILL, ECHO)
}

// WontEcho sends a Telnet command indicating refusal to perform the ECHO option.
func (t *Telnet) WontEcho() {
	t.SendCommand(WONT, ECHO)
}

// DoWindowSize sends a Telnet DO command to negotiate the Window Size option with the remote host.
func (t *Telnet) DoWindowSize() {
	t.SendCommand(DO, WS)
}

// DoTerminalType sends a Telnet command to request the terminal type from the client in compliance with RFC 884.
func (t *Telnet) DoTerminalType() {
	// See http://tools.ietf.org/html/rfc884
	t.SendCommand(DO, TT, IAC, SB, TT, 1, IAC, SE) // 1 = SEND
}

// SendCommand sends a sequence of IOCode commands over the Telnet connection by writing the built command to the connection.
func (t *Telnet) SendCommand(codes ...IOCode) {
	_, _ = t.conn.Write(t.BuildCommand(codes...))
}

// BuildCommand constructs a byte slice representing a Telnet command using the provided IOCodes and the IAC prefix.
func (t *Telnet) BuildCommand(codes ...IOCode) []byte {
	command := make([]byte, len(codes)+1)
	command[0] = codeToByte[IAC]

	for i, code := range codes {
		command[i+1] = codeToByte[code]
	}

	return command
}
