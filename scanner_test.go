package main

import (
	"context"

	v1 "k8s.io/api/core/v1"
)

type TestPodLister struct {
	Pods []v1.Pod
	Err error
}

func (t TestPodLister) ListPods(ctx context.Context, namespace string) ([]v1.Pod, error) {
	return t.Pods, t.Err
}

// func TestScanNamespace(t *testing.T) {
// 	// Create anonymous slice of struct since its used only once
// 	tests := []struct {
// 		name string
// 		pods []v1.Pod
// 		expected []PodIssue
// 	}{
// 		{
// 			name: "Detect CrashLoopBackOff",
// 			pods: []v1.Pod,
// 		}
// 	}
// }