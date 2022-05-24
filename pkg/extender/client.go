package extender

import (
	"fmt"
	lvmclient "github.com/openebs/lvm-localpv/pkg/generated/clientset/internalclientset"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// getK8sClient returns K8s clientset by taking kubeconfig as an argument
func getK8sClient() (*kubernetes.Clientset, error) {
	// use the current context in kubeconfig
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, errors.Wrap(err, "Could not get in cluster config")
	}
	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "Could not get new config")
	}
	return clientset, nil
}

func getLVMClient() (*lvmclient.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, errors.Wrap(err, "Could not get in cluster config")
	}
	client, err := lvmclient.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("could not get new config: %v", err)
	}
	return client, nil
}
