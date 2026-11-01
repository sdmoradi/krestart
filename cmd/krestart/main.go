package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	k8s "github.com/sdmoradi/krestart/internal/kubernetes"
	"github.com/sdmoradi/krestart/internal/output"
	"github.com/sdmoradi/krestart/internal/restart"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	namespace := flag.String("n", "", "Kubernetes namespace (default: all namespaces)")
	podName := flag.String("p", "", "Kubernetes Pod name")
	kubeconfig := flag.String("kubeconfig", "", "Path to kubeconfig file")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `restart-info - Kubernetes Container Restart Inspector

Inspect Kubernetes Pods, identify container restarts, and investigate
termination reasons, exit codes, and related events.

USAGE:
  kubectl restart-info [flags]

FLAGS:
  -n string
        Kubernetes namespace (default: all namespaces)

  -p string
        Kubernetes Pod name (use with -n to identify the Pod's namespace)

  -kubeconfig string
        Path to kubeconfig file (default: use standard kubeconfig loading rules)

  -h, -help
        Display this help information

EXAMPLES:
  # Inspect restarted containers across all namespaces
  kubectl restart-info

  # Inspect restarted containers in a specific namespace
  kubectl restart-info -n production

  # Inspect a specific Pod and its related events
  kubectl restart-info -n production -p api-7d9f8c

  # Use a specific kubeconfig file
  kubectl restart-info -kubeconfig ~/.kube/production

WHAT YOU GET:
  - Container restart counts
  - Last termination reason and exit code
  - Related Kubernetes events

`)
	}

	flag.Parse()

	clientset, err := k8s.NewClient(*kubeconfig)
	if err != nil {
		return err
	}

	podClient := k8s.NewPodClient(clientset)
	ctx := context.Background()

	if *podName != "" {
		pod, err := podClient.GetPod(ctx, *namespace, *podName)
		if err != nil {
			return fmt.Errorf("get pod %q: %w", *podName, err)
		}

		result := restart.AnalyzePod(pod)
		output.PrintPodDetails(result)

		fmt.Println()
		fmt.Println("Events:")
		fmt.Println("-------")

		events, err := podClient.ListPodEvents(ctx, pod)
		if err != nil {
			return fmt.Errorf("list pod events: %w", err)
		}

		eventResults := restart.AnalyzeEvents(events)
		output.PrintPodEvents(eventResults)

		return nil
	}

	pods, err := podClient.ListPods(ctx, *namespace)
	if err != nil {
		return fmt.Errorf("list pods: %w", err)
	}

	results := restart.AnalyzePods(pods)
	output.PrintRestartedPods(results)

	return nil
}
