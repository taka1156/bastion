package app

import "github.com/taka1156/bastion/internal/domain/entity"

type initializeWorkflow interface {
	Execute(outputDir string) error
}

type sshWorkflow interface {
	Execute(config entity.BastionConfig) error
}

type rsyncWorkflow interface {
	Execute(config entity.BastionConfig, localPath string) error
}
