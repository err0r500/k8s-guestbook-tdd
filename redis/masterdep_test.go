package redis_test

import (
	"sync"
	"testing"

	"github.com/err0r500/k8s-guestbook-tdd/redis"
	"github.com/pulumi/pulumi/sdk/v2/go/common/resource"
	p "github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	"github.com/stretchr/testify/assert"
)

type mocks struct{}

// Create the mock.
func (mocks) NewResource(typeToken, name string, inputs resource.PropertyMap, provider, id string) (string, resource.PropertyMap, error) {
	return name + "_id", inputs, nil
}

func (mocks) Call(token string, args resource.PropertyMap, provider string) (resource.PropertyMap, error) {
	return args, nil
}

func TestRedisMasterDeployment(t *testing.T) {
	err := p.RunErr(func(c *p.Context) error {
		d, _, err := redis.MasterDeployment(c)
		assert.NoError(t, err)

		wg := sync.WaitGroup{}
		wg.Add(1)

		p.All(d.URN(), d.Metadata.Elem().Name()).ApplyT(func(all []interface{}) error {
			n := all[1].(*string)
			assert.True(t, redis.IsValidDNS(*n))
			wg.Done()
			return nil
		})

		wg.Wait()
		return nil

	}, p.WithMocks("project", "stack", mocks{}))

	assert.NoError(t, err)

}
