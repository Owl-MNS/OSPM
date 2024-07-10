package config

import (
	"os"
	"strings"
)

type ClientPolicy struct {
	OrganizationSoftDeleteWhiteListedIPs     string
	OrganizationHardDeleteWhiteListedIPs     string
	ListAllOrganizationWhiteListedIPs        string
	UndoOrganizationSoftDeleteWhiteListedIPs string
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

	loadedClientPolicies.ListAllOrganizationWhiteListedIPs = os.Getenv("ORGANIZATION_LIST_ALL_CLIENT_WHITELIST_IP")
	loadedClientPolicies.ListAllOrganizationWhiteListedIPs = strings.ReplaceAll(loadedClientPolicies.ListAllOrganizationWhiteListedIPs, " ", "")
	if loadedClientPolicies.ListAllOrganizationWhiteListedIPs == "" {
		loadedClientPolicies.ListAllOrganizationWhiteListedIPs = "0.0.0.0/0"
	}

	loadedClientPolicies.UndoOrganizationSoftDeleteWhiteListedIPs = os.Getenv("UNDO_ORGANIZATION_SOFT_DELETE_CLIENT_WHITELIST_IP")
	loadedClientPolicies.UndoOrganizationSoftDeleteWhiteListedIPs = strings.ReplaceAll(loadedClientPolicies.UndoOrganizationSoftDeleteWhiteListedIPs, " ", "")
	if loadedClientPolicies.UndoOrganizationSoftDeleteWhiteListedIPs == "" {
		loadedClientPolicies.UndoOrganizationSoftDeleteWhiteListedIPs = "0.0.0.0/0"
	}

	return loadedClientPolicies
}
