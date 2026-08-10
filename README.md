# krestart

A lightweight Kubernetes CLI tool to find restarted Pods and show their restart count, reason, and exit code.

## Requirements

* Go 1.XX+
* Kubernetes cluster
* Valid `~/.kube/config`

## Build

```bash
git clone https://github.com/YOUR_USERNAME/krestart.git
cd krestart

go build -o krestart .
```

## Run

Check restarted Pods across all namespaces:

```bash
./krestart
```

Check a specific namespace:

```bash
./krestart -n production
```

Check a specific Pod:

```bash
./krestart -n production -p pod-name
```

## Install

Install the binary globally:

```bash
sudo install -m 0755 krestart /usr/local/bin/krestart
```

Then run:

```bash
krestart
```

Or for a specific namespace:

```bash
krestart -n production
```

## Example

```text
NAMESPACE    POD                                           RESTARTS   REASON               EXIT CODE
---------------------------------------------------------------------------------------------------
stage        cms-adminapi-6ddcf7d544-fttb                 1          Error                137
stage        cms-adminapi-6ddcf7d544-q4vvt                 1          Error                137
stage        frontend-admin-5c5444ddf5-257c2               3          Error                  1
```

## Development

Run directly:

```bash
go run .
```

Run for a specific namespace:

```bash
go run . -n production
```

Format:

```bash
gofmt -w main.go
```