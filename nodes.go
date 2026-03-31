package main

import (
	"context"
	"log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func (r *ScannerReport) handleNodes(clientset *kubernetes.Clientset, ctx context.Context) {
	// --- Handle Node logic --- //
	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error listing nodes: %s", err)
	}

	r.NodeCount = len(nodes.Items)

	for _, node := range nodes.Items {
		var status string
		version := node.Status.NodeInfo.KubeletVersion


		for _, condition := range node.Status.Conditions {
			if condition.Type == v1.NodeReady {
				switch condition.Status {
					case v1.ConditionTrue:
						status = "Ready"
					case v1.ConditionFalse:
						status = "Not Ready"
					default:
						status = "Unknown"
				}
			}
		}
		// Add gathered data to Slice of Nodes
		r.Nodes = append(r.Nodes, Node{
			Name: node.Name,
			Status: status,
			Version: version,
		})
	}
}