package main

import (
	"github.com/err0r500/k8s-guestbook-tdd/redis"
	p "github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func main() {
	p.Run(func(c *p.Context) error {
		redis.MasterDeployment(c)
		redis.SlaveDeployment(c)
		return nil
	})
}
