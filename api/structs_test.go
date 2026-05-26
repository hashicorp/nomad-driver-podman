// Copyright IBM Corp. 2019, 2026
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"encoding/json"
	"net"
	"strings"
	"testing"

	"github.com/shoenig/test/must"
)

func TestSpecGenerator_NetworksJSONSerialization(t *testing.T) {
	ip := net.ParseIP("10.88.0.100")
	ipv6 := net.ParseIP("fd00:dead:beef::100")

	testCases := []struct {
		name        string
		networkName string
		ips         []*net.IP
		expectKey   string
	}{
		{
			name:        "default network serializes as Networks key",
			networkName: "default",
			ips:         []*net.IP{&ip, &ipv6},
			expectKey:   "Networks",
		},
		{
			name:        "custom network name used as map key",
			networkName: "localv6",
			ips:         []*net.IP{&ipv6},
			expectKey:   "localv6",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			spec := SpecGenerator{
				ContainerNetworkConfig: ContainerNetworkConfig{
					Networks: map[string]PerNetworkOptions{
						tc.networkName: {StaticIPs: tc.ips},
					},
				},
			}

			jsonBytes, err := json.Marshal(spec)
			must.NoError(t, err)

			jsonStr := string(jsonBytes)
			must.True(t, !strings.Contains(jsonStr, `"newNetworks"`))
			must.True(t, strings.Contains(jsonStr, `"`+tc.expectKey+`"`))

			for _, ip := range tc.ips {
				must.True(t, strings.Contains(jsonStr, ip.String()))
			}
		})
	}
}
