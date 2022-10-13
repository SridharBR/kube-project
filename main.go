package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	var kubeconfig *string
	var version *string
	var scale *int

	kubeconfig = flag.String("kubeconfig", "", "Absolute path for kubeconfig file")
	version = flag.String("version", "", "Nginx version")
	scale = flag.Int("scale", 0, "No of pods")

	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	result, depCliGetErr := deploymentsClient.Get(context.Background(), "kubeproject", metav1.GetOptions{})

	// If Get fails, we dont have any deployments with this name. So create a new one.
	if depCliGetErr != nil {
		deployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name: "kubeproject",
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: GetIntPtr(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "nginx",
					},
				},
				Template: apiv1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": "nginx",
						},
					},
					Spec: apiv1.PodSpec{
						Containers: []apiv1.Container{
							{
								Name:  "nginx_server",
								Image: "nginx:1.12",
								Ports: []apiv1.ContainerPort{
									{
										Name:          "http",
										Protocol:      apiv1.ProtocolTCP,
										ContainerPort: 80,
									},
								},
							},
						},
					},
				},
			},
		}

		fmt.Println("Creating a new deployment from go tool")

		result, depCliCreateErr := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
		if depCliCreateErr != nil {
			panic(depCliCreateErr)
		}
		fmt.Println("Created deployment : ", result.GetObjectMeta().GetName())
	} else { // If deployment exists, update version and pods as received from CLI
		updateFlag := false
		if *scale != 0 {
			result.Spec.Replicas = GetIntPtr(int32(*scale))
			updateFlag = true

		}
		if *version != "" {
			result.Spec.Template.Spec.Containers[0].Image = GetNginxVersion(*version)
			updateFlag = true
		}

		if updateFlag {
			_, depCliUpdateErr := deploymentsClient.Update(context.TODO(), result, metav1.UpdateOptions{})
			if depCliUpdateErr != nil {
				panic(depCliUpdateErr)
			} else {
				fmt.Println("Updated pods successfully")
			}
		}
	}

}

// Generate nginx version string
func GetNginxVersion(ver string) string {
	strSlc := make([]string, 0, 2)
	strSlc = append(strSlc, "nginx")
	strSlc = append(strSlc, ver)
	str1 := strings.Join(strSlc, ":")

	return str1
}

func GetIntPtr(i int32) *int32 {
	return &i
}
