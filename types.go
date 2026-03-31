package main

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type PodIssue struct {
	Namespace string `json:"namespace"`
	Pod string `json:"pod"`
	Issue string `json:"issue"`
	RestartCount int32 `json:"restart_count"`
}

type Node struct {
	Name string `json:"name"`
	Status string `json:"status"`
	Version string `json:"version"`
}

type ScannerReport struct {
	CurrentContextName string `json:"current_context_name"`
	NodeCount int `json:"node_count"`
	Nodes []Node `json:"nodes"`
	NamespaceCount int `json:"namespace_count"`
	PodIssues []PodIssue `json:"pod_issues"`
}

type PodLister interface {
	ListPods(ctx context.Context, namespace string) ([]v1.Pod, error)
}

type K8sPodLister struct {
	clientset *kubernetes.Clientset
}

type TestPodLister struct {
	Pods []v1.Pod
	Err error
}

var (
	// System-namespaces
	systemNamespaces []string = []string{
		"cilium-secrets",
		"kube-node-lease",
		"kube-public",
		"kube-system",
		"longhorn-system",
		"vso-system",
		"traefik",
		"velero",
		"cnpg-system",
		"cert-manager",
	}

	// Statuses
	containerWaitingStatuses []string = []string{
		"ContainerCreating",
		"CrashLoopBackOff",
		"ErrImagePull",
		"ImagePullBackOff",
		"CreateContainerConfigError",
		"InvalidImageName",
		"CreateContainerError",
	}
	containerTerminatedStatuses []string = []string{
		"OOMKilled",
		"Error",
		"Completed",
		"ContainerCannotRun",
		"DeadlineExceeded",
	}
)

const (
	// Treshold to mark container's issue: HighRestartCount
	containerRestartCountTreshold int32 = 5
	// Maximum number of goroutines that scan namespaces in parallel
	maxParallelRequests int32 = 5
)

func NewScannerReport() *ScannerReport {
	// Inizialize Slice
	return &ScannerReport{
		PodIssues: make([]PodIssue, 0),
		Nodes: make([]Node, 0),
	}
}