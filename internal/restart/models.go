package restart

import "time"

type ContainerRestart struct {
	ContainerName string
	RestartCount  int32
	Reason        string
	ExitCode      int32
	FinishedAt    time.Time
}

type PodRestart struct {
	Namespace      string
	PodName        string
	TotalRestarts  int32
	ContainerItems []ContainerRestart
}

type PodEvent struct {
	LastSeen time.Time
	Type     string
	Reason   string
	Message  string
}