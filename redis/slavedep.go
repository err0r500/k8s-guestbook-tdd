package redis

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/core/v1"
	p "github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func SlaveDeployment(ctx *p.Context) (*appsv1.Deployment, *corev1.Service, error) {
	labels := map[string]string{
		"app":  "redis",
		"role": "slave",
		"tier": "backend",
	}
	port := 6379

	depDef, svcDef := DeploySvcDef_SingleContainer(
		Deploy_SingleContainer{
			Name:     "redis-slave",
			Labels:   labels,
			Replicas: 2,
			Container: ContainerSpec{
				Name:  "redis",
				Image: "gcr.io/google_samples/gb-redisslave:v3",
				Port:  port,
			},
		},
		&port,
	)

	d, err := appsv1.NewDeployment(ctx, "redisSlaveDeployment", &depDef)
	if err != nil {
		return nil, nil, err
	}

	s, err := corev1.NewService(ctx, "redisSlaveSvc", &svcDef)
	if err != nil {
		return nil, nil, err
	}

	return d, s, nil
}
