package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeployment(t *testing.T) {
	dName := "redis-master"

	dd, err := deploymentsByName(dName)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(dd), 1)
	if t.Failed() {
		return
	}

	d := dd[0]
	assert.Equal(t, dName, d.ObjectMeta.Name)
	assert.Equal(t, int32(1), *d.Spec.Replicas)
	assert.Equal(t, int32(1), d.Status.Replicas)
	assert.Equal(t, "redis", d.Spec.Template.Labels["app"])
	assert.True(t, dockerImageIsVersionned(d.Spec.Template.Spec.Containers[0].Image))
}
