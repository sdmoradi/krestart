# krestart

A lightweight Kubernetes CLI tool to find restarted containers and show their restart reason and exit code.

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

```bash
./krestart
```

Or install it globally:

```bash
sudo install -m 0755 krestart /usr/local/bin/krestart
```

Then run:

```bash
krestart -n production
```

## Example

```text
Container: api       Restarts:12  Reason:OOMKilled  ExitCode:137
Container: worker    Restarts: 4  Reason:Error      ExitCode:1
```

## Development

Run directly:

```bash
go run .
```


Format:

```bash
gofmt -w main.go
```