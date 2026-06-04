# bastion

`bastion.json` による SSH ホスト管理・ファイル同期 CLI ツール。

## curl でインストール

最新リリースを自動でダウンロードし、`/usr/local/bin` にインストールします。

```bash
curl -fsSL https://raw.githubusercontent.com/taka1156/bastion/master/scripts/install.sh | bash
```

インストール先を変更する場合:

```bash
curl -fsSL https://raw.githubusercontent.com/taka1156/bastion/master/scripts/install.sh | INSTALL_DIR=$HOME/.local/bin bash
```

インストール後、`bsn` をショートカットエイリアスとして利用できます。

## アーキテクチャ

- Domain: ルールとモデル
  - `internal/domain/entity`
- App: 構成とオーケストレーション
  - `internal/app`
- Input adapters: CLI / JSON 入力
  - `internal/input`
- Infra: SSH セッション・SFTP/rsync・ファイル書き込み
  - `internal/infra/ssh`
  - `internal/infra/sftp`
  - `internal/infra/filewriter`
- Workflow: ユースケース
  - `internal/workflow/initialize`
  - `internal/workflow/ssh`
  - `internal/workflow/rsync`
- Entry Point: CLI
  - `cmd/bastion`

依存は常に内側へのみ向きます。

## 使い方

### bastion.json の初期化

`init` サブコマンドで `bastion.json` のテンプレートを生成します。

```bash
bastion init
```

| オプション | 既定値 | 説明 |
|---|---|---|
| `-output` | `.` | `bastion.json` の出力先ディレクトリ |

### SSH セッションを開始

```bash
bastion ssh
```

`bastion.json` に定義された最初のホストへのインタラクティブ SSH セッションを開きます。

### rsync でファイルを同期（SFTP）

```bash
bastion rsync -local ./src
```

| オプション | 既定値 | 説明 |
|---|---|---|
| `-local` | `./local` | リモートホストに同期するローカルディレクトリ |

### バージョン表示

```bash
bastion version
```

### bastion のアップデート

```bash
bastion update
```

| オプション | 既定値 | 説明 |
|---|---|---|
| `-lang` | (自動検出) | CLI メッセージの言語（`en` または `ja`） |

## bastion.json フォーマット

JSON Schema バリデーションに対応したエディターでは補完・検証が利用できます。

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

### ホストフィールド

| フィールド | 必須 | 説明 |
|---|---|---|
| `name` | ✓ | ホストの表示名 |
| `ip` | ✓ | 対象サーバーの IP アドレス |
| `user` | ✓ | SSH ログインユーザー |
| `port` | | SSH ポート番号（既定値: `22`） |
| `key` | | SSH 秘密鍵ファイルのパス |
| `cloudflare` | | Cloudflare Tunnel 設定（下記参照） |

### Cloudflare Tunnel フィールド

| フィールド | 必須 | 説明 |
|---|---|---|
| `use_tunnel` | ✓ | Cloudflare Tunnel 経由で SSH をルーティングするか |
| `tunnel_token` | | Cloudflare Tunnel トークン |
| `subdomain` | | トンネルホスト名のサブドメイン |
| `domain` | | トンネルホスト名のドメイン |

## 開発

### 実行

```bash
go run ./cmd/bastion
# または
make run
```

### ビルド

```bash
make build
```

### テスト

```bash
make test
```

### クロスプラットフォームビルド

```bash
make dist
```

`dist/` 以下に以下のバイナリが生成されます:

- `bastion_linux_amd64.tar.gz`
- `bastion_linux_arm64.tar.gz`
- `bastion_darwin_amd64.tar.gz`
- `bastion_darwin_arm64.tar.gz`
- `bastion_windows_amd64.exe`

## ライセンス

MIT
