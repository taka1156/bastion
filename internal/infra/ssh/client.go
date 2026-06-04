package ssh

import (
	"fmt"
	"os"
	"path/filepath"

	gossh "golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"

	"github.com/taka1156/bastion/internal/domain/entity"
)

func buildSSHConfig(host entity.HostConfig) (*gossh.ClientConfig, error) {
	key, err := os.ReadFile(host.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to read SSH key: %w", err)
	}
	signer, err := gossh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SSH key: %w", err)
	}

	knownHostsPath := filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
	hostKeyCallback, err := knownhosts.New(knownHostsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load known_hosts: %w", err)
	}

	return &gossh.ClientConfig{
		User:            host.User,
		Auth:            []gossh.AuthMethod{gossh.PublicKeys(signer)},
		HostKeyCallback: hostKeyCallback,
	}, nil
}

// Dial establishes an SSH connection to the given host.
func Dial(host entity.HostConfig) (*gossh.Client, error) {
	cfg, err := buildSSHConfig(host)
	if err != nil {
		return nil, err
	}
	return gossh.Dial("tcp", host.IP+":"+host.Port, cfg)
}

// ExecuteSession opens an interactive SSH session (bash) on the first host of config.
func ExecuteSession(config entity.BastionConfig) error {
	client, err := Dial(config.Hosts[0])
	if err != nil {
		return fmt.Errorf("failed to connect to SSH server: %w", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close SSH client: %v\n", err)
		}
	}()

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer func() {
		if err := session.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close SSH session: %v\n", err)
		}
	}()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	if err := session.Run("bash"); err != nil {
		return fmt.Errorf("failed to run SSH command: %w", err)
	}

	return nil
}
