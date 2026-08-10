package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	namespace := flag.String("n", "", "Kubernetes namespace")
	podName := flag.String("p", "", "Kubernetes pod name")
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

	ctx := context.Background()

	// If a pod name is provided, show details for that pod.
	if *podName != "" {
		showPod(ctx, clientset, *namespace, *podName)
		return
	}

	// Otherwise list restarted pods.
	showRestartedPods(ctx, clientset, *namespace)
}

func exitCodeMeaning(exitCode int32, reason string) string {
	switch exitCode {
	case 0:
		return "Completed / Success"

	case 1:
		return "Application error"

	case 2:
		return "Misuse / invalid argument"

	case 125:
		return "Container runtime error"

	case 126:
		return "Permission denied / not executable"

	case 127:
		return "Command not found"

	case 137:
		if reason == "OOMKilled" {
			return "SIGKILL / OOMKilled"
		}
		return "SIGKILL"

	case 139:
		return "SIGSEGV / segmentation fault"

	case 143:
		return "SIGTERM / graceful termination"

	case 255:
		return "Specific error,cri,kublet,infra,node,unknown "

	default:
		return "Unknown"
	}
}

func showRestartedPods(
	ctx context.Context,
	clientset *kubernetes.Clientset,
	namespace string,
) {
	pods, err := clientset.CoreV1().
		Pods(namespace).
		List(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf(
		"%-12s %-45s %-10s %-20s %-10s %-30s\n",
		"NAMESPACE",
		"POD",
		"RESTARTS",
		"REASON",
		"EXIT CODE",
		"EXIT CODE MEANING",
	)

	fmt.Println("----------------------------------------------------------------------------------------------------------------")

	for _, pod := range pods.Items {
		var totalRestarts int32
		reason := "-"
		exitCode := int32(0)

		for _, container := range pod.Status.ContainerStatuses {
			totalRestarts += container.RestartCount

			if container.RestartCount > 0 &&
				container.LastTerminationState.Terminated != nil {

				termination := container.LastTerminationState.Terminated
				reason = termination.Reason
				exitCode = termination.ExitCode
			}
		}

		if totalRestarts == 0 {
			continue
		}

		meaning := exitCodeMeaning(exitCode, reason)

		fmt.Printf(
			"%-12s %-45s %-10d %-20s %-10d %-30s\n",
			pod.Namespace,
			pod.Name,
			totalRestarts,
			reason,
			exitCode,
			meaning,
		)
	}
}

func showPod(
	ctx context.Context,
	clientset *kubernetes.Clientset,
	namespace string,
	podName string,
) {
	// Get Pod
	pod, err := clientset.CoreV1().
		Pods(namespace).
		Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("NAMESPACE: %s\n", pod.Namespace)
	fmt.Printf("POD:       %s\n", pod.Name)
	fmt.Printf("STATUS:    %s\n", pod.Status.Phase)

	fmt.Println()

	fmt.Println("Restart Information:")
	fmt.Println("--------------------")

	var totalRestarts int32

	for _, container := range pod.Status.ContainerStatuses {
		totalRestarts += container.RestartCount
	}

	fmt.Printf("Total Restarts: %d\n", totalRestarts)

	for _, container := range pod.Status.ContainerStatuses {
		if container.RestartCount == 0 {
			continue
		}

		reason := "-"
		exitCode := int32(0)

		if container.LastTerminationState.Terminated != nil {
			termination := container.LastTerminationState.Terminated

			reason = termination.Reason
			exitCode = termination.ExitCode

			fmt.Printf(
				"Reason: %s, ExitCode: %d\n",
				reason,
				exitCode,
			)

			if !termination.FinishedAt.IsZero() {
				fmt.Printf(
					"Finished: %s\n",
					termination.FinishedAt.Format(time.RFC3339),
				)
			}
		}
	}

	fmt.Println()
	fmt.Println("Events:")
	fmt.Println("-------")

	showPodEvents(ctx, clientset, pod)
}

func showPodEvents(
	ctx context.Context,
	clientset *kubernetes.Clientset,
	pod *v1.Pod,
) {
	events, err := clientset.CoreV1().
		Events(pod.Namespace).
		List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Unable to get events: %v\n", err)
		return
	}

	var podEvents []v1.Event

	for _, event := range events.Items {
		if event.InvolvedObject.Kind != "Pod" {
			continue
		}

		if event.InvolvedObject.Name != pod.Name {
			continue
		}

		podEvents = append(podEvents, event)
	}

	sort.Slice(podEvents, func(i, j int) bool {
		return podEvents[i].LastTimestamp.Before(
			&podEvents[j].LastTimestamp,
		)
	})

	if len(podEvents) == 0 {
		fmt.Println("No events found.")
		return
	}

	fmt.Printf(
		"%-20s %-10s %-20s %s\n",
		"LAST SEEN",
		"TYPE",
		"REASON",
		"MESSAGE",
	)

	for _, event := range podEvents {
		lastSeen := event.LastTimestamp.Time

		if lastSeen.IsZero() && !event.EventTime.IsZero() {
			lastSeen = event.EventTime.Time
		}

		fmt.Printf(
			"%-20s %-10s %-20s %s\n",
			lastSeen.Format(time.RFC3339),
			event.Type,
			event.Reason,
			event.Message,
		)
	}
}
