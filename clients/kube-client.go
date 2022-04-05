package clients

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KubeClient interface to execute the kubectl command
type KubeClient interface {
	UpdateCronImage(jobName string, image string, tag string)
	UpdateDeploymentImage(deployment string, image string, tag string)
}

type kubeClient struct {
	client    *kubernetes.Clientset
	namespace string
}

// NewKubeClient returns a new instance of kubeClient
func NewKubeClient(encodedConfig string, namespace string) KubeClient {
	f, err := os.Create(".kubeConfig")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fileBytes, err := base64.StdEncoding.DecodeString(encodedConfig)
	if err != nil {
		log.Fatal(err)
	}

	n, err := f.Write(fileBytes)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Wrote %d bytes\n", n)

	err = f.Sync()
	if err != nil {
		log.Fatal(err)
	}

	config, err := clientcmd.BuildConfigFromFlags("", ".kubeConfig")
	if err != nil {
		panic(err)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return &kubeClient{
		client:    client,
		namespace: namespace,
	}
}

func (k *kubeClient) UpdateCronImage(jobName string, image string, tag string) {
	imageName := fmt.Sprintf("%s:%s", image, tag)
	job, err := k.client.BatchV1().CronJobs(k.namespace).Get(context.Background(), jobName, metav1.GetOptions{})
	if err != nil {
		log.Fatalf("Unable to get cron job %s\n", err.Error())
	}

	log.Printf("Found cron job %s\n", job.Name)

	for i := range job.Spec.JobTemplate.Spec.Template.Spec.Containers {
		log.Printf("Changing container image from %s to %s", job.Spec.JobTemplate.Spec.Template.Spec.Containers[i].Image, imageName)
		job.Spec.JobTemplate.Spec.Template.Spec.Containers[i].Image = imageName
	}

	log.Printf("Updating cron job %s\n", job.Name)
	job, err = k.client.BatchV1().CronJobs(k.namespace).Update(context.Background(), job, metav1.UpdateOptions{})
	if err != nil {
		log.Fatalf("Unable to get updated cron job %s\n", err.Error())
	}

	for i := range job.Spec.JobTemplate.Spec.Template.Spec.Containers {
		log.Printf("Container %s image is set to %s", job.Spec.JobTemplate.Spec.Template.Spec.Containers[i].Name, imageName)
		job.Spec.JobTemplate.Spec.Template.Spec.Containers[i].Image = imageName
	}
}

func (k *kubeClient) UpdateDeploymentImage(deploymentName string, image string, tag string) {
	imageName := fmt.Sprintf("%s:%s", image, tag)
	deployment, err := k.client.AppsV1().Deployments(k.namespace).Get(context.Background(), deploymentName, metav1.GetOptions{})
	if err != nil {
		log.Fatalf("Unable to get deployment %s\n", err.Error())
	}

	log.Printf("Found deployment %s\n", deployment.Name)

	for i := range deployment.Spec.Template.Spec.Containers {
		log.Printf("Changing container image from %s to %s", deployment.Spec.Template.Spec.Containers[i].Image, imageName)
		deployment.Spec.Template.Spec.Containers[i].Image = imageName
	}

	log.Printf("Updating deployment %s\n", deployment.Name)
	deployment, err = k.client.AppsV1().Deployments(k.namespace).Update(context.Background(), deployment, metav1.UpdateOptions{})
	if err != nil {
		log.Fatalf("Unable to get updated deployment %s\n", err.Error())
	}

	for i := range deployment.Spec.Template.Spec.Containers {
		log.Printf("Container %s image is set to %s", deployment.Spec.Template.Spec.Containers[i].Name, imageName)
		deployment.Spec.Template.Spec.Containers[i].Image = imageName
	}
}
