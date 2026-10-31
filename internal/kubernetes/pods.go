package kubernetes

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodClient struct {
	client kubernetes.Interface
}

func NewPodClient(client kubernetes.Interface) *PodClient {
	return &PodClient{
		client: client,
	}
}

func (p *PodClient) ListPods(
	ctx context.Context,
	namespace string,
) ([]corev1.Pod, error) {
	pods, err := p.client.CoreV1().
		Pods(namespace).
		List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return pods.Items, nil
}

func (p *PodClient) GetPod(
	ctx context.Context,
	namespace string,
	name string,
) (*corev1.Pod, error) {
	return p.client.CoreV1().
		Pods(namespace).
		Get(ctx, name, metav1.GetOptions{})
}

func (p *PodClient) ListPodEvents(
	ctx context.Context,
	pod *corev1.Pod,
) ([]corev1.Event, error) {
	events, err := p.client.CoreV1().
		Events(pod.Namespace).
		List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	var podEvents []corev1.Event

	for _, event := range events.Items {
		if event.InvolvedObject.Kind != "Pod" {
			continue
		}

		if event.InvolvedObject.Name != pod.Name {
			continue
		}

		if event.InvolvedObject.UID != pod.UID {
			continue
		}

		podEvents = append(podEvents, event)
	}

	return podEvents, nil
}