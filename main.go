package main

import (
	"log"
	"os"

	"github.com/JustSomeHack/ci-kubernetes/clients"
)

func main() {
	kubeConfig := os.Getenv("KUBE_CONFIG")
	if kubeConfig == "" {
		log.Fatalln("No KUBE_CONFIG found!")
	}

	namespace := os.Getenv("NAMESPACE")
	if namespace == "" {
		namespace = "default"
	}
	log.Printf("Using namespace %s\n", namespace)

	action := os.Getenv("ACTION")
	if action == "" {
		log.Fatalln("No ACTION found")
	}

	kind := os.Getenv("KIND")
	if kind == "" {
		log.Fatalln("No KIND found")
	}

	name := os.Getenv("NAME")
	if name == "" {
		log.Fatalln("No NAME found")
	}

	imageName := os.Getenv("IMAGE_NAME")
	if imageName == "" {
		log.Fatalln("No IMAGE_NAME found")
	}

	imageTag := os.Getenv("IMAGE_TAG")
	if imageTag == "" {
		log.Fatalln("No IMAGE_TAG found")
	}

	client := clients.NewKubeClient(kubeConfig, namespace)

	switch action {
	case "update":
		switch kind {
		case "cronjob":
			client.UpdateCronImage(name, imageName, imageTag)
			break
		case "deployment":
			client.UpdateDeploymentImage(name, imageName, imageTag)
			break
		}
		break
	}
}
