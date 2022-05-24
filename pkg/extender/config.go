package extender

import (
	"github.com/openebs/openebsctl/pkg/util"
	"k8s.io/client-go/rest"
	"log"
	"os"
	"path/filepath"
)

// GetOutofClusterKubeConfig creates returns a clientset for the kubeconfig &
// sets the env variable for the same
func GetOutofClusterKubeConfig() {
	var kubeconfig *string
	// config file not provided, auto detect from the host OS
	if util.Kubeconfig == "" {
		if home := homeDir(); home != "" {
			cfg := filepath.Join(home, ".kube", "config")
			kubeconfig = &cfg
		} else {
			log.Fatal(`kubeconfig not provided, Please provide config file path with "--kubeconfig" flag`)
		}
	} else {
		// Get the kubeconfig file from CLI args
		kubeconfig = &util.Kubeconfig
	}
	err := os.Setenv("KUBECONFIG", *kubeconfig)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("KUBECONFIG")
}

func getKubeConfig() *rest.Config {
	// use the current context in kubeconfig
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Print("Could not load in cluster config")
		GetOutofClusterKubeConfig()
		os.Getenv("KUBECONFIG")
		// TODO : make interface of ClientSet and put K8sClient and lvmClient in it and pass it in this function
		// TODO : to create clientset like this, [clientset, err := kubernetes.NewForConfig(config)]
		return nil
	}
	return config
}
