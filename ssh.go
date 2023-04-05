package xk6ssh

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/spf13/afero"

	"golang.org/x/crypto/ssh"
)

// K6SSH is the main export of the k6 extension
type K6SSH struct {
	Session *ssh.Session
	Client  *ssh.Client
	Config  *ssh.ClientConfig
	Out     *bytes.Buffer
	Stdin   io.WriteCloser
	fs      afero.Fs
}

// ConnectionOptions provides configuration for the SSH session
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

	key, err := afero.ReadFile(k6ssh.fs, pk)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(signer), nil
}

// Connect starts and SSH session with the provided options
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
		Config: ssh.Config{},
		User:   options.Username,
		Auth:   []ssh.AuthMethod{authMethod},
		// #nosec G106
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

func (k6ssh *K6SSH) defaultKeyPath() string {
	//nolint: forbidigo
	home := os.Getenv("HOME")
	if len(home) > 0 {
		return path.Join(home, ".ssh/id_rsa")
	}
	return ""
}

// Run executes a remote command over SSH
func (k6ssh *K6SSH) Run(command string) (string, error) {
	session, err := k6ssh.Client.NewSession()
	if err != nil {
		return "", err
	}
	defer func() {
		_ = session.Close()
	}()
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	err = session.Run(command)
	return stdoutBuf.String(), err
}
