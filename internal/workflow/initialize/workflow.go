package initialize

import (
	"github.com/taka1156/bastion/internal/domain/entity"
)

type configWriter interface {
	Write(config entity.BastionConfig, outputDir string) error
}

type Workflow struct {
	writer configWriter
}

func NewWorkflow(writer configWriter) *Workflow {
	return &Workflow{writer: writer}
}

func (w *Workflow) Execute(outputDir string) error {
	return w.writer.Write(defaultBastionConfig(), outputDir)
}

func defaultBastionConfig() entity.BastionConfig {
	cloudflare := entity.CloudflareConfig{
		EnableTunnel: true,
		TunnelToken:  "xxx",
		Subdomain:    "example",
		Domain:       "example.com",
	}
	host := entity.HostConfig{
		Name: "production",
		IP:   "xxxx.xxxx.xxxx.xxxx",
		User: "root",
		Port: "22",
		Key:  "~/.ssh/id_rsa",
	}
	return entity.BastionConfig{
		Schema:     "./bastion.schema.json",
		Cloudflare: &cloudflare,
		Hosts:      []entity.HostConfig{host},
	}
}
