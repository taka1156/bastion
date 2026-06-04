package app

import (
	"fmt"
	"os"

	"github.com/taka1156/bastion/internal/infra/filewriter"
	"github.com/taka1156/bastion/internal/input"
	"github.com/taka1156/bastion/internal/workflow/initialize"
	rsyncwf "github.com/taka1156/bastion/internal/workflow/rsync"
	sshwf "github.com/taka1156/bastion/internal/workflow/ssh"
)

// Version is set via ldflags at build time.
var Version = "dev"

type workflowCases struct {
	initialize initializeWorkflow
	ssh        sshWorkflow
	rsync      rsyncWorkflow
}

type App struct {
	clientInput *input.ClientInput
	jsonInput   *input.JsonInput
	flows       workflowCases
}

func NewApp() *App {
	writer := filewriter.NewLocalFileWriter()

	return &App{
		clientInput: input.NewClientInput(),
		jsonInput:   input.NewJsonInput(),
		flows: workflowCases{
			initialize: initialize.NewWorkflow(writer),
			ssh:        sshwf.NewWorkflow(),
			rsync:      rsyncwf.NewWorkflow(),
		},
	}
}

func (a *App) Run() error {
	args := os.Args
	clientConfig := a.clientInput.GetInput(args)

	switch clientConfig.Mode.CommandlineModeValue() {
	case "init":
		if err := a.flows.initialize.Execute(clientConfig.OutputDirValue()); err != nil {
			return fmt.Errorf("init: %w", err)
		}
		return nil
	case "version":
		fmt.Println(Version)
		return nil
	case "update":
		fmt.Println("Update: work in progress...")
		return nil
	case "sftp":
		fmt.Println("SFTP: work in progress...")
		return nil
	}

	// Subcommands that require loading bastion.json
	config, err := a.jsonInput.Load("bastion.json")
	if err != nil {
		return err
	}

	switch clientConfig.Mode.CommandlineModeValue() {
	case "ssh":
		return a.flows.ssh.Execute(config)
	case "rsync":
		return a.flows.rsync.Execute(config, clientConfig.LocalPathValue())
	default:
		fmt.Fprintf(os.Stderr, "No valid mode selected. Use 'init', 'ssh', 'rsync', 'sftp', 'update', or 'version' subcommands.\n")
		return nil
	}
}
