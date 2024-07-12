package organization

import (
	"errors"
	"ospm/config"
	"ospm/internal/models"
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
		{
			name:           "whitelist is empty. In this case, the output should be false",
			whitelist:      "",
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

func TestClientIPCanUndoOrganizationSoftDelete(t *testing.T) {
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
		{
			name:           "whitelist is empty. In this case, the output should be false",
			whitelist:      "",
			clientIP:       "192.168.1.12",
			expectedResult: false,
		},
	}

	config.LoadOSPMConfigs()

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			config.OSPM.ClientPolicies.UndoOrganizationSoftDeleteWhiteListedIPs = tc.whitelist
			testResult := ClientIPCanUndoOrganizationSoftDelete(tc.clientIP)
			assert.Equal(t, tc.expectedResult, testResult)
		})
	}

}

func TestClientIPCanHardDeleteOrganization(t *testing.T) {
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
		{
			name:           "whitelist is empty. In this case, the output should be false",
			whitelist:      "",
			clientIP:       "192.168.1.12",
			expectedResult: false,
		},
	}

	config.LoadOSPMConfigs()
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			config.OSPM.ClientPolicies.OrganizationHardDeleteWhiteListedIPs = tc.whitelist
			testResult := ClientIPCanHardDeleteOrganization(tc.clientIP)
			assert.Equal(t, tc.expectedResult, testResult)
		})

	}

}

func TestClientIPCanSoftDeleteOrganization(t *testing.T) {
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
		{
			name:           "whitelist is empty. In this case, the output should be false",
			whitelist:      "",
			clientIP:       "192.168.1.12",
			expectedResult: false,
		},
	}

	config.LoadOSPMConfigs()

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			config.OSPM.ClientPolicies.OrganizationSoftDeleteWhiteListedIPs = tc.whitelist
			testResult := ClientIPCanSoftDeleteOrganization(tc.clientIP)
			assert.Equal(t, tc.expectedResult, testResult)
		})
	}
}

