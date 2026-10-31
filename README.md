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
- Support standard Kubernetes kubeconfig loading and context selection through the configuration file.
- Specify a custom kubeconfig file with `--kubeconfig`.
- Run as a standalone CLI or a `kubectl` plugin.

## Requirements

- Go 1.22 or later to build from source.
- `kubectl` installed to use the plugin interface.
- Access to a Kubernetes cluster.
- A valid Kubernetes configuration, typically `~/.kube/config`.
- Appropriate Kubernetes RBAC permissions to list Pods and read Pod events.

## Build

Clone the repository and build the standalone binary:

```bash
git clone https://github.com/sdmoradi/krestart.git
cd krestart

GOTOOLCHAIN=local go build -o krestart ./cmd/krestart
```

To build the `kubectl` plugin binary:

```bash
GOTOOLCHAIN=local go build -o kubectl-krestart ./cmd/krestart
```

## Usage

### Standalone CLI

Find restarted Pods across all namespaces:

```bash
./krestart
```

Find restarted Pods in a specific namespace:

```bash
./krestart -n production
```

Inspect a specific Pod:

```bash
./krestart -n production -p pod-name
```

Use a custom kubeconfig:

```bash
./krestart --kubeconfig /path/to/kubeconfig -n production
```

### kubectl Plugin

After installing the `kubectl-krestart` executable, invoke it through `kubectl`.

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

The plugin accepts the same arguments as the standalone CLI.

## Install

### Install as a kubectl plugin

Build the plugin binary:

```bash
go build -o kubectl-krestart ./cmd/krestart
```

Create a directory for user-installed executables and install the binary:

```bash
mkdir -p "$HOME/.local/bin"
install -m 0755 kubectl-krestart "$HOME/.local/bin/kubectl-krestart"
```

Add the directory to your `PATH` if it is not already there:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

To persist this setting for Bash:

```bash
grep -qxF 'export PATH="$HOME/.local/bin:$PATH"' ~/.bashrc \
  || echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
```

Verify that the plugin is discoverable:

```bash
kubectl plugin list
```

Run it:

```bash
kubectl krestart
```

### Install as a standalone CLI

Build the standalone binary:

```bash
go build -o krestart ./cmd/krestart
```

Install it system-wide:

```bash
sudo install -m 0755 krestart /usr/local/bin/krestart
```

Run it from any directory:

```bash
krestart -n production
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

**Note:** Selecting a context through `kubectl config use-context` changes the current context in the kubeconfig. If you need to target a context without changing the current selection, explicit context support should be implemented and tested in `krestart`.

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
│   │   ├── models.go
│   │   └── analyzer.go
│   └── output/
│       └── table.go
├── go.mod
├── go.sum
├── README.md
└── LICENSE
```

### Run Without Building

```bash
go run ./cmd/krestart
```

Run for a specific namespace:

```bash
go run ./cmd/krestart -n production
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