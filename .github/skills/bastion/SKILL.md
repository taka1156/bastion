# Bastion Project Skill

## Overview

**bastion** is a Go CLI tool for SSH host management and file synchronization, driven by a single `bastion.json` configuration file. It supports SSH sessions, rsync-style file sync via SFTP, and Cloudflare Tunnel routing.

Module: `github.com/taka1156/bastion`
Entry point binary: `cmd/bastion`

---

## Architecture

Dependencies only point inward: `cmd` → `app` → `workflow` → `domain`, with `input` and `infra` also depending only on `domain/entity`.

```
cmd/bastion/
  main.go                         Entry point — calls app.NewApp().Run() only

internal/
  domain/
    entity/
      config.go                   BastionConfig, HostConfig, CloudflareConfig,
                                  CommandlineMode (closed union), ClientConfig + accessors

  input/
    cli.go                        CLI flag parsing — ClientInput.GetInput(args)
    json.go                       bastion.json loader — JsonInput.Load(path)

  app/
    app.go                        DI wiring (NewApp) + Run() dispatch loop
    interfaces.go                 initializeWorkflow / sshWorkflow / rsyncWorkflow interfaces

  workflow/
    initialize/workflow.go        init subcommand — writes default bastion.json via filewriter
    ssh/workflow.go               ssh subcommand — delegates to infra/ssh
    rsync/workflow.go             rsync subcommand — delegates to infra/sftp

  infra/
    ssh/client.go                 SSH dial + interactive session (golang.org/x/crypto/ssh)
    sftp/client.go                SFTP-based rsync-style sync (github.com/pkg/sftp)
    filewriter/json.go            JSON file writer — LocalFileWriter.Write(config, outputDir)
```

---

## Key Types

| Type | Package | Purpose |
|---|---|---|
| `BastionConfig` | `domain/entity` | Root config: `$schema`, `Cloudflare *CloudflareConfig`, `Hosts []HostConfig` |
| `HostConfig` | `domain/entity` | SSH target: `name`, `ip`, `user`, `port`, `key` |
| `CloudflareConfig` | `domain/entity` | Tunnel settings: `enableTunnel`, `tunnelToken`, `subdomain`, `domain` |
| `CommandlineMode` | `domain/entity` | Closed union of modes: `Initialize`, `Ssh`, `Rsync`, `Sftp`, `Update`, `Version` |
| `ClientConfig` | `domain/entity` | Parsed CLI flags with nil-safe accessors (`OutputDirValue()`, `LocalPathValue()`, `LangValue()`) |
| `ClientInput` | `input` | Parses `os.Args` into `ClientConfig` using `flag.FlagSet` per subcommand |
| `JsonInput` | `input` | Loads and JSON-decodes `bastion.json` into `BastionConfig` |
| `LocalFileWriter` | `infra/filewriter` | Writes indented JSON to `{outputDir}/bastion.json` |
| `App` | `app` | Owns `ClientInput`, `JsonInput`, and all workflow instances |

---

## Subcommands

| Subcommand | Workflow | Flags | Status |
|---|---|---|---|
| `init` | `workflow/initialize` | `-output <dir>` (default: `.`) | Implemented |
| `ssh` | `workflow/ssh` | — | Implemented |
| `rsync` | `workflow/rsync` | `-local <dir>` (default: `./local`) | Implemented |
| `sftp` | — | — | Work in progress |
| `update` | — | `-lang en\|ja` | Work in progress |
| `version` | — | — | Implemented (prints `app.Version`) |

`app.Version` is injected at build time via ldflags: `-X github.com/taka1156/bastion/internal/app.Version=$(VERSION)`

---

## Configuration File: bastion.json

```json
{
  "$schema": "./bastion.schema.json",
  "hosts": [
    {
      "name": "production",
      "ip": "xxx.xxx.xxx.xxx",
      "user": "ec2-user",
      "port": 22,
      "key": "key/prod.pem",
      "cloudflare": {
        "use_tunnel": true,
        "tunnel_token": "xxx",
        "subdomain": "prod",
        "domain": "example.com"
      }
    }
  ]
}
```

Schema validation is defined in `bastion.schema.json` (JSON Schema draft-07). `hosts` is required and must have at least one entry. `name`, `ip`, `user` are required per host.

---

## Build & Development Commands

| Command | Description |
|---|---|
| `make run` | `go run ./cmd/bastion` |
| `make build` | Build to `bin/bastion` with version ldflags |
| `make test` | `go test ./...` |
| `make test-cover` | Coverage report → `cover.html` |
| `make fmt` | `go fmt ./...` + `golangci-lint run ./...` |
| `make dist` | Cross-compile to `dist/` for linux/darwin/windows × amd64/arm64 |

Scripts:
- `scripts/build.sh` — single-target cross-compile; called by `make dist`
- `scripts/install.sh` — curl installer; downloads latest release from GitHub
- `scripts/tagged.sh` — creates and optionally pushes a git tag

---

## SSH Implementation Notes

- Uses `golang.org/x/crypto/ssh` for SSH connection and session management
- Authenticates with a public key read from `host.Key` path
- Validates server identity against `~/.ssh/known_hosts`
- `infra/ssh.ExecuteSession` always connects to `config.Hosts[0]`
- `infra/sftp.ExecuteRsync` does rsync-like sync: skips unchanged files (same size + mtime), sets mtime after upload

---

## Conventions

- All workflow structs accept interfaces, not concrete types — enables testing with mocks
- `app/interfaces.go` defines the interface boundary between `app` and `workflow`
- `workflow/initialize` depends on a `configWriter` interface (not `filewriter.LocalFileWriter` directly)
- No external CLI frameworks — uses stdlib `flag.FlagSet` per subcommand
- Error wrapping follows `fmt.Errorf("context: %w", err)` pattern throughout
