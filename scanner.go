package main

import (
	"context"
	"log"
	"slices"
	"sync"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func (r *ScannerReport) handleNamespaces(clientset *kubernetes.Clientset, ctx context.Context, printAllNamespaces bool) {
	// --- List Namespaces; Pods and Container issues logic --- //
	namespaces, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Couldn't list all namespaces: %s", err)
	}

	ch := make(chan []PodIssue)
	var wg sync.WaitGroup
	var nsCount int
	// Creaeting buffered channel to fight rate limiting
	sem := make(chan struct{}, maxParallelRequests)

	for _, ns := range namespaces.Items {
		// * Skip system namespaces when parameter `printAllNamespaces` is False (default). It is handled by `Contains` methods.
		if !printAllNamespaces && slices.Contains(systemNamespaces, ns.Name) {
			continue
		}
		nsCount += 1

		wg.Add(1)
		go func(name string) {
			sem <- struct{}{} // acquire - block if X goroutines already running
			defer func() { <- sem }() // release - free slot when done
			defer wg.Done()
			// Debug
			// fmt.Printf("Verifying namespace: %s\n", name)
			issues := scanNamespace(clientset, ctx, name)
			ch <- issues
		} (ns.Name)
	}
	// Save namespace count
	r.NamespaceCount = nsCount

	// Close channel when all goroutines finish
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Collect results from channel
	for issues := range ch {
		r.PodIssues = append(r.PodIssues, issues...)
	}
}

func scanNamespace(clientset *kubernetes.Clientset, ctx context.Context, nsName string) []PodIssue {
	pods, err := clientset.CoreV1().Pods(nsName).List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error listing pods from namespace %s: %s", nsName, err)
	}

	// Pod iteration and status handling
	var issues []PodIssue
	for _, pod := range pods.Items {
		var issue string
		var restartCount int32
		cs := pod.Status.ContainerStatuses

		if pod.Status.Phase == v1.PodPending {
			issue = "Pending"
		}

		// Container iteration and status handling
		for _, status := range cs {
			// Check for waiting status
			if status.State.Waiting != nil && slices.Contains(containerWaitingStatuses, status.State.Waiting.Reason) {
				issue = status.State.Waiting.Reason
			} else if status.State.Terminated != nil && slices.Contains(containerTerminatedStatuses, status.State.Terminated.Reason) {
				// Check for terminated status
				issue = status.State.Terminated.Reason
			} else if status.RestartCount > containerRestartCountTreshold {
				restartCount = status.RestartCount
				issue = "HighRestartCount"
			}

			if issue != "" {
				issues = append(issues, PodIssue{
					Namespace: nsName,
					Pod: pod.Name,
					Issue: issue,
					RestartCount: restartCount,
				})
			}
		}
	}
	return issues
}