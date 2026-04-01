package main

import (
	"context"

	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TestPodLister struct {
	Pods []v1.Pod
	Err error
}

func (t TestPodLister) ListPods(ctx context.Context, namespace string) ([]v1.Pod, error) {
	return t.Pods, t.Err
}

func TestScanNamespace(t *testing.T) {
	// Create anonymous slice of structs since its used only once
	tests := []struct {
		name string
		pods []v1.Pod
		expected []PodIssue
	}{
		{
			name: "Detect CrashLoopBackOff",
			pods: []v1.Pod{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "crash-pod",
						Namespace: "test",
					},
					Status: v1.PodStatus{
						Phase: v1.PodRunning,
						ContainerStatuses: []v1.ContainerStatus{
							{
								State: v1.ContainerState{
									Waiting: &v1.ContainerStateWaiting{
										Reason: "CrashLoopBackOff",
									},
								},
							},
						},
					},
				},
			},
			expected: []PodIssue{
				{Namespace: "test", Pod: "crash-pod", Issue: "CrashLoopBackOff"},
			},
		},
		{
			name: "Healty pods without issues",
			pods: []v1.Pod{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "healthy-pod",
					Namespace: "test",
				},
				Status: v1.PodStatus{
					Phase: v1.PodRunning,
				},
			},
		},
		expected: nil,
		},
	}
	// Test against anonymous slice of structs
	ctx := context.TODO()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fakeLister := TestPodLister{Pods: test.pods}
			got := scanNamespace(fakeLister, ctx, "test")

			// Compare length
			if len(got) != len(test.expected) {
				t.Fatalf("got: %d, expected: %d", len(got), len(test.expected))
			}

			for i := range got {
				if got[i].Issue != test.expected[i].Issue {
					t.Fatalf("Issue: %d, got: %q, want: %q", i, got[i].Issue, test.expected[i].Issue)
				}
			}
		})
	}
}