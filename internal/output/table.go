package output

import (
	"fmt"
	"github.com/sdmoradi/krestart/internal/restart"
	"os"
	"text/tabwriter"
	"time"
)

func PrintRestartedPods(pods []restart.PodRestart) {
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)

	fmt.Fprintln(w, "NAMESPACE\tPOD\tCONTAINER\tRESTARTS\tREASON\tEXIT CODE\tEXIT CODE MEANING")

	for _, pod := range pods {
		for _, container := range pod.ContainerItems {
			fmt.Fprintf(
				w,
				"%s\t%s\t%s\t%d\t%s\t%d\t%s\n",
				truncate(pod.Namespace, 15),
				truncate(pod.PodName, 35),
				truncate(container.ContainerName, 22),
				container.RestartCount,
				truncate(container.Reason, 15),
				container.ExitCode,
				truncate(
					restart.ExitCodeMeaning(
						container.ExitCode,
						container.Reason,
					),
					28,
				),
			)
		}
	}

	if err := w.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to print table: %v\n", err)
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
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