// x
func Test(t *testing.T) {
	type testCase struct {
		name                string
		organizationDetails models.Organization
		expectedResult      error
	}

	testCases := []testCase{
		{
			name:           "the organization details has the values that meets the value rules",
			expectedResult: nil,
			organizationDetails: models.Organization{
				Details: models.OrganizationDetails{
					Name:    "Sample Organization 4",
					Address: "123 Sample Street",
					Email:   "info4@sample.org",
					Mobile:  "1234567890124",
					Phone:   "0987654321",
				},
				Owner: models.OrganizationOwner{
					Type:            "legal",
					Name:            "Ario2 Ahmadi",
					Address:         "456 Owner Avenue",
					Email:           "ario2@example.com",
					Mobile:          "23456789010234",
					Phone:           "1234567891",
					LegalNationalID: "AB1234562",
				},
				Balance:                  0.0,
				AllowNagativeBalance:     false,
				NegativeBalanceThreshold: 0.0,
			},
		},
		{
			name:           "the organization details has a negative balance value that should be rejected",
			expectedResult: errors.New("organization balance can not accept any values but 0 while creating the organization"),
			organizationDetails: models.Organization{
				Details: models.OrganizationDetails{
					Name:    "Sample Organization 4",
					Address: "123 Sample Street",
					Email:   "info4@sample.org",
					Mobile:  "1234567890124",
					Phone:   "0987654321",
				},
				Owner: models.OrganizationOwner{
					Type:            "legal",
					Name:            "Ario2 Ahmadi",
					Address:         "456 Owner Avenue",
					Email:           "ario2@example.com",
					Mobile:          "23456789010234",
					Phone:           "1234567891",
					LegalNationalID: "AB1234562",
				},
				Balance:                  -100.0,
				AllowNagativeBalance:     false,
				NegativeBalanceThreshold: 0.0,
			},
		},
		{
			name:           "the organization details has a int 0 value that should be accepted",
			expectedResult: nil,
			organizationDetails: models.Organization{
				Details: models.OrganizationDetails{
					Name:    "Sample Organization 4",
					Address: "123 Sample Street",
					Email:   "info4@sample.org",
					Mobile:  "1234567890124",
					Phone:   "0987654321",
				},
				Owner: models.OrganizationOwner{
					Type:            "legal",
					Name:            "Ario2 Ahmadi",
					Address:         "456 Owner Avenue",
					Email:           "ario2@example.com",
					Mobile:          "23456789010234",
					Phone:           "1234567891",
					LegalNationalID: "AB1234562",
				},
				Balance:                  0,
				AllowNagativeBalance:     false,
				NegativeBalanceThreshold: 0.0,
			},
		},
		{
			name:           "the organization details has a negative int value that should be rejected",
			expectedResult: errors.New("organization balance can not accept any values but 0 while creating the organization"),
			organizationDetails: models.Organization{
				Details: models.OrganizationDetails{
					Name:    "Sample Organization 4",
					Address: "123 Sample Street",
					Email:   "info4@sample.org",
					Mobile:  "1234567890124",
					Phone:   "0987654321",
				},
				Owner: models.OrganizationOwner{
					Type:            "legal",
					Name:            "Ario2 Ahmadi",
					Address:         "456 Owner Avenue",
					Email:           "ario2@example.com",
					Mobile:          "23456789010234",
					Phone:           "1234567891",
					LegalNationalID: "AB1234562",
				},
				Balance:                  -100,
				AllowNagativeBalance:     false,
				NegativeBalanceThreshold: 0.0,
			},
		},
		{
			name:           "the organization details has a positive int value that should be rejected",
			expectedResult: errors.New("organization balance can not accept any values but 0 while creating the organization"),
			organizationDetails: models.Organization{
				Details: models.OrganizationDetails{
					Name:    "Sample Organization 4",
					Address: "123 Sample Street",
					Email:   "info4@sample.org",
					Mobile:  "1234567890124",
					Phone:   "0987654321",
				},
				Owner: models.OrganizationOwner{
					Type:            "legal",
					Name:            "Ario2 Ahmadi",
					Address:         "456 Owner Avenue",
					Email:           "ario2@example.com",
					Mobile:          "23456789010234",
					Phone:           "1234567891",
					LegalNationalID: "AB1234562",
				},
				Balance:                  100,
				AllowNagativeBalance:     false,
				NegativeBalanceThreshold: 0.0,
			},
		},
		{
			name:           "the organization details has a positive balance value that should be rejected",
			expectedResult: errors.New("organization balance can not accept any values but 0 while creating the organization"),
			organizationDetails: models.Organization{
				Details: models.OrganizationDetails{
					Name:    "Sample Organization 4",
					Address: "123 Sample Street",
					Email:   "info4@sample.org",
					Mobile:  "1234567890124",
					Phone:   "0987654321",
				},
				Owner: models.OrganizationOwner{
					Type:            "legal",
					Name:            "Ario2 Ahmadi",
					Address:         "456 Owner Avenue",
					Email:           "ario2@example.com",
					Mobile:          "23456789010234",
					Phone:           "1234567891",
					LegalNationalID: "AB1234562",
				},
				Balance:                  100.0,
				AllowNagativeBalance:     false,
				NegativeBalanceThreshold: 0.0,
			},
		},
		{
			name:           "the organization AllowNagativeBalance is true that should be rejected",
			expectedResult: errors.New("organization AllowNagativeBalance can not be true while creating the organization"),
			organizationDetails: models.Organization{
				Details: models.OrganizationDetails{
					Name:    "Sample Organization 4",
					Address: "123 Sample Street",
					Email:   "info4@sample.org",
					Mobile:  "1234567890124",
					Phone:   "0987654321",
				},
				Owner: models.OrganizationOwner{
					Type:            "legal",
					Name:            "Ario2 Ahmadi",
					Address:         "456 Owner Avenue",
					Email:           "ario2@example.com",
					Mobile:          "23456789010234",
					Phone:           "1234567891",
					LegalNationalID: "AB1234562",
				},
				Balance:                  0.0,
				AllowNagativeBalance:     true,
				NegativeBalanceThreshold: 0.0,
			},
		},
		{
			name:           "the organization NegativeBalanceThreshold is 0 that should be accepted",
			expectedResult: nil,
			organizationDetails: models.Organization{
				Details: models.OrganizationDetails{
					Name:    "Sample Organization 4",
					Address: "123 Sample Street",
					Email:   "info4@sample.org",
					Mobile:  "1234567890124",
					Phone:   "0987654321",
				},
				Owner: models.OrganizationOwner{
					Type:            "legal",
					Name:            "Ario2 Ahmadi",
					Address:         "456 Owner Avenue",
					Email:           "ario2@example.com",
					Mobile:          "23456789010234",
					Phone:           "1234567891",
					LegalNationalID: "AB1234562",
				},
				Balance:                  0.0,
				AllowNagativeBalance:     false,
				NegativeBalanceThreshold: 0.0,
			},
		},
		{
			name:           "the organization NegativeBalanceThreshold is negative that should be rejected",
			expectedResult: errors.New("organization NegativeBalanceThreshold can not accept any values but 0 while creating the organization"),
			organizationDetails: models.Organization{
				Details: models.OrganizationDetails{
					Name:    "Sample Organization 4",
					Address: "123 Sample Street",
					Email:   "info4@sample.org",
					Mobile:  "1234567890124",
					Phone:   "0987654321",
				},
				Owner: models.OrganizationOwner{
					Type:            "legal",
					Name:            "Ario2 Ahmadi",
					Address:         "456 Owner Avenue",
					Email:           "ario2@example.com",
					Mobile:          "23456789010234",
					Phone:           "1234567891",
					LegalNationalID: "AB1234562",
				},
				Balance:                  0.0,
				AllowNagativeBalance:     false,
				NegativeBalanceThreshold: -100,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			testResult := DetailsCheck(&tc.organizationDetails)
			if testResult == nil {
				assert.Equal(t, tc.expectedResult, testResult)
			} else {
				assert.Contains(t, testResult.Error(), tc.expectedResult.Error())

			}
		})
	}
}

// x
func TestShorten(t *testing.T) {

}
