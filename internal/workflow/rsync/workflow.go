package rsync

import (
	"github.com/taka1156/bastion/internal/domain/entity"
	sftpinfra "github.com/taka1156/bastion/internal/infra/sftp"
)

type Workflow struct{}

func NewWorkflow() *Workflow {
	return &Workflow{}
}

func (w *Workflow) Execute(config entity.BastionConfig, localPath string) error {
	return sftpinfra.ExecuteRsync(config, localPath)
}
