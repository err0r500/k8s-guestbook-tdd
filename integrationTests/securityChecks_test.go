package tests

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
)

// TODO : fix the little pointer game with security contexts
func TestSecurityContexts(t *testing.T) {
	pp, err := podsInNamespace("default")
	assert.Equal(t, err, nil)
	if t.Failed() {
		return
	}

	for _, p := range pp {
		for _, c := range p.Spec.Containers {
			assert.True(t, isSafe(c.SecurityContext, p.Spec.SecurityContext), fmt.Sprintf("%s is not safe in pod %s", c.Name, p.Name))
		}
	}
}

func isSafe(c *v1.SecurityContext, p *v1.PodSecurityContext) bool {
	if c == nil {
		return runsAsNonRootPod(p) || runsAsUserNotRootPod(p)
	}

	if privilegeEscalationAllowed(c) {
		log.Println("possible privilege escalation found")
		return false
	}

	return runAsNonRootEnforced(*c, p)
}

func runAsNonRootEnforced(c v1.SecurityContext, p *v1.PodSecurityContext) bool {
	runAsNonRootEnforced := runsAsNonRootContainer(c) || (c.RunAsNonRoot == nil && runsAsNonRootPod(p))
	runAsNotRootUserEnforced := runsAsUserNotRootContainer(c) || (c.RunAsUser == nil && runsAsUserNotRootPod(p))
	return runAsNonRootEnforced || runAsNotRootUserEnforced
}

func privilegeEscalationAllowed(c *v1.SecurityContext) bool {
	return c.AllowPrivilegeEscalation != nil && *c.AllowPrivilegeEscalation
}

func runsAsNonRootContainer(c v1.SecurityContext) bool {
	return c.RunAsNonRoot != nil && *c.RunAsNonRoot
}

func runsAsUserNotRootContainer(c v1.SecurityContext) bool {
	return c.RunAsUser != nil && *c.RunAsUser != int64(0)
}

func runsAsNonRootPod(c *v1.PodSecurityContext) bool {
	return c.RunAsNonRoot != nil && *c.RunAsNonRoot
}

func runsAsUserNotRootPod(c *v1.PodSecurityContext) bool {
	return c.RunAsUser != nil && *c.RunAsUser != int64(0)
}
