package main

import (
	"log"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func SetKubeconfigPath(kubeconfig *string) string {
	// --- If user did not specified kubeconfig path - use the default path --- //
	kubeconfigPath := *kubeconfig
	if *kubeconfig == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Error setting default user home dir: %s", err)
		}
		kubeconfigPath = filepath.Join(homeDir, ".kube", "config")
	}

	return kubeconfigPath
}

func loadKubeconfig(kubeconfigPath string) *kubernetes.Clientset {
	// --- Load Kubeconfig --- //
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		log.Fatalf("Error loading kubeconfig: %s", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Clientset: %s", err)
	}

	return clientset
}