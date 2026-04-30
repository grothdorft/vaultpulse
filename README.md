# vaultpulse

> CLI tool that monitors HashiCorp Vault secret leases and sends alerts before expiration via configurable webhooks.

---

## Installation

```bash
go install github.com/youruser/vaultpulse@latest
```

Or download a pre-built binary from the [releases page](https://github.com/youruser/vaultpulse/releases).

---

## Usage

Set your Vault address and token, then run `vaultpulse` with a config file:

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.xxxxxxxxxxxxxxxx"

vaultpulse watch --config config.yaml
```

**Example `config.yaml`:**

```yaml
check_interval: 5m
alert_threshold: 24h
webhooks:
  - url: "https://hooks.slack.com/services/your/webhook/url"
    type: slack
  - url: "https://your-custom-endpoint.example.com/alert"
    type: generic
leases:
  - path: "secret/data/my-app/db-credentials"
  - path: "aws/creds/my-role"
```

**Available commands:**

```
vaultpulse watch      Start monitoring leases and sending alerts
vaultpulse list       List all tracked leases and their expiration times
vaultpulse version    Print the current version
```

---

## Configuration

| Field | Description | Default |
|---|---|---|
| `check_interval` | How often to poll Vault for lease status | `5m` |
| `alert_threshold` | How far before expiration to trigger an alert | `24h` |
| `webhooks` | List of webhook targets for notifications | `[]` |

---

## License

[MIT](LICENSE)