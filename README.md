# krestart

A lightweight Kubernetes CLI tool and `kubectl` plugin for identifying Pods with container restarts and investigating container termination reasons, exit codes, restart counts, and related Kubernetes events.

Use it as a standalone CLI or directly through `kubectl`:

```bash
restart-info -n production
kubectl restart-info -n production
```

## Features

- Find Pods with restarted containers across all namespaces or within a specific namespace.
- Inspect restart details for a specific Pod.
- Display container restart counts, termination reasons, and exit codes.
- Show exit code meanings to help diagnose container failures.
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

### Install using Krew

[Krew](https://krew.sigs.k8s.io/) is the plugin manager for `kubectl`.

If you already have Krew installed, install `restart-info` directly:

```bash
kubectl krew install restart-info
```

Verify the installation:

```bash
kubectl plugin list
```

Run the plugin:

```bash
kubectl restart-info
```

Upgrade to the latest published version:

```bash
kubectl krew upgrade restart-info
```

Uninstall:

```bash
kubectl krew uninstall restart-info
```

For Krew installation instructions, see the [official Krew documentation](https://krew.sigs.k8s.io/docs/user-guide/setup/install/).

### Install using a local Krew manifest

If you want to install the plugin directly from a manifest file before submitting it to the official Krew index:

```bash
kubectl krew install --manifest=./plugins/restart-info.yaml
```

Replace the manifest path with the location of your local manifest.

This method is useful for testing manifest changes before submitting them to the Krew index.

### Download a prebuilt release

Download the appropriate archive for your operating system and architecture from the [GitHub Releases](https://github.com/sdmoradi/krestart/releases) page.

Available release archives:

| Platform | Architecture | Archive |
|---|---|---|
| Linux | amd64 | `krestart_0.1.2_linux_amd64.tar.gz` |
| Linux | arm64 | `krestart_0.1.2_linux_arm64.tar.gz` |
| macOS | amd64 | `krestart_0.1.2_darwin_amd64.tar.gz` |
| macOS | arm64 | `krestart_0.1.2_darwin_arm64.tar.gz` |
| Windows | amd64 | `krestart_0.1.2_windows_amd64.zip` |

Verify the downloaded archive against `checksums.txt` from the same release before installing it.

Each archive contains:

- `restart-info` — standalone CLI executable.
- `kubectl-restart-info` — executable for the kubectl plugin interface.
- `LICENSE`
- `README.md`

### Install as a standalone CLI

On Linux or macOS, extract the downloaded archive and install the standalone binary:

```bash
tar -xzf krestart_0.1.2_linux_amd64.tar.gz

sudo install -m 0755 restart-info /usr/local/bin/restart-info
```

Replace the archive name with the one matching your platform and architecture.

Verify the installation:

```bash
restart-info --help
```

Run:

```bash
restart-info -n production
```

### Install manually as a kubectl plugin

Extract the release archive and install the plugin executable:

```bash
tar -xzf krestart_0.1.2_linux_amd64.tar.gz

mkdir -p "$HOME/.local/bin"
install -m 0755 kubectl-restart-info "$HOME/.local/bin/kubectl-restart-info"
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
kubectl restart-info -n production
```

**Windows:** Extract the Windows archive and place `kubectl-restart-info.exe` or `restart-info.exe` in a directory included in your `PATH`.

## Usage

### Help

Display available commands, flags, descriptions, and examples:

```bash
restart-info --help
```

Or through kubectl:

```bash
kubectl restart-info --help
```

### Find restarted Pods

Find restarted Pods across all namespaces:

```bash
kubectl restart-info
```

Find restarted Pods in a specific namespace:

```bash
kubectl restart-info -n production
```

### Inspect a specific Pod

Display detailed container restart information and related Kubernetes events:

```bash
kubectl restart-info -n production -p api-7d9f8c
```

### Use a custom kubeconfig

```bash
kubectl restart-info --kubeconfig /path/to/kubeconfig -n production
```

The standalone CLI accepts the same arguments:

```bash
restart-info -n production
restart-info -n production -p api-7d9f8c
restart-info --kubeconfig /path/to/kubeconfig -n production
```

## Kubernetes Configuration

`restart-info` uses the standard Kubernetes client configuration loading rules.

The configuration is loaded from:

- An explicitly supplied `--kubeconfig` path, when provided.
- The `KUBECONFIG` environment variable, when configured.
- The default kubeconfig location, typically `~/.kube/config`.

The current context in the selected configuration determines the cluster and credentials used by the tool.

For example:

```bash
kubectl config get-contexts
kubectl config use-context staging

kubectl restart-info -n production
```

You can also use a separate kubeconfig file without changing your default configuration:

```bash
kubectl restart-info --kubeconfig /path/to/kubeconfig -n production
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
go build -o restart-info ./cmd/krestart
```

Build the kubectl plugin:

```bash
go build -o kubectl-restart-info ./cmd/krestart
```

### Run Without Building

```bash
go run ./cmd/krestart
```

Run for a specific namespace:

```bash
go run ./cmd/krestart -n production
```

Display help:

```bash
go run ./cmd/krestart --help
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
├── plugins/
│   └── restart-info.yaml
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

Contributions, bug reports, and suggestions are welcome.

For code changes, run the tests and build checks before submitting a pull request.

## License

See [LICENSE](LICENSE) for license information.