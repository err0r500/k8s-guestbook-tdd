package redis

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/core/v1"
	p "github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func MasterDeployment(ctx *p.Context) (*appsv1.Deployment, *corev1.Service, error) {
	depDef, svcDef := DeploySvcDef_SingleContainer(
		Deploy_SingleContainer{
			Name: "redis-master",
			Labels: map[string]string{
				"app":  "redis",
				"role": "master",
				"tier": "backend",
			},
			Replicas: 1,
			Container: ContainerSpec{
				Name:  "redis",
				Image: "redis:6.0.7-alpine",
				Port:  6379,
			},
		},
		nil,
	)

	d, err := appsv1.NewDeployment(ctx, "redisMasterDeployment", &depDef)
	if err != nil {
		return nil, nil, err
	}

	s, err := corev1.NewService(ctx, "redisMasterSvc", &svcDef)
	if err != nil {
		return nil, nil, err
	}

	return d, s, nil
}
