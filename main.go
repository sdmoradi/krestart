package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// Command-line flags
	namespace := flag.String("n", "", "Kubernetes namespace")
	flag.Parse()

	// kubeconfig location
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	kubeconfig := filepath.Join(home, ".kube", "config")

	// Build Kubernetes config
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}

	// Create Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// List pods
	pods, err := clientset.CoreV1().
		Pods(*namespace).
		List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf(
		"%-12s %-45s %-10s %-20s %-10s\n",
		"NAMESPACE",
		"POD",
		"RESTARTS",
		"REASON",
		"EXIT CODE",
	)

	fmt.Println("-----------------------------------------------------------------------------------------------")

	for _, pod := range pods.Items {
		var totalRestarts int32
		reason := "-"
		var exitCode int32

		for _, container := range pod.Status.ContainerStatuses {
			totalRestarts += container.RestartCount

			if container.LastTerminationState.Terminated != nil {
				reason = container.LastTerminationState.Terminated.Reason
				exitCode = container.LastTerminationState.Terminated.ExitCode
			} else if container.State.Waiting != nil {
				reason = container.State.Waiting.Reason
			}
		}

		if totalRestarts == 0 {
			continue
		}

		fmt.Printf(
			"%-12s %-45s %-10d %-20s %-10d\n",
			pod.Namespace,
			pod.Name,
			totalRestarts,
			reason,
			exitCode,
		)
	}
}
