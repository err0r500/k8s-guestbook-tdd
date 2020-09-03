package tests

import (
	"context"
	"log"
	"strings"

	typesv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var clientset *kubernetes.Clientset

func init() {
	conf, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		log.Fatal(err)
	}

	clientset, err = kubernetes.NewForConfig(conf)
	if err != nil {
		log.Fatal(err)
	}
}

func deploymentsByName(name string) ([]typesv1.Deployment, error) {
	deployList, err := clientset.AppsV1().
		Deployments(metav1.NamespaceDefault).
		List(context.Background(),
			metav1.ListOptions{
				FieldSelector: "metadata.name=" + name,
			})

	if err != nil {
		return nil, err
	}

	return deployList.Items, nil
}

func dockerImageIsVersionned(tag string) bool {
	splitted := strings.Split(tag, ":")
	if len(splitted) <= 1 {
		return false
	}
	return splitted[1] != "latest"
}
