# portwatch

A lightweight CLI tool that monitors open ports and alerts on unexpected changes.

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

## Usage

Start monitoring with default settings (scans every 60 seconds):

```bash
portwatch start
```

Specify a custom scan interval and alert on any changes:

```bash
portwatch start --interval 30 --notify
```

Run a one-time scan and display all open ports:

```bash
portwatch scan
```

Example output:

```
[+] Monitoring started — baseline captured (12 open ports)
[!] ALERT: New port detected — 8080/tcp (2024-11-01 14:32:07)
[!] ALERT: Port closed — 3306/tcp (2024-11-01 14:35:22)
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--interval` | `60` | Scan interval in seconds |
| `--notify` | `false` | Enable desktop notifications |
| `--log` | `""` | Path to write alert log file |
| `--ports` | `all` | Comma-separated port range to watch |

## License

MIT © [yourusername](https://github.com/yourusername)