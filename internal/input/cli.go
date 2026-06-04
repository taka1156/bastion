package input

import (
	"flag"
	"fmt"
	"os"

	"github.com/taka1156/bastion/internal/domain/entity"
)

type ClientInput struct{}

func NewClientInput() *ClientInput {
	return &ClientInput{}
}

func (ci *ClientInput) GetInput(args []string) entity.ClientConfig {
	clientConfig := entity.ClientConfig{}

	if len(args) > 1 {
		switch args[1] {
		case "init":
			initCmd := flag.NewFlagSet("init", flag.ExitOnError)
			clientConfig.OutputDir = initCmd.String("output", ".", "output directory for generated files")
			initCmd.Usage = func() {
				fmt.Fprintf(os.Stderr, "Usage: %s init [options]\n\n", args[0])
				fmt.Fprintf(os.Stderr, "Initialize setting JSON\n")
				initCmd.PrintDefaults()
			}
			_ = initCmd.Parse(args[2:])
			clientConfig.Mode = entity.Initialize
			return clientConfig
		case "ssh":
			sshCmd := flag.NewFlagSet("ssh", flag.ExitOnError)
			sshCmd.Usage = func() {
				fmt.Fprintf(os.Stderr, "Usage: %s ssh\n\n", args[0])
				fmt.Fprintf(os.Stderr, "Start SSH session with bastion\n")
				sshCmd.PrintDefaults()
			}
			_ = sshCmd.Parse(args[2:])
			clientConfig.Mode = entity.Ssh
			return clientConfig
		case "rsync":
			rsyncCmd := flag.NewFlagSet("rsync", flag.ExitOnError)
			clientConfig.LocalPath = rsyncCmd.String("local", "./local", "local directory to sync to remote")
			rsyncCmd.Usage = func() {
				fmt.Fprintf(os.Stderr, "Usage: %s rsync [options]\n\n", args[0])
				fmt.Fprintf(os.Stderr, "Sync a local directory to the remote host\n")
				rsyncCmd.PrintDefaults()
			}
			_ = rsyncCmd.Parse(args[2:])
			clientConfig.Mode = entity.Rsync
			return clientConfig
		case "sftp":
			sftpCmd := flag.NewFlagSet("sftp", flag.ExitOnError)
			sftpCmd.Usage = func() {
				fmt.Fprintf(os.Stderr, "Usage: %s sftp\n\n", args[0])
				fmt.Fprintf(os.Stderr, "Start SFTP session with bastion\n")
				sftpCmd.PrintDefaults()
			}
			_ = sftpCmd.Parse(args[2:])
			clientConfig.Mode = entity.Sftp
			return clientConfig
		case "update":
			updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
			clientConfig.Lang = updateCmd.String("lang", "", "language for CLI messages (en/ja, default: auto-detect)")
			updateCmd.Usage = func() {
				fmt.Fprintf(os.Stderr, "Usage: %s update [options]\n\n", args[0])
				fmt.Fprintf(os.Stderr, "Update bastion to the latest version\n")
				updateCmd.PrintDefaults()
			}
			_ = updateCmd.Parse(args[2:])
			clientConfig.Mode = entity.Update
			return clientConfig
		case "version":
			versionCmd := flag.NewFlagSet("version", flag.ExitOnError)
			versionCmd.Usage = func() {
				fmt.Fprintf(os.Stderr, "Usage: %s version\n\n", args[0])
				fmt.Fprintf(os.Stderr, "Show current version of bastion\n")
				versionCmd.PrintDefaults()
			}
			_ = versionCmd.Parse(args[2:])
			clientConfig.Mode = entity.Version
			return clientConfig
		}
	}

	return clientConfig
}
