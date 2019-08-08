package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app  = kingpin.New("oc-utils", "A command-line utility to provide extra functionality for OpenShift management.")
	prod = app.Flag("prod", "Use production Openshift (default is non-prod).").Bool()

	scaleup     = app.Command("scaleup", "Scale up all the pods in a given namespace.")
	namespaceUp = scaleup.Arg("namespace", "Openshift namespace to execute command in.").Required().String()

	scaledown     = app.Command("scaledown", "Scale down all the pods in a given namespace.")
	namespaceDown = scaledown.Arg("namespace", "Openshift namespace to execute command in.").Required().String()

	restartPods = app.Command("restartpods",
		"Restart all the pods in a namespace that are older than given threshold.")
	namespaceRestart = restartPods.Arg("namespace", "Openshift namespace to execute command in.").Required().String()
	threshold        = restartPods.Arg("threshold",
		"Threshold to determine whether a pod should be restarted").Required().Int()
)

func main() {

	// Get Openshift token for the current user
	out, err := exec.Command("oc", "whoami", "-t").Output()
	if err != nil {
		fmt.Printf("\nYou need to login to OpenShift before running this!\n")
		fmt.Printf("Use command 'oc login https://osenp.mlcinsurance.com.au:8443' for example\n\n")
		os.Exit(1)
	}
	token := strings.TrimSpace(string(out))

	argsParse := kingpin.MustParse(app.Parse(os.Args[1:]))
	baseURL := ""

	// Check if using prod or non-prod Openshift cluster
	if *prod {
		baseURL = "https://ose.mlcinsurance.com.au:8443"
	} else {
		baseURL = "https://osenp.mlcinsurance.com.au:8443"
	}

	switch argsParse {
	// Scale up pods
	case scaleup.FullCommand():
		getDeployments(*namespaceUp, token, "scale-up", baseURL)

	// Scale down pods
	case scaledown.FullCommand():
		getDeployments(*namespaceDown, token, "scale-down", baseURL)

	// Restart old pods
	case restartPods.FullCommand():
		restartOldPods(*namespaceRestart, token, *threshold, baseURL)
	}
}
