package ssh

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/naturesh/mcloud/internal/core"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type Client struct {
	conn *ssh.Client
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Run(cmd string) error {
	session, err := c.conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	return session.Run(cmd)
}

func Connect(ip string) (*Client, error) {
	var authMethods []ssh.AuthMethod

	if sock := os.Getenv("SSH_AUTH_SOCK"); sock != "" {
		if conn, err := net.Dial("unix", sock); err == nil {
			authMethods = append(authMethods, ssh.PublicKeysCallback(agent.NewClient(conn).Signers))
		}
	}

	if len(authMethods) == 0 {
		return nil, core.ErrSSHAgentNotFound
	}

	config := &ssh.ClientConfig{
		User:            "root",
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	conn, err := ssh.Dial("tcp", ip+":22", config)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", core.ErrSSHConnect, err)
	}

	return &Client{conn: conn}, nil
}
