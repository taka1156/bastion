package entity

// CommandlineMode represents the subcommand mode passed to the CLI.
// The unexported val field prevents external construction of arbitrary values,
// simulating a closed union type: Initialize | Update | Version.
type CommandlineMode struct {
	val string
}

// Predefined CommandlineMode values. These are the only valid modes.
var (
	Initialize = CommandlineMode{val: "init"}
	Ssh        = CommandlineMode{val: "ssh"}
	Rsync      = CommandlineMode{val: "rsync"}
	Sftp       = CommandlineMode{val: "sftp"}
	Update     = CommandlineMode{val: "update"}
	Version    = CommandlineMode{val: "version"}
)

func (c CommandlineMode) CommandlineModeValue() string {
	switch c.val {
	case Initialize.val:
		return "init"
	case Ssh.val:
		return "ssh"
	case Rsync.val:
		return "rsync"
	case Sftp.val:
		return "sftp"
	case Update.val:
		return "update"
	case Version.val:
		return "version"
	default:
		return "undefined"
	}
}

type ClientConfig struct {
	Mode      CommandlineMode
	OutputDir *string
	LocalPath *string
	Lang      *string
}

func (c ClientConfig) OutputDirValue() string {
	if c.OutputDir == nil {
		return "."
	}
	return *c.OutputDir
}

func (c ClientConfig) LocalPathValue() string {
	if c.LocalPath == nil {
		return "./local"
	}
	return *c.LocalPath
}

func (c ClientConfig) LangValue() string {
	if c.Lang == nil {
		return ""
	}
	return *c.Lang
}

type CloudflareConfig struct {
	EnableTunnel bool   `json:"enableTunnel"`
	TunnelToken  string `json:"tunnelToken"`
	Subdomain    string `json:"subdomain"`
	Domain       string `json:"domain"`
}

type HostConfig struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
	User string `json:"user"`
	Port string `json:"port"`
	Key  string `json:"key"`
}

type BastionConfig struct {
	Schema     string            `json:"$schema"`
	Cloudflare *CloudflareConfig `json:"cloudflare,omitempty"`
	Hosts      []HostConfig      `json:"hosts"`
}
