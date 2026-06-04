package ssh

import (
	"github.com/taka1156/bastion/internal/domain/entity"
	sshinfra "github.com/taka1156/bastion/internal/infra/ssh"
)

type Workflow struct{}

func NewWorkflow() *Workflow {
	return &Workflow{}
}

func (w *Workflow) Execute(config entity.BastionConfig) error {
	return sshinfra.ExecuteSession(config)
}
