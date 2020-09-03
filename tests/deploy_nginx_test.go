package deploy_nginx_test

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestDeployment(t *testing.T) {
	deployList, err := clientset.AppsV1().
		Deployments(metav1.NamespaceDefault).
		List(context.Background(),
			metav1.ListOptions{
				FieldSelector: "metadata.name=nginx",
			})

	assert.Equal(t, err, nil)
	assert.Equal(t, len(deployList.Items), 1)

	deploy := deployList.Items[0]
	assert.Equal(t, deploy.Name, "nginx")
	assert.Equal(t, *deploy.Spec.Replicas, int32(2))
}
