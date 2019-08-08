package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// ScaleStruct - json payload for scale resource returned by Openshift rest api
type ScaleStruct struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   struct {
		Name              string    `json:"name"`
		Namespace         string    `json:"namespace"`
		SelfLink          string    `json:"selfLink"`
		UID               string    `json:"uid"`
		ResourceVersion   string    `json:"resourceVersion"`
		CreationTimestamp time.Time `json:"creationTimestamp"`
	} `json:"metadata"`
	Spec struct {
		Replicas int `json:"replicas"`
	} `json:"spec"`
	Status struct {
		Replicas int `json:"replicas"`
		Selector struct {
			App              string `json:"app"`
			Deploymentconfig string `json:"deploymentconfig"`
		} `json:"selector"`
		TargetSelector string `json:"targetSelector"`
	} `json:"status"`
}

// DeploymentConfig - json payload for single deploymentconfig returned by Openshift rest api
type DeploymentConfig struct {
	Metadata struct {
		Name              string    `json:"name"`
		Namespace         string    `json:"namespace"`
		SelfLink          string    `json:"selfLink"`
		UID               string    `json:"uid"`
		ResourceVersion   string    `json:"resourceVersion"`
		Generation        int       `json:"generation"`
		CreationTimestamp time.Time `json:"creationTimestamp"`
		Labels            struct {
			App      string `json:"app"`
			Group    string `json:"group"`
			Template string `json:"template"`
		} `json:"labels"`
		Annotations struct {
			KubectlKubernetesIoLastAppliedConfiguration string `json:"kubectl.kubernetes.io/last-applied-configuration"`
		} `json:"annotations"`
	} `json:"metadata"`
	Spec struct {
		Strategy struct {
			Type          string `json:"type"`
			RollingParams struct {
				UpdatePeriodSeconds int             `json:"updatePeriodSeconds"`
				IntervalSeconds     int             `json:"intervalSeconds"`
				TimeoutSeconds      int             `json:"timeoutSeconds"`
				MaxUnavailable      string          `json:"maxUnavailable"`
				MaxSurge            json.RawMessage `json:"maxSurge"`
			} `json:"rollingParams"`
			Resources struct {
			} `json:"resources"`
			ActiveDeadlineSeconds int `json:"activeDeadlineSeconds"`
		} `json:"strategy"`
		Triggers []struct {
			Type string `json:"type"`
		} `json:"triggers"`
		Replicas int  `json:"replicas"`
		Test     bool `json:"test"`
		Selector struct {
			App              string `json:"app"`
			Deploymentconfig string `json:"deploymentconfig"`
		} `json:"selector"`
		Template struct {
			Metadata struct {
				CreationTimestamp interface{} `json:"creationTimestamp"`
				Labels            struct {
					App              string `json:"app"`
					DeploymentTime   string `json:"deployment_time"`
					Deploymentconfig string `json:"deploymentconfig"`
					GitCommit        string `json:"git_commit"`
				} `json:"labels"`
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
				} `json:"containers"`
				RestartPolicy                 string `json:"restartPolicy"`
				TerminationGracePeriodSeconds int    `json:"terminationGracePeriodSeconds"`
				DNSPolicy                     string `json:"dnsPolicy"`
				NodeSelector                  struct {
					Env   string `json:"env"`
					Mule  string `json:"mule"`
					Owner string `json:"owner"`
				} `json:"nodeSelector"`
				SecurityContext struct {
				} `json:"securityContext"`
			} `json:"spec"`
		} `json:"template"`
	} `json:"spec"`
	Status struct {
		LatestVersion       int `json:"latestVersion"`
		ObservedGeneration  int `json:"observedGeneration"`
		Replicas            int `json:"replicas"`
		UpdatedReplicas     int `json:"updatedReplicas"`
		AvailableReplicas   int `json:"availableReplicas"`
		UnavailableReplicas int `json:"unavailableReplicas"`
		Details             struct {
			Message string `json:"message"`
			Causes  []struct {
				Type string `json:"type"`
			} `json:"causes"`
		} `json:"details"`
		Conditions []struct {
			Type               string    `json:"type"`
			Status             string    `json:"status"`
			LastUpdateTime     time.Time `json:"lastUpdateTime"`
			LastTransitionTime time.Time `json:"lastTransitionTime"`
			Reason             string    `json:"reason,omitempty"`
			Message            string    `json:"message"`
		} `json:"conditions"`
		ReadyReplicas int `json:"readyReplicas"`
	} `json:"status"`
}

// DeploymentConfigList - json payload for all deploymentconfigs returned by Openshift rest api
type DeploymentConfigList struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   struct {
		SelfLink        string `json:"selfLink"`
		ResourceVersion string `json:"resourceVersion"`
	} `json:"metadata"`
	Items []DeploymentConfig `json:"items"`
}

func getDeployments(namespace string, token string, command string, baseURL string) {
	url := fmt.Sprintf(baseURL+"/oapi/v1/namespaces/%s/deploymentconfigs", namespace)
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

	var deploymentList DeploymentConfigList
	err = json.Unmarshal(body, &deploymentList)
	if err != nil {
		fmt.Println("There was an error:", err)
	}

	if len(deploymentList.Items) == 0 {
		fmt.Printf("\nThere are no deployments on %s namespace, it is either empty or doesn't exist.\n\n", namespace)
	} else {
		// Get a list of names of each deploymentconfig in the namespace
		names := make([]string, 0)
		for _, deployment := range deploymentList.Items {
			names = append(names, deployment.Metadata.Name)
		}

		scalePods(names, namespace, token, command, client, baseURL)

		fmt.Printf("\nSuccessfully scaled pods.\n")
	}
}

func scalePods(deployNames []string, namespace string, token string, command string,
	client *http.Client, baseURL string) {

	fmt.Printf("Scaling pods on %s namespace\n\n", namespace)

	for _, deploymentName := range deployNames {
		// Scale each deploymentconfig
		url := fmt.Sprintf(baseURL+"/oapi/v1/namespaces/%s/deploymentconfigs/%s/scale",
			namespace, deploymentName)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("Authorization", "Bearer "+token)
		resp, _ := client.Do(req)
		body, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}

		var scale ScaleStruct
		errUnmarshal := json.Unmarshal(body, &scale)
		if errUnmarshal != nil {
			fmt.Println("There was an error:", errUnmarshal)
		}

		if command == "scale-up" {
			fmt.Printf("Scaling deployment %s from %d -> %d\n", deploymentName, scale.Spec.Replicas, scale.Spec.Replicas+1)
			scale.Spec.Replicas++
		} else {
			if scale.Spec.Replicas == 0 {
				continue // Already zero replicas - can't scale below zero
			}
			fmt.Printf("Scaling deployment %s from %d -> %d\n", deploymentName, scale.Spec.Replicas, scale.Spec.Replicas-1)
			scale.Spec.Replicas--
		}

		payload, _ := json.Marshal(scale)
		req, _ = http.NewRequest("PUT", url, bytes.NewBuffer(payload))
		req.Header.Add("Authorization", "Bearer "+token)
		resp, _ = client.Do(req)
		body, readErr = ioutil.ReadAll(resp.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}
	}
}
