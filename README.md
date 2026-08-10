# krestart

A lightweight Kubernetes CLI tool for identifying restarted containers
and showing their termination reasons and exit codes.

## Features

- List restarted containers across all namespaces
- Show container restart counts
- Show termination reasons
- Show exit codes
- Uses the local Kubernetes kubeconfig
- Lightweight and dependency-minimal

## Requirements

- Go 1.XX+
- Access to a Kubernetes cluster
- Valid kubeconfig

## Installation

Clone the repository:

```bash
git clone https://github.com/sdmoradi/krestart.git
cd krestart