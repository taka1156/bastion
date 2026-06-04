package sftp

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	pkgsftp "github.com/pkg/sftp"

	"github.com/taka1156/bastion/internal/domain/entity"
	sshinfra "github.com/taka1156/bastion/internal/infra/ssh"
)

func uploadFile(sftpClient *pkgsftp.Client, root *os.Root, localPath, relPath string, info os.FileInfo) error {
	remotePath := filepath.Join(localPath, relPath)

	// rsync-like: skip if remote file has same size and mtime
	if remoteInfo, err := sftpClient.Stat(remotePath); err == nil {
		if remoteInfo.Size() == info.Size() && remoteInfo.ModTime().Equal(info.ModTime()) {
			fmt.Printf("Skipped (unchanged): %s\n", relPath)
			return nil
		}
	}

	localFile, err := root.Open(relPath)
	if err != nil {
		return fmt.Errorf("failed to open local file: %w", err)
	}
	defer func() {
		if err := localFile.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close local file: %v\n", err)
		}
	}()

	remoteFile, err := sftpClient.Create(remotePath)
	if err != nil {
		return fmt.Errorf("failed to create remote file: %w", err)
	}
	defer func() {
		if err := remoteFile.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close remote file: %v\n", err)
		}
	}()

	if _, err := io.Copy(remoteFile, localFile); err != nil {
		return fmt.Errorf("failed to upload %s: %w", relPath, err)
	}

	// rsync-like: sync mtime after upload
	if err := sftpClient.Chtimes(remotePath, info.ModTime(), info.ModTime()); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set mtime for %s: %v\n", relPath, err)
	}

	fmt.Printf("Uploaded: %s\n", relPath)
	return nil
}

// ExecuteRsync syncs localPath to the first host in config via SFTP.
func ExecuteRsync(config entity.BastionConfig, localPath string) error {
	conn, err := sshinfra.Dial(config.Hosts[0])
	if err != nil {
		return fmt.Errorf("failed to connect to SSH server: %w", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close SSH connection: %v\n", err)
		}
	}()

	sftpClient, err := pkgsftp.NewClient(conn)
	if err != nil {
		return fmt.Errorf("failed to create SFTP client: %w", err)
	}
	defer func() {
		if err := sftpClient.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close SFTP client: %v\n", err)
		}
	}()

	root, err := os.OpenRoot(localPath)
	if err != nil {
		return fmt.Errorf("failed to open local root: %w", err)
	}
	defer func() {
		if err := root.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close local root: %v\n", err)
		}
	}()

	if err := filepath.Walk(localPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(localPath, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		if info.IsDir() {
			return sftpClient.MkdirAll(filepath.Join(localPath, relPath))
		}

		return uploadFile(sftpClient, root, localPath, relPath, info)
	}); err != nil {
		return fmt.Errorf("failed to upload files: %w", err)
	}

	return nil
}
