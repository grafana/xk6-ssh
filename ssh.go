package xk6ssh

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"golang.org/x/crypto/ssh"
)

type RSAKey struct {
	file     string
	password string
}

// SSH is the main export of k6 docker extension
type K6SSH struct {
	Version        string
	Session        *ssh.Session
	Client         *ssh.Client
	Config         *ssh.ClientConfig
	Out            *bytes.Buffer
	Stdin          io.WriteCloser
	ConnectOptions ConnectionOptions
}

type ConnectionOptions struct {
	RsaKey       string
	Host         string
	Port         int
	Username     string
	Password     string
	SudoPassword string
}

type RunOptions struct {
	Sudo bool
}

func (k6ssh *K6SSH) rsaKeyAuthMethod(options ConnectionOptions) (ssh.AuthMethod, error) {
	var pk string
	if options.RsaKey != "" {
		pk = options.RsaKey
	} else {
		pk = k6ssh.defaultKeyPath()
	}

	key, err := ioutil.ReadFile(pk)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(signer), nil
}

func (k6ssh *K6SSH) Connect(options ConnectionOptions) error {
	k6ssh.ConnectOptions = options
	var authMethod ssh.AuthMethod
	var err error
	if options.Password != "" {
		authMethod = ssh.Password(options.Password)
	} else {
		authMethod, err = k6ssh.rsaKeyAuthMethod(options)
		if err != nil {
			return err
		}
	}

	k6ssh.Config = &ssh.ClientConfig{
		Config:          ssh.Config{},
		User:            options.Username,
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		ClientVersion:   "",
		Timeout:         0,
	}

	addr := fmt.Sprintf("%s:%d", options.Host, options.Port)
	client, err := ssh.Dial("tcp", addr, k6ssh.Config)
	if err != nil {
		return err
	}
	k6ssh.Client = client
	return nil
}

func (conn *K6SSH) defaultKeyPath() string {
	home := os.Getenv("HOME")
	if len(home) > 0 {
		return path.Join(home, ".ssh/id_rsa")
	}
	return ""
}

func (k6ssh *K6SSH) Run(cmds []string, options RunOptions) (string, error) {
	session, err := k6ssh.Client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	err = session.RequestPty("xterm", 80, 40, modes)
	if err != nil {
		return "", err
	}

	in, err := session.StdinPipe()
	if err != nil {
		return "", err
	}

	out, err := session.StdoutPipe()
	if err != nil {
		return "", err
	}

	var output []byte

	go func(in io.WriteCloser, out io.Reader, output *[]byte) {
		var (
			line string
			r    = bufio.NewReader(out)
		)
		for {
			b, err := r.ReadByte()
			if err != nil {
				break
			}

			*output = append(*output, b)

			if b == byte('\n') {
				line = ""
				continue
			}

			line += string(b)
			if options.Sudo {
				if strings.HasPrefix(line, "[sudo] password for ") && strings.HasSuffix(line, ": ") {
					_, err = in.Write([]byte(k6ssh.ConnectOptions.SudoPassword + "\n"))
					if err != nil {
						break
					}
				}
			}
		}
	}(in, out, &output)

	cmd := strings.Join(cmds, "; ")
	_, err = session.Output(cmd)
	if err != nil {
		return "", err
	}

	return string(output), nil
}
