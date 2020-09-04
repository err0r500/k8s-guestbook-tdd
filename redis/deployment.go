package redis

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/meta/v1"
	p "github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

type ContainerSpec struct {
	Name  string
	Image string
	Port  int
}

type Deploy_SingleContainer struct {
	Name      string
	Labels    map[string]string
	Container ContainerSpec
	Replicas  int
}

type Svc_SingleTarget struct {
	Name           string
	SelectorLabels map[string]string
	TargetPort     int
	Port           *int
}

func SvcDef_SingleTarget(svc Svc_SingleTarget) corev1.ServiceArgs {
	port := svc.TargetPort
	if svc.Port != nil {
		port = *svc.Port
	}

	return corev1.ServiceArgs{
		Metadata: metav1.ObjectMetaArgs{
			Name: p.String(svc.Name),
		},
		Spec: corev1.ServiceSpecArgs{
			Ports: corev1.ServicePortArray{
				corev1.ServicePortArgs{
					Port:       p.Int(port),
					TargetPort: p.Int(svc.TargetPort),
				},
			},
			Selector: stringMapToPulumi(svc.SelectorLabels),
		},
	}
}

func DeploySvcDef_SingleContainer(dep Deploy_SingleContainer, svcPort *int) (appsv1.DeploymentArgs, corev1.ServiceArgs) {
	return DeploymentDef(dep), SvcDef_SingleTarget(
		Svc_SingleTarget{
			Name:           dep.Name,
			SelectorLabels: dep.Labels,
			TargetPort:     dep.Container.Port,
			Port:           svcPort,
		},
	)
}

func DeploymentDef(dep Deploy_SingleContainer) appsv1.DeploymentArgs {
	labels := stringMapToPulumi(dep.Labels)

	containers := corev1.ContainerArray{
		corev1.ContainerArgs{
			SecurityContext: corev1.SecurityContextArgs{
				AllowPrivilegeEscalation: p.Bool(false),
				RunAsUser:                p.Int(1000),
			},
			Name:  p.String(dep.Container.Name),
			Image: p.String(dep.Container.Image),
			Resources: corev1.ResourceRequirementsArgs{
				Limits: p.StringMap{
					"cpu":    p.String("100m"),
					"memory": p.String("100Mi"),
				},
			},
			Ports: corev1.ContainerPortArray{
				&corev1.ContainerPortArgs{
					ContainerPort: p.Int(dep.Container.Port),
				},
			},
		},
	}

	return appsv1.DeploymentArgs{
		Metadata: metav1.ObjectMetaArgs{
			Name: p.String(dep.Name),
		},
		Spec: appsv1.DeploymentSpecArgs{
			Selector: metav1.LabelSelectorArgs{
				MatchLabels: labels,
			},
			Replicas: p.Int(dep.Replicas),
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: metav1.ObjectMetaArgs{
					Labels: labels,
				},
				Spec: &corev1.PodSpecArgs{
					SecurityContext: corev1.PodSecurityContextArgs{
						RunAsUser: p.Int(1000),
					},
					Containers: containers,
				},
			},
		},
	}
}
