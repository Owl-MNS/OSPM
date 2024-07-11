//go:build critical

package organization

import (
	"ospm/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientIPCanListAllOrganization(t *testing.T) {
	type testCase struct {
		name           string
		whitelist      string
		clientIP       string
		expectedResult bool
	}

	testCases := []testCase{
		{
			name:           "client ip exists in the whitelist but whitelist is using absolute IP instead of ip range. In this case, the output should be true",
			whitelist:      "172.16.1.5/32,192.168.1.12",
			clientIP:       "192.168.1.12",
			expectedResult: true,
		},
		{
			name:           "client ip does not exist in the whitelist. In this case, the output should be false",
			whitelist:      "172.16.1.5/32,192.168.1.13/32",
			clientIP:       "192.168.1.12",
			expectedResult: false,
		},
	}

	config.LoadOSPMConfigs()

	for _, tc := range testCases {
		// ensures that each subtest gets its own copy of the test case.
		// This is important because the loop variable tc would otherwise be shared across all goroutines,
		// leading to race conditions and incorrect test results
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			config.OSPM.ClientPolicies.ListAllOrganizationWhiteListedIPs = tc.whitelist
			output := ClientIPCanListAllOrganization(tc.clientIP)

			assert.Equal(t, tc.expectedResult, output)
		})
	}

}
