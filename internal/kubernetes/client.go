package kubernetes

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// NewClient creates a Kubernetes client using standard kubeconfig loading rules.
func NewClient(kubeconfigPath string) (*kubernetes.Clientset, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()

	// Use an explicitly supplied kubeconfig path when provided.
	if kubeconfigPath != "" {
		loadingRules.ExplicitPath = kubeconfigPath
	}

	configOverrides := &clientcmd.ConfigOverrides{}

	configLoader := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		loadingRules,
		configOverrides,
	)

	config, err := configLoader.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("load Kubernetes configuration: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("create Kubernetes client: %w", err)
	}

	return clientset, nil
}
