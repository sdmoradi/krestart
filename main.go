package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	// kubeconfig location
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	kubeconfig := filepath.Join(home, ".kube", "config")

	// Build config
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}

	// Create client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// List pods in all namespaces
	pods, err := clientset.CoreV1().
		Pods("").
		List(context.Background(), metav1.ListOptions{})

	if err != nil {
		panic(err)
	}

	fmt.Printf("%-20s %-45s %-10s\n", "NAMESPACE", "POD", "RESTARTS")
	fmt.Println("---------------------------------------------------------------------------------------")

	for _, pod := range pods.Items {
		var totalRestarts int32

		for _, container := range pod.Status.ContainerStatuses {

			totalRestarts += container.RestartCount

			reason := "-"
			exitCode := int32(0)

			if container.LastTerminationState.Terminated != nil {
				reason = container.LastTerminationState.Terminated.Reason
				exitCode = container.LastTerminationState.Terminated.ExitCode
			} else if container.State.Waiting != nil {
				reason = container.State.Waiting.Reason
			}

			if container.RestartCount > 0 {
				fmt.Printf(
					"Container: %-20s Restarts:%2d  Reason:%-20s ExitCode:%d\n",
					container.Name,
					container.RestartCount,
					reason,
					exitCode,
				)
			}
		}
	}
}
