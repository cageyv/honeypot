# honeypot

A minimal honeypot binary that exits with a configurable exit code. It can be used to mimic popular binaries in a container or on a local system to detect malicious activity.

## Usage

### Basic Usage

Run the binary with default exit code (222):

```bash
./honeypot
```

### Custom Exit Code via Flag

Specify a custom exit code using the `-honeypot-exit-code` flag:

```bash
./honeypot -honeypot-exit-code=201
```

### Custom Exit Code via Environment Variable

Set the `HONEYPOT_EXIT_CODE` environment variable:

```bash
HONEYPOT_EXIT_CODE=111 ./honeypot
```

## Security Use Cases

- **Container Protection**: Copy this binary into a scratch container as a decoy for popular binaries. If a malicious AI agent or script attempts to execute it, the container will exit immediately, alerting you to potential intrusion.
- **Local System Protection**: Place this binary in your local system to protect against unauthorized executions. If an attacker tries to run it, the system will exit with the configured code, helping you identify suspicious activity.
- **Honeypot for Attackers**: Use this binary as a honeypot to detect and analyze potential hacking attempts. By monitoring exit codes, you can gather insights into attack patterns.

## Development

### Prerequisites

- Go 1.x

### Building

Build the binary:

```bash
go build -o honeypot cmd/honeypot/main.go
```

### Testing

Run tests:

```bash
go test ./cmd/honeypot -v
```

### CI/CD

This project uses GitHub Actions for CI/CD. The workflow runs tests and builds the binary using GoReleaser. See `.github/workflows/ci.yml` for details.

### Contributing

1. Fork the repository.
2. Create a feature branch.
3. Commit your changes.
4. Push to the branch.
5. Create a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.