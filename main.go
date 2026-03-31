package main

import (
	"context"
	"flag"

	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// --- Create instance of ScannerReport and inizialize Slices --- //
	report := NewScannerReport()

	// --- Set flags --- //
	kubeconfig := flag.String("kubeconfig", "", "Path to your Kubeconfig")
	printAllNamespaces := flag.Bool("all-namespaces", false, "Include all namespaces or print only non-system ones. Default: `False`")
	format := flag.String("format", "", "Output format style. Options: json | table. Default: table")
	flag.Parse()

	// --- Handle Kubeconfig --- //
	kubeconfigPath := SetKubeconfigPath(kubeconfig)
	report.CurrentContextName = clientcmd.GetConfigFromFileOrDie(kubeconfigPath).CurrentContext
	clientset := loadKubeconfig(kubeconfigPath)

	// --- Every kubernetes API call requires context --- //
	// ctx := context.WithTimeout()
	ctx := context.TODO()

	// --- Handle Nodes, Namespaces --- //
	report.handleNodes(clientset,ctx)
	report.handleNamespaces(clientset, ctx, *printAllNamespaces)

	// --- Print Report --- //
	report.RenderReport(format)
}