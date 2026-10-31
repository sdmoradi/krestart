# krestart

A lightweight Kubernetes CLI tool to identify Pods with container restarts and inspect restart counts, termination reasons, exit codes, and related Kubernetes events.

## Features

- Find restarted Pods across all namespaces or within a specific namespace.
- Inspect restart details for a specific Pod.
- Display container termination reasons, exit codes, and exit code meanings.
- View Kubernetes events associated with a specific Pod.
- Use standard Kubernetes kubeconfig loading and context selection.
- Specify a custom kubeconfig file when needed.

## Requirements

- Go 1.22 or later.
- Access to a Kubernetes cluster.
- A valid Kubernetes configuration, typically `~/.kube/config`.
- Appropriate Kubernetes RBAC permissions to list Pods and read Pod events.

## Build

Clone the repository and build the CLI:

```bash
git clone https://github.com/sdmoradi/krestart.git
cd krestart

go build -o krestart ./cmd/krestart
```

Run the following command to verify the binary:

```bash
./krestart --help
```

## Usage

### Find restarted Pods across all namespaces

```bash
./krestart
```

### Find restarted Pods in a specific namespace

```bash
./krestart -n production
```

### Inspect a specific Pod

```bash
./krestart -n production -p pod-name
```

This displays the Pod's restart details and related Kubernetes events.

### Use a custom kubeconfig file

```bash
./krestart --kubeconfig /path/to/kubeconfig
```

Combine it with a namespace and Pod name if needed:

```bash
./krestart --kubeconfig /path/to/kubeconfig -n production -p pod-name
```

## Kubernetes Configuration

`krestart` uses the standard Kubernetes client configuration loading rules.

Configuration is loaded from:

- An explicitly supplied `--kubeconfig` path, when provided.
- The `KUBECONFIG` environment variable, when configured.
- The default kubeconfig location, typically `~/.kube/config`.

The current context in the selected configuration determines which cluster and credentials are used.

For example, to use a specific kubeconfig:

```bash
./krestart --kubeconfig ~/.kube/config -n production
```

To use a different context, select it in your kubeconfig or use the standard `kubectl config use-context` command.

## Install

After building the binary, install it into `/usr/local/bin`:

```bash
sudo install -m 0755 krestart /usr/local/bin/krestart
```

You can then run it from any directory:

```bash
krestart
```

Or inspect a specific namespace:

```bash
krestart -n production
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

**Note:** Exit code `137` indicates termination by `SIGKILL`; it does not, by itself, prove that a container was killed because of an out-of-memory condition. Check the termination reason and Pod events for additional context.

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

Run against a specific namespace:

```bash
go run ./cmd/krestart -n production
```

Use a custom kubeconfig:

```bash
go run ./cmd/krestart --kubeconfig /path/to/kubeconfig
```

### Format Code

Format the Go source files:

```bash
gofmt -w cmd/ internal/
```

### Run Tests

```bash
go test ./...
```

### Build

```bash
go build ./...
```

## License

See [LICENSE](LICENSE) for license information.