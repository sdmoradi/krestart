package output

import (
	"fmt"
	"time"

	"github.com/sdmoradi/krestart/internal/restart"
)

func PrintRestartedPods(pods []restart.PodRestart) {
	fmt.Printf(
		"%-12s %-45s %-12s %-20s %-10s %-30s\n",
		"NAMESPACE",
		"POD",
		"CONTAINER",
		"REASON",
		"EXIT CODE",
		"EXIT CODE MEANING",
	)

	fmt.Println("----------------------------------------------------------------------------------------------------------------")

	for _, pod := range pods {
		for _, container := range pod.ContainerItems {
			fmt.Printf(
				"%-12s %-45s %-12s %-20s %-10d %-30s\n",
				pod.Namespace,
				pod.PodName,
				container.ContainerName,
				container.Reason,
				container.ExitCode,
				restart.ExitCodeMeaning(
					container.ExitCode,
					container.Reason,
				),
			)
		}
	}
}

func PrintPodDetails(pod restart.PodRestart) {
	fmt.Printf("NAMESPACE: %s\n", pod.Namespace)
	fmt.Printf("POD:       %s\n", pod.PodName)
	fmt.Printf("Total Restarts: %d\n", pod.TotalRestarts)

	fmt.Println()
	fmt.Println("Container Restart Information:")
	fmt.Println("------------------------------")

	for _, container := range pod.ContainerItems {
		fmt.Printf("Container: %s\n", container.ContainerName)
		fmt.Printf("Restarts:  %d\n", container.RestartCount)
		fmt.Printf("Reason:    %s\n", container.Reason)
		fmt.Printf("Exit Code: %d\n", container.ExitCode)

		if !container.FinishedAt.IsZero() {
			fmt.Printf(
				"Finished:  %s\n",
				container.FinishedAt.Format(time.RFC3339),
			)
		}

		fmt.Println()
	}
}

func PrintPodEvents(events []restart.PodEvent) {
	if len(events) == 0 {
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

	for _, event := range events {
		fmt.Printf(
			"%-20s %-10s %-20s %s\n",
			event.LastSeen.Format(time.RFC3339),
			event.Type,
			event.Reason,
			event.Message,
		)
	}
}