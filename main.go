package main

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/meta/v1"
	p "github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func myNginxDeployment(ctx *p.Context) error {
	nginxInput := p.String("nginx")
	appLabels := p.StringMap{
		"app": nginxInput,
	}

	deployment, err := appsv1.NewDeployment(ctx, "my-nginx-deployment", &appsv1.DeploymentArgs{
		Metadata: metav1.ObjectMetaArgs{
			Name: nginxInput,
		},
		Spec: appsv1.DeploymentSpecArgs{
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: appLabels,
			},
			Replicas: p.Int(2),
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: appLabels,
				},
				Spec: &corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						corev1.ContainerArgs{
							Name:  nginxInput,
							Image: p.String("nginx:1.19.2"),
						}},
				},
			},
		},
	})
	if err != nil {
		return err
	}

	ctx.Export("name", deployment.Metadata.Elem().Name())
	return nil
}

func main() {
	p.Run(myNginxDeployment)
}
