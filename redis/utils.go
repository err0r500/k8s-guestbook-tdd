package redis

import (
	"regexp"

	p "github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func IsValidDNS(s string) bool {
	return regexp.MustCompile(`^[a-z0-9\-\.\_]+$`).MatchString(s) && len(s) < 256
}

func stringMapToPulumi(ll map[string]string) p.StringMap {
	pm := p.StringMap{}
	for k, v := range ll {
		pm[k] = p.String(v)
	}

	return pm
}
