package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeployment(t *testing.T) {
	dName := "nginx"

	dd, err := deploymentsByName(dName)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(dd), 1)
	if t.Failed() {
		return
	}

	d := dd[0]
	assert.Equal(t, d.Name, dName)
	assert.GreaterOrEqual(t, *d.Spec.Replicas, int32(2))
	assert.GreaterOrEqual(t, d.Status.Replicas, int32(1))
	assert.Equal(t, d.Spec.Template.Labels["app"], dName)
	assert.True(t, dockerImageIsVersionned(d.Spec.Template.Spec.Containers[0].Image))
}
