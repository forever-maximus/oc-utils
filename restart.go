package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// Pod - json payload for single pod returned by Openshift rest api
type Pod struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   struct {
		Name              string    `json:"name"`
		GenerateName      string    `json:"generateName"`
		Namespace         string    `json:"namespace"`
		SelfLink          string    `json:"selfLink"`
		UID               string    `json:"uid"`
		ResourceVersion   string    `json:"resourceVersion"`
		CreationTimestamp time.Time `json:"creationTimestamp"`
		Labels            struct {
			App              string `json:"app"`
			Deployment       string `json:"deployment"`
			DeploymentTime   string `json:"deployment_time"`
			Deploymentconfig string `json:"deploymentconfig"`
			GitCommit        string `json:"git_commit"`
		} `json:"labels"`
		Annotations struct {
			KubernetesIoCreatedBy                    string `json:"kubernetes.io/created-by"`
			OpenshiftIoDeploymentConfigLatestVersion string `json:"openshift.io/deployment-config.latest-version"`
			OpenshiftIoDeploymentConfigName          string `json:"openshift.io/deployment-config.name"`
			OpenshiftIoDeploymentName                string `json:"openshift.io/deployment.name"`
			OpenshiftIoScc                           string `json:"openshift.io/scc"`
		} `json:"annotations"`
	} `json:"metadata"`
	Spec struct {
		Volumes []struct {
			Name      string `json:"name"`
			ConfigMap struct {
				Name        string `json:"name"`
				DefaultMode int    `json:"defaultMode"`
			} `json:"configMap,omitempty"`
			Secret struct {
				SecretName  string `json:"secretName"`
				DefaultMode int    `json:"defaultMode"`
			} `json:"secret,omitempty"`
		} `json:"volumes"`
		Containers []struct {
			Name  string `json:"name"`
			Image string `json:"image"`
			Ports []struct {
				Name          string `json:"name"`
				ContainerPort int    `json:"containerPort"`
				Protocol      string `json:"protocol"`
			} `json:"ports"`
			Env []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"env"`
			Resources struct {
				Limits struct {
					CPU    string `json:"cpu"`
					Memory string `json:"memory"`
				} `json:"limits"`
				Requests struct {
					CPU    string `json:"cpu"`
					Memory string `json:"memory"`
				} `json:"requests"`
			} `json:"resources"`
			VolumeMounts []struct {
				Name      string `json:"name"`
				ReadOnly  bool   `json:"readOnly,omitempty"`
				MountPath string `json:"mountPath"`
			} `json:"volumeMounts"`
			LivenessProbe struct {
				HTTPGet struct {
					Path   string `json:"path"`
					Port   int    `json:"port"`
					Scheme string `json:"scheme"`
				} `json:"httpGet"`
				InitialDelaySeconds int `json:"initialDelaySeconds"`
				TimeoutSeconds      int `json:"timeoutSeconds"`
				PeriodSeconds       int `json:"periodSeconds"`
				SuccessThreshold    int `json:"successThreshold"`
				FailureThreshold    int `json:"failureThreshold"`
			} `json:"livenessProbe"`
			ReadinessProbe struct {
				HTTPGet struct {
					Path   string `json:"path"`
					Port   int    `json:"port"`
					Scheme string `json:"scheme"`
				} `json:"httpGet"`
				InitialDelaySeconds int `json:"initialDelaySeconds"`
				TimeoutSeconds      int `json:"timeoutSeconds"`
				PeriodSeconds       int `json:"periodSeconds"`
				SuccessThreshold    int `json:"successThreshold"`
				FailureThreshold    int `json:"failureThreshold"`
			} `json:"readinessProbe"`
			TerminationMessagePath string `json:"terminationMessagePath"`
			ImagePullPolicy        string `json:"imagePullPolicy"`
			SecurityContext        struct {
				Capabilities struct {
					Drop []string `json:"drop"`
				} `json:"capabilities"`
				Privileged     bool `json:"privileged"`
				SeLinuxOptions struct {
					Level string `json:"level"`
				} `json:"seLinuxOptions"`
			} `json:"securityContext"`
		} `json:"containers"`
		RestartPolicy                 string `json:"restartPolicy"`
		TerminationGracePeriodSeconds int    `json:"terminationGracePeriodSeconds"`
		DNSPolicy                     string `json:"dnsPolicy"`
		NodeSelector                  struct {
			Env   string `json:"env"`
			Mule  string `json:"mule"`
			Owner string `json:"owner"`
		} `json:"nodeSelector"`
		ServiceAccountName string `json:"serviceAccountName"`
		ServiceAccount     string `json:"serviceAccount"`
		NodeName           string `json:"nodeName"`
		SecurityContext    struct {
			SeLinuxOptions struct {
				Level string `json:"level"`
			} `json:"seLinuxOptions"`
		} `json:"securityContext"`
		ImagePullSecrets []struct {
			Name string `json:"name"`
		} `json:"imagePullSecrets"`
	} `json:"spec"`
	Status struct {
		Phase      string `json:"phase"`
		Conditions []struct {
			Type               string      `json:"type"`
			Status             string      `json:"status"`
			LastProbeTime      interface{} `json:"lastProbeTime"`
			LastTransitionTime time.Time   `json:"lastTransitionTime"`
		} `json:"conditions"`
		HostIP            string    `json:"hostIP"`
		PodIP             string    `json:"podIP"`
		StartTime         time.Time `json:"startTime"`
		ContainerStatuses []struct {
			Name  string `json:"name"`
			State struct {
				Running struct {
					StartedAt time.Time `json:"startedAt"`
				} `json:"running"`
			} `json:"state"`
			LastState struct {
			} `json:"lastState"`
			Ready        bool   `json:"ready"`
			RestartCount int    `json:"restartCount"`
			Image        string `json:"image"`
			ImageID      string `json:"imageID"`
			ContainerID  string `json:"containerID"`
		} `json:"containerStatuses"`
	} `json:"status"`
}

// PodList - json payload for all pods returned by Openshift rest api
type PodList struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   struct {
		SelfLink        string `json:"selfLink"`
		ResourceVersion string `json:"resourceVersion"`
	} `json:"metadata"`
	Items []Pod `json:"items"`
}

func restartOldPods(namespace string, token string, threshold int, baseURL string) {
	// Get all pods in the given namespace
	url := fmt.Sprintf(baseURL+"/api/v1/namespaces/%s/pods", namespace)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	} else if resp == nil {
		log.Fatal("The http request couldn't resolve hostname - you're probably not on the VPN?")
	} else if resp.StatusCode == 401 {
		fmt.Printf("\nYou need to login to OpenShift before running this!\n")
		fmt.Printf("Use command 'oc login https://osenp.mlcinsurance.com.au:8443' for example\n\n")
		os.Exit(1)
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var podList PodList
	err = json.Unmarshal(body, &podList)
	if err != nil {
		fmt.Println("There was an error:", err)
	}

	if len(podList.Items) == 0 {
		fmt.Printf("\nThere are no pods on %s namespace.\n\n", namespace)
	} else {
		// Check the age of each pod - restart any older than 3 days
		restartCount := 0
		thresholdHours := float64(threshold) * 24.0

		for _, pod := range podList.Items {
			diff := time.Now().Sub(pod.Status.StartTime).Hours()
			if diff > thresholdHours {
				deleteURL := url + "/" + pod.Metadata.Name
				req, _ = http.NewRequest("DELETE", deleteURL, nil)
				req.Header.Add("Authorization", "Bearer "+token)
				resp, _ = client.Do(req)
				body, readErr = ioutil.ReadAll(resp.Body)
				if readErr != nil {
					log.Fatal(readErr)
				}
				fmt.Printf("Restarting pod %s\n", pod.Metadata.Name)
				restartCount++
			}
		}

		if restartCount == 0 {
			fmt.Printf("\nNone of the pods on %s namespace are older than %d days\n\n", namespace, threshold)
		} else {
			fmt.Printf("\nSuccessfully restarted %d pods on %s namespace\n\n", restartCount, namespace)
		}
	}
}
