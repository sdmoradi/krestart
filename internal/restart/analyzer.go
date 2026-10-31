package restart

import (

	"time"
	"sort"
	corev1 "k8s.io/api/core/v1"
)

func ExitCodeMeaning(exitCode int32, reason string) string {
	if reason == "OOMKilled" {
		return "OOMKilled"
	}

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
		return "SIGKILL"
	case 139:
		return "SIGSEGV / segmentation fault"
	case 143:
		return "SIGTERM / graceful termination"
	default:
		return "Unknown"
	}
}

func AnalyzePod(pod *corev1.Pod) PodRestart {
	result := PodRestart{
		Namespace: pod.Namespace,
		PodName:   pod.Name,
	}

	for _, container := range pod.Status.ContainerStatuses {
		result.TotalRestarts += container.RestartCount

		if container.RestartCount == 0 {
			continue
		}

		item := ContainerRestart{
			ContainerName: container.Name,
			RestartCount:  container.RestartCount,
			Reason:        "-",
		}

		if termination := container.LastTerminationState.Terminated; termination != nil {
			item.Reason = termination.Reason
			item.ExitCode = termination.ExitCode
			item.FinishedAt = termination.FinishedAt.Time
		}

		result.ContainerItems = append(result.ContainerItems, item)
	}

	return result
}

func AnalyzePods(pods []corev1.Pod) []PodRestart {
	var results []PodRestart

	for i := range pods {
		result := AnalyzePod(&pods[i])

		if result.TotalRestarts == 0 {
			continue
		}

		results = append(results, result)
	}

	return results
}

func AnalyzeEvents(events []corev1.Event) []PodEvent {
	results := make([]PodEvent, 0, len(events))

	for _, event := range events {
		var lastSeen time.Time

		if event.Series != nil {
			lastSeen = event.Series.LastObservedTime.Time
		}

		if lastSeen.IsZero() {
			lastSeen = event.EventTime.Time
		}

		if lastSeen.IsZero() {
			lastSeen = event.LastTimestamp.Time
		}

		if lastSeen.IsZero() {
			lastSeen = event.CreationTimestamp.Time
		}

		results = append(results, PodEvent{
			LastSeen: lastSeen,
			Type:     event.Type,
			Reason:   event.Reason,
			Message:  event.Message,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].LastSeen.Before(results[j].LastSeen)
	})

	return results
}