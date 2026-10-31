# krestart

A lightweight Kubernetes CLI tool and `kubectl` plugin for identifying Pods with container restarts and inspecting their restart counts, termination reasons, exit codes, and related Kubernetes events.

Use it as a standalone CLI or directly through `kubectl`:

```bash
krestart -n production
kubectl krestart -n production
```

## Features

- Find restarted Pods across all namespaces or within a specific namespace.
- Inspect restart details for a specific Pod.
- Display container termination reasons, exit codes, and exit code meanings.
- View Kubernetes events associated with a specific Pod.
- Use standard Kubernetes kubeconfig loading and current-context selection.
- Specify a custom kubeconfig file with `--kubeconfig`.
- Run as a standalone CLI or a `kubectl` plugin.

## Requirements

- `kubectl` installed to use the plugin interface.
- Access to a Kubernetes cluster.
- A valid Kubernetes configuration, typically `~/.kube/config`.
- Appropriate Kubernetes RBAC permissions to list Pods and read Pod events.

Go 1.22 or later is only required when building from source.

## Installation

### Download a prebuilt release

Download the appropriate archive for your operating system and architecture from the [GitHub Releases](https://github.com/sdmoradi/krestart/releases) page.

Available release archives:

| Platform | Architecture | Archive |
|---|---|---|
| Linux | amd64 | `krestart_0.1.0_Linux_x86_64.tar.gz` |
| Linux | arm64 | `krestart_0.1.0_Linux_arm64.tar.gz` |
| macOS | amd64 | `krestart_0.1.0_Darwin_x86_64.tar.gz` |
| macOS | arm64 | `krestart_0.1.0_Darwin_arm64.tar.gz` |
| Windows | amd64 | `krestart_0.1.0_Windows_x86_64.zip` |

Verify the downloaded archive against `checksums.txt` from the same release before installing it.

Each archive contains:

- `krestart` — standalone CLI executable.
- `kubectl-krestart` — executable for the kubectl plugin interface.
- `LICENSE`
- `README.md`

### Install as a standalone CLI

On Linux or macOS, extract the downloaded archive and install the standalone binary:

```bash
tar -xzf krestart_0.1.0_Linux_x86_64.tar.gz

sudo install -m 0755 krestart /usr/local/bin/krestart
```

Replace the archive name with the one matching your platform and architecture.

Verify the installation:

```bash
krestart --help
```

Run:

```bash
krestart -n production
```

### Install as a kubectl plugin

Extract the release archive and install the plugin executable:

```bash
tar -xzf krestart_0.1.0_Linux_x86_64.tar.gz

mkdir -p "$HOME/.local/bin"
install -m 0755 kubectl-krestart "$HOME/.local/bin/kubectl-krestart"
```

Ensure `$HOME/.local/bin` is included in your `PATH`:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

To persist this setting for Bash:

```bash
grep -qxF 'export PATH="$HOME/.local/bin:$PATH"' ~/.bashrc \
  || echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
```

Verify that kubectl discovers the plugin:

```bash
kubectl plugin list
```

Run it:

```bash
kubectl krestart -n production
```

The plugin accepts the same arguments as the standalone CLI.

**Windows:** Extract the Windows archive and place `kubectl-krestart.exe` or `krestart.exe` in a directory included in your `PATH`.

## Usage

### Standalone CLI

Find restarted Pods across all namespaces:

```bash
krestart
```

Find restarted Pods in a specific namespace:

```bash
krestart -n production
```

Inspect a specific Pod:

```bash
krestart -n production -p pod-name
```

Use a custom kubeconfig:

```bash
krestart --kubeconfig /path/to/kubeconfig -n production
```

### kubectl Plugin

Find restarted Pods across all namespaces:

```bash
kubectl krestart
```

Find restarted Pods in a specific namespace:

```bash
kubectl krestart -n production
```

Inspect a specific Pod:

```bash
kubectl krestart -n production -p pod-name
```

Use a custom kubeconfig:

```bash
kubectl krestart --kubeconfig /path/to/kubeconfig -n production
```

## Kubernetes Configuration

`krestart` uses the standard Kubernetes client configuration loading rules.

The configuration is loaded from:

- An explicitly supplied `--kubeconfig` path, when provided.
- The `KUBECONFIG` environment variable, when configured.
- The default kubeconfig location, typically `~/.kube/config`.

The current context in the selected configuration determines the cluster and credentials used by the tool.

For example:

```bash
kubectl config get-contexts
kubectl config use-context staging

kubectl krestart -n production
```

You can also use a separate kubeconfig file without changing your default configuration:

```bash
kubectl krestart --kubeconfig /path/to/kubeconfig -n production
```

## Example Output

The following is an illustrative example of restarted Pod output:

```text
NAMESPACE    POD                                           RESTARTS   REASON               EXIT CODE  EXIT CODE MEANING
----------------------------------------------------------------------------------------------------------------
stage        cms-adminapi-6ddcf7d544-2j8kc                 1          Error                137        SIGKILL
stage        cms-adminapi-6ddcf7d544-fttb8                 2          Error                137        SIGKILL
stage        cms-adminapi-6ddcf7d544-q4vvt                 1          Error                137        SIGKILL
stage        frontend-admin-5c5444ddf5-257c2               5          Error                1          Application error
stage        frontend-admin-5c5444ddf5-5lr42               4          Error                1          Application error
stage        frontend-admin-5c5444ddf5-8hqpz               6          Error                1          Application error
```

**Note:** Exit code `137` indicates termination by `SIGKILL`; it does not, by itself, prove an out-of-memory condition. Check the container termination reason and Pod events for additional context.

## Development

### Requirements

- Go 1.22 or later.
- Git.

### Build from Source

Clone the repository:

```bash
git clone https://github.com/sdmoradi/krestart.git
cd krestart
```

Build the standalone CLI:

```bash
go build -o krestart ./cmd/krestart
```

Build the kubectl plugin:

```bash
go build -o kubectl-krestart ./cmd/krestart
```

### Run Without Building

```bash
go run ./cmd/krestart
```

Run for a specific namespace:

```bash
go run ./cmd/krestart -n production
```

### Project Structure

```text
krestart/
├── cmd/
│   └── krestart/
│       └── main.go
├── internal/
│   ├── kubernetes/
│   │   ├── client.go
│   │   └── pods.go
│   ├── restart/
│   │   ├── analyzer.go
│   │   └── analyzer_test.go
│   └── output/
│       └── table.go
├── go.mod
├── go.sum
├── README.md
└── LICENSE
```

### Format Code

```bash
gofmt -w cmd/ internal/
```

### Run Tests

```bash
go test ./...
```

### Build All Packages

```bash
go build ./...
```

## Contributing

Contributions, bug reports, and suggestions are welcome. For code changes, run the tests and build checks before submitting a pull request.

## License

See [LICENSE](LICENSE) for license information.