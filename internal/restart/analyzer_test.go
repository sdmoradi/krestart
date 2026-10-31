package restart

import (
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestExitCodeMeaning(t *testing.T) {
	tests := []struct {
		name     string
		exitCode int32
		reason   string
		expected string
	}{
		{"success", 0, "", "Completed / Success"},
		{"application error", 1, "", "Application error"},
		{"invalid argument", 2, "", "Misuse / invalid argument"},
		{"runtime error", 125, "", "Container runtime error"},
		{"permission denied", 126, "", "Permission denied / not executable"},
		{"command not found", 127, "", "Command not found"},
		{"killed", 137, "", "SIGKILL"},
		{"segmentation fault", 139, "", "SIGSEGV / segmentation fault"},
		{"terminated", 143, "", "SIGTERM / graceful termination"},
		{"unknown", 255, "", "Unknown"},
		{"OOMKilled overrides exit code", 137, "OOMKilled", "OOMKilled"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExitCodeMeaning(tt.exitCode, tt.reason)
			if got != tt.expected {
				t.Errorf(
					"ExitCodeMeaning(%d, %q) = %q; want %q",
					tt.exitCode,
					tt.reason,
					got,
					tt.expected,
				)
			}
		})
	}
}

func TestAnalyzePod(t *testing.T) {
	finishedAt := time.Date(2026, 10, 31, 10, 0, 0, 0, time.UTC)

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wallet-adminapi-123",
			Namespace: "production",
		},
		Status: corev1.PodStatus{
			ContainerStatuses: []corev1.ContainerStatus{
				{
					Name:         "api",
					RestartCount: 3,
					LastTerminationState: corev1.ContainerState{
						Terminated: &corev1.ContainerStateTerminated{
							Reason:     "OOMKilled",
							ExitCode:   137,
							FinishedAt: metav1.NewTime(finishedAt),
						},
					},
				},
				{
					Name:         "sidecar",
					RestartCount: 0,
				},
				{
					Name:         "worker",
					RestartCount: 2,
					LastTerminationState: corev1.ContainerState{
						Terminated: &corev1.ContainerStateTerminated{
							Reason:   "Error",
							ExitCode: 1,
						},
					},
				},
			},
		},
	}

	got := AnalyzePod(pod)

	if got.Namespace != "production" {
		t.Errorf("Namespace = %q; want %q", got.Namespace, "production")
	}

	if got.PodName != "wallet-adminapi-123" {
		t.Errorf("PodName = %q; want %q", got.PodName, "wallet-adminapi-123")
	}

	if got.TotalRestarts != 5 {
		t.Errorf("TotalRestarts = %d; want 5", got.TotalRestarts)
	}

	if len(got.ContainerItems) != 2 {
		t.Fatalf("len(ContainerItems) = %d; want 2", len(got.ContainerItems))
	}

	api := got.ContainerItems[0]
	if api.ContainerName != "api" ||
		api.RestartCount != 3 ||
		api.Reason != "OOMKilled" ||
		api.ExitCode != 137 ||
		!api.FinishedAt.Equal(finishedAt) {
		t.Errorf("unexpected api container result: %+v", api)
	}

	worker := got.ContainerItems[1]
	if worker.ContainerName != "worker" ||
		worker.RestartCount != 2 ||
		worker.Reason != "Error" ||
		worker.ExitCode != 1 {
		t.Errorf("unexpected worker container result: %+v", worker)
	}
}

func TestAnalyzePods(t *testing.T) {
	pods := []corev1.Pod{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "restarted-pod",
				Namespace: "production",
			},
			Status: corev1.PodStatus{
				ContainerStatuses: []corev1.ContainerStatus{
					{Name: "api", RestartCount: 2},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "healthy-pod",
				Namespace: "production",
			},
			Status: corev1.PodStatus{
				ContainerStatuses: []corev1.ContainerStatus{
					{Name: "api", RestartCount: 0},
				},
			},
		},
	}

	got := AnalyzePods(pods)

	if len(got) != 1 {
		t.Fatalf("len(results) = %d; want 1", len(got))
	}

	if got[0].PodName != "restarted-pod" {
		t.Errorf("PodName = %q; want %q", got[0].PodName, "restarted-pod")
	}

	if got[0].TotalRestarts != 2 {
		t.Errorf("TotalRestarts = %d; want 2", got[0].TotalRestarts)
	}
}

func TestAnalyzeEvents(t *testing.T) {
	base := time.Date(2026, 10, 31, 10, 0, 0, 0, time.UTC)

	seriesTime := base.Add(3 * time.Minute)
	eventTime := base.Add(2 * time.Minute)
	lastTimestamp := base.Add(time.Minute)
	creationTime := base

	events := []corev1.Event{
		{
			Reason: "SeriesEvent",
			Series: &corev1.EventSeries{
				LastObservedTime: metav1.NewMicroTime(seriesTime),
			},
		},
		{
			Reason:    "EventTime",
			EventTime: metav1.NewMicroTime(eventTime),
		},
		{
			Reason:        "LastTimestamp",
			LastTimestamp: metav1.NewTime(lastTimestamp),
		},
		{
			Reason: "CreationTimestamp",
			ObjectMeta: metav1.ObjectMeta{
				CreationTimestamp: metav1.NewTime(creationTime),
			},
		},
	}

	got := AnalyzeEvents(events)

	if len(got) != 4 {
		t.Fatalf("len(results) = %d; want 4", len(got))
	}

	expected := []struct {
		reason string
		time   time.Time
	}{
		{"CreationTimestamp", creationTime},
		{"LastTimestamp", lastTimestamp},
		{"EventTime", eventTime},
		{"SeriesEvent", seriesTime},
	}

	for i, want := range expected {
		if got[i].Reason != want.reason {
			t.Errorf(
				"results[%d].Reason = %q; want %q",
				i,
				got[i].Reason,
				want.reason,
			)
		}

		if !got[i].LastSeen.Equal(want.time) {
			t.Errorf(
				"results[%d].LastSeen = %v; want %v",
				i,
				got[i].LastSeen,
				want.time,
			)
		}
	}
}

func TestAnalyzeEventsTimestampFallback(t *testing.T) {
	base := time.Date(2026, 10, 31, 10, 0, 0, 0, time.UTC)

	events := []corev1.Event{
		{
			Reason: "CreationTimestamp",
			ObjectMeta: metav1.ObjectMeta{
				CreationTimestamp: metav1.NewTime(base),
			},
		},
		{
			Reason: "NoTimestamp",
		},
	}

	got := AnalyzeEvents(events)

	if len(got) != 2 {
		t.Fatalf("len(results) = %d; want 2", len(got))
	}

	// Events are sorted by LastSeen, so zero timestamps come first.
	if got[0].Reason != "NoTimestamp" || !got[0].LastSeen.IsZero() {
		t.Errorf("unexpected zero timestamp result: %+v", got[0])
	}

	if got[1].Reason != "CreationTimestamp" || !got[1].LastSeen.Equal(base) {
		t.Errorf("unexpected creation timestamp result: %+v", got[1])
	}
}
