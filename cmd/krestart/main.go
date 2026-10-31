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
	namespace := flag.String("n", "", "Kubernetes namespace")
	podName := flag.String("p", "", "Kubernetes pod name")
	kubeconfig := flag.String("kubeconfig", "", "Path to kubeconfig file")

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

		// Event conversion and sorting will be implemented in Phase 7.
		_ = events

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