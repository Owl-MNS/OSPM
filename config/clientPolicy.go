package config

import (
	"os"
	"strings"
)

type ClientPolicy struct {
	OrganizationSoftDeleteWhiteListedIPs string
	OrganizationHardDeleteWhiteListedIPs string
}

func LoadClientPolicies() *ClientPolicy {
	loadedClientPolicies := &ClientPolicy{}

	loadedClientPolicies.OrganizationSoftDeleteWhiteListedIPs = os.Getenv("ORGANIZATION_SOFT_DELETE_CLIENT_WHITELIST_IP")
	loadedClientPolicies.OrganizationSoftDeleteWhiteListedIPs = strings.ReplaceAll(loadedClientPolicies.OrganizationSoftDeleteWhiteListedIPs, " ", "")
	if loadedClientPolicies.OrganizationSoftDeleteWhiteListedIPs == "" {
		loadedClientPolicies.OrganizationSoftDeleteWhiteListedIPs = "0.0.0.0/0"
	}

	loadedClientPolicies.OrganizationHardDeleteWhiteListedIPs = os.Getenv("ORGANIZATION_HARD_DELETE_CLIENT_WHITELIST_IP")
	loadedClientPolicies.OrganizationHardDeleteWhiteListedIPs = strings.ReplaceAll(loadedClientPolicies.OrganizationHardDeleteWhiteListedIPs, " ", "")
	if loadedClientPolicies.OrganizationHardDeleteWhiteListedIPs == "" {
		loadedClientPolicies.OrganizationHardDeleteWhiteListedIPs = "0.0.0.0/0"
	}

	return loadedClientPolicies
}
