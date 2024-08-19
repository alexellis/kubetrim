package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/alexellis/kubetrim/pkg"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	_ "k8s.io/client-go/plugin/pkg/client/auth/exec"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"  // Import for GCP auth
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc" // Import for OIDC (often used with EKS)
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"k8s.io/client-go/util/homedir"
)

var (
	writeFile bool
	force     bool
)

func main() {

	// set usage:

	flag.Usage = func() {

		fmt.Printf("kubetrim (%s %s) Copyright Alex Ellis (c) 2024\n\n", pkg.Version, pkg.GitCommit)
		fmt.Print("Sponsor Alex on GitHub: https://github.com/sponsors/alexellis\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "kubetrim removes contexts & clusters from your kubeconfig if they are not accessible.\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Set the kubeconfig file to use through the KUBECONFIG environment variable\n")
		fmt.Fprintf(flag.CommandLine.Output(), "or the default location will be used: ~/.kube/config\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")

		flag.PrintDefaults()
	}

	flag.BoolVar(&writeFile, "write", true, "Write changes to the kubeconfig file")
	flag.BoolVar(&force, "force", false, "Force delete all contexts, even if all are unreachable")
	flag.Parse()

	// Load the kubeconfig file
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	if configPath := os.Getenv("KUBECONFIG"); configPath != "" {
		kubeconfig = configPath
	}

	// Load the kubeconfig
	config, err := clientcmd.LoadFromFile(kubeconfig)
	if err != nil {
		fmt.Printf("Error loading kubeconfig: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("kubetrim (%s %s) by Alex Ellis \n\nLoaded: %s. Checking..\n", pkg.Version, pkg.GitCommit, kubeconfig)

	st := time.Now()
	// List of contexts to be deleted
	var contextsToDelete []string

	// Enumerate and check all contexts
	for contextName := range config.Contexts {
		fmt.Printf("  - %s: ", contextName)

		// Set the context for the current iteration
		clientConfig := clientcmd.NewNonInteractiveClientConfig(*config, contextName, &clientcmd.ConfigOverrides{}, nil)
		restConfig, err := clientConfig.ClientConfig()
		if err != nil {
			fmt.Printf("Failed to create REST config (%v)\n", err)
			contextsToDelete = append(contextsToDelete, contextName)
			continue
		}

		// Create Kubernetes client
		clientset, err := kubernetes.NewForConfig(restConfig)
		if err != nil {
			fmt.Printf("Failed to create clientset (%v)\n", err)
			contextsToDelete = append(contextsToDelete, contextName)
			continue
		}

		// Check if the context is working
		err = checkCluster(clientset)
		if err != nil {
			fmt.Printf("❌ - (%v)\n", err)
			contextsToDelete = append(contextsToDelete, contextName)
		} else {
			fmt.Println("✅")
		}
	}

	// Delete the contexts that are not working
	for _, contextName := range contextsToDelete {
		deleteContextAndCluster(config, contextName)
	}

	if writeFile {

		if len(contextsToDelete) == len(config.Contexts) && !force {
			fmt.Println("No contexts are working, the Internet may be down, use --force to delete all contexts anyway.")
			os.Exit(1)
		}
		if len(contextsToDelete) > 0 {
			// Save the modified kubeconfig
			if err = clientcmd.WriteToFile(*config, kubeconfig); err != nil {
				fmt.Printf("Error saving updated kubeconfig: %v\n", err)
				os.Exit(1)
			}
		}
		fmt.Printf("Updated: %s (in %s).\n", kubeconfig, time.Since(st).Round(time.Millisecond))
	}
}

// checkCluster tries to list nodes in the cluster to verify if the context is working
func checkCluster(clientset *kubernetes.Clientset) error {
	_, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to connect to cluster: %v", err)
	}
	return nil
}

// deleteContextAndCluster removes the context and its associated cluster from the config
func deleteContextAndCluster(config *clientcmdapi.Config, contextName string) {

	// Get the cluster name associated with the context
	clusterName := config.Contexts[contextName].Cluster

	// Delete the context
	delete(config.Contexts, contextName)

	// Check if any other contexts use this cluster
	clusterUsed := false
	for _, ctx := range config.Contexts {
		if ctx.Cluster == clusterName {
			clusterUsed = true
			break
		}
	}

	// If the cluster is not used by any other context, delete it
	if !clusterUsed {
		delete(config.Clusters, clusterName)
	}
}
