package tests

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
)

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
		return safePodSecurityContext(p)
	}

	if privilegeEscalationAllowed(c) {
		log.Println("possible privilege escalation found")
		return false
	}

	if !runAsNonRootEnforced(*c, p) {
		log.Println("run as non root not enforced")
		return false
	}

	return true
}

// TODO : fix the little pointer game
func runAsNonRootEnforced(c v1.SecurityContext, p *v1.PodSecurityContext) bool {
	runAsNonRoot := false
	if c.RunAsNonRoot == nil && c.RunAsUser == nil {
		log.Println("pod level")
		runAsNonRoot = safePodSecurityContext(p)
	} else {
		log.Println("container level")
		runAsNonRoot = safeContainerSecurityContext(&c)
	}
	return runAsNonRoot
}

// container
func privilegeEscalationAllowed(c *v1.SecurityContext) bool {
	return c.AllowPrivilegeEscalation != nil && *c.AllowPrivilegeEscalation
}

func safeContainerSecurityContext(c *v1.SecurityContext) bool {
	return runsAsNonRootContainer(c) || runsAsUserNotRootContainer(c)
}

func runsAsNonRootContainer(c *v1.SecurityContext) bool {
	return c.RunAsNonRoot != nil && *c.RunAsNonRoot
}

func runsAsUserNotRootContainer(c *v1.SecurityContext) bool {
	return c.RunAsUser != nil && *c.RunAsUser != int64(0)
}

// pod
func safePodSecurityContext(c *v1.PodSecurityContext) bool {
	return runsAsNonRoot(c) && runsAsUserNotRoot(c)
}

func runsAsNonRoot(c *v1.PodSecurityContext) bool {
	return c.RunAsNonRoot != nil && *c.RunAsNonRoot
}

func runsAsUserNotRoot(c *v1.PodSecurityContext) bool {
	return c.RunAsUser != nil && *c.RunAsUser != int64(0)
}
