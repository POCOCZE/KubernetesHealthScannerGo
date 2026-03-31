package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
)

func (r *ScannerReport) renderTable() {
	// Sort slice of structs then print table
	sort.Slice(r.PodIssues, func(i, j int) bool {
		return r.PodIssues[i].Namespace < r.PodIssues[j].Namespace
	})

	// --- Print table result --- //
	fmt.Printf("Connected to cluster: %s\n", r.CurrentContextName)
	fmt.Printf("Nodes found: %v\n\n", r.NodeCount)

	fmt.Printf("%-30s%-10s%s\n", "NAME", "STATUS", "VERSION")
	for _, row := range r.Nodes {
		fmt.Printf("%-30s%-10s%s\n", row.Name, row.Status, row.Version)
	}

	fmt.Printf("\nScanning pods across %v namespaces...\n\n", r.NamespaceCount)

	fmt.Printf("%-20s%-60s%-24s%s\n", "NAMESPACE", "POD", "ISSUE", "RESTARTS")
	for _, row := range r.PodIssues {
		fmt.Printf("%-20s%-60s%-24s%v\n", row.Namespace, row.Pod, row.Issue, row.RestartCount)
	}
}

func (r *ScannerReport) renderJSON() {
	JSON, err := json.MarshalIndent(r, "", " ")
	if err != nil {
		log.Fatalf("Error rendering json output: %s", err)
	}
	fmt.Println(string(JSON))
}

func (r *ScannerReport) RenderReport(format *string) {
	switch *format {
		case "table", "":
			r.renderTable()
		case "json":
			r.renderJSON()
		default:
			log.Fatalf("Error: Got invalid format.")
	}
}