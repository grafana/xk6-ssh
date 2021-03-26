package xk6ssh

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"golang.org/x/crypto/ssh"
)

type RSAKey struct {
	file     string
	password string
}

// SSH is the main export of k6 docker extension
type K6SSH struct {
	Version string
	Session *ssh.Session
	Client  *ssh.Client
	Config  *ssh.ClientConfig
	Out     *bytes.Buffer
	Stdin   io.WriteCloser
}

type ConnectionOptions struct {
	RsaKey   string
	Host     string
	Port     int
	Username string
	Password string
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

func (k6ssh *K6SSH) Run(command string) (string, error) {
	session, err := k6ssh.Client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	err = session.Run(command)
	return stdoutBuf.String(), err
}
