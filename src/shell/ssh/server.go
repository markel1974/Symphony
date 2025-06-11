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

package ssh

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"github.com/markel1974/c64emu/src/shell/adaptiveticker"
	"github.com/markel1974/c64emu/src/shell/cli"
	"github.com/markel1974/c64emu/src/shell/context"
	"github.com/markel1974/c64emu/src/shell/interfaces"
	"github.com/markel1974/c64emu/src/shell/terminal"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"os"
)

type Server struct {
	ticker             *adaptiveticker.AdaptiveTicker
	template           *cli.Command
	prompt             string
	addr               string
	factory            *terminal.EquipmentFactory
	authorized         map[string]bool
	config             *ssh.ServerConfig
	initialized        bool
	privateKeyFilename string
	debug              bool
	auth               interfaces.IAuthenticator
	autosave           bool
}

func NewServer(ticker *adaptiveticker.AdaptiveTicker, auth interfaces.IAuthenticator, port int, autosave bool) *Server {
	return &Server{
		ticker:             ticker,
		addr:               fmt.Sprintf("%s:%d", "127.0.0.1", port),
		factory:            terminal.NewEquipmentFactory(),
		authorized:         make(map[string]bool),
		auth:               auth,
		initialized:        false,
		privateKeyFilename: "id_rsa",
		debug:              false,
		autosave:           autosave,
	}
}

func (r *Server) Setup() {
	if r.initialized {
		return
	}

	if authorizedKeys, err := os.ReadFile("authorized_keys"); err == nil {
		for len(authorizedKeys) > 0 {
			pubKey, _, _, rest, err := ssh.ParseAuthorizedKey(authorizedKeys)
			if err != nil {
				log.Println(err)
				break
			}
			r.authorized[string(pubKey.Marshal())] = true
			authorizedKeys = rest
		}
	}

	r.config = &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			if r.auth.Authenticate(c.User(), string(pass)) {
				return nil, nil
			}
			return nil, fmt.Errorf("password rejected for %q", c.User())
		},

		PublicKeyCallback: func(c ssh.ConnMetadata, pubKey ssh.PublicKey) (*ssh.Permissions, error) {
			if r.authorized[string(pubKey.Marshal())] {
				return &ssh.Permissions{
					Extensions: map[string]string{"pubkey-fp": ssh.FingerprintSHA256(pubKey)},
				}, nil
			}
			return nil, fmt.Errorf("unknown public key for %q", c.User())
		},
	}

	var signer ssh.Signer

	if _, err := os.Stat(r.privateKeyFilename); err == nil {
		privateBytes, err := os.ReadFile(r.privateKeyFilename)
		if err != nil {
			log.Fatal("Failed to parse private key: ", err)
		}

		if signer, err = ssh.ParsePrivateKey(privateBytes); err != nil {
			log.Fatal("Failed to create signer: ", err)
		}

		//log.Println("Private key", r.privateKeyFilename ,"successfully loaded")
	} else {
		//log.Println("Trying to generate Private key...")

		private, err := r.generatePrivateKey(4096)
		if err != nil {
			log.Fatal("Failed to generate private key: ", err)
		}

		if signer, err = ssh.NewSignerFromKey(private); err != nil {
			log.Fatal("Failed to create signer: ", err)
		}

		if err = r.savePrivateKey(r.privateKeyFilename, private); err != nil {
			log.Fatal("Failed to save Private key: ", err)
		}

		//log.Println("Private key", r.privateKeyFilename, "successfully generated")
	}

	r.config.AddHostKey(signer)

	r.initialized = true
}

func (r *Server) SetPrompt(prompt string) {
	r.prompt = prompt
}

func (r *Server) SetTemplate(template *cli.Command) {
	r.template = template
}

func (r *Server) Start() {
	r.Setup()

	listener, err := net.Listen("tcp", r.addr)
	if err != nil {
		log.Fatal("failed to listen for connection: ", err)
	}

	for {
		nConn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			return
		}

		go r.handleConnection(nConn)
	}
}

func (r *Server) AsyncStart() {
	go func() {
		r.Start()
	}()
}

func (r *Server) handleConnection(nConn net.Conn) {
	conn, channels, reqs, err := ssh.NewServerConn(nConn, r.config)
	if err != nil {
		log.Println("failed to handshake: ", err)
		return
	}

	if r.debug {
		log.Println("Connected from", string(conn.ClientVersion()))
	}

	// The incoming Request channel must be serviced.
	go ssh.DiscardRequests(reqs)

	for newChannel := range channels {
		// Channels have a type, depending on the application level protocol intended.
		// In the case of a shell, the type is "session" and ServerShell may be used to present a terminal interface.
		if newChannel.ChannelType() != "session" {
			_ = newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}
		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Println("Could not accept channel:", err)
			continue
		}

		r.auth.IsAuthenticated()

		ctx := context.NewContext(r.ticker, channel, channel, r.auth, r.template, r.prompt, r.autosave)
		ioAdapter := interfaces.IInputOutput(ctx)
		term := r.factory.Create("VT100", ioAdapter, -1)
		ctx.Setup(term)
		//ctx.setAuthenticatedState()
		//ctx.SetEnterKey(10)

		// out-of-band requests
		go func(in <-chan *ssh.Request) {
			for req := range in {
				switch req.Type {
				case "shell":
					if len(req.Payload) == 0 {
						_ = req.Reply(true, nil)
					}
				case "pty-req":
					termLen := req.Payload[3]
					w, h := r.parseSize(req.Payload[termLen+4:])
					ctx.SetScreenSize(int(w), int(h))
					_ = req.Reply(true, nil)
				case "window-change":
					w, h := r.parseSize(req.Payload)
					ctx.SetScreenSize(int(w), int(h))
				}
			}
		}(requests)

		ctx.Exec()

		_ = channel.Close()
	}
}

func (r *Server) generatePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	reader := rand.Reader
	return rsa.GenerateKey(reader, bitSize)
}

func (r *Server) savePrivateKey(filename string, key *rsa.PrivateKey) error {
	keyBytes := x509.MarshalPKCS1PrivateKey(key)
	keyPem := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyBytes})
	return os.WriteFile(filename, keyPem, 0600)
}

func (r *Server) parseSize(b []byte) (uint32, uint32) {
	w := binary.BigEndian.Uint32(b)
	h := binary.BigEndian.Uint32(b[4:])
	return w, h
}
