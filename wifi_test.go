// Copyright © 2021 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package apppayload_test

import (
	"reflect"
	"testing"

	apppayload "go.thethings.network/lorawan-application-payload"
)

func TestWiFiInference(t *testing.T) {
	for _, tc := range []struct {
		Name   string
		Input  map[string]interface{}
		Output []apppayload.AccessPoint
		Ok     bool
	}{
		{
			Name: "MalformedBSSID",
			Input: map[string]interface{}{
				"access_points": []interface{}{
					map[string]interface{}{
						"bssid": "14:60:80:9a:19:58t",
						"rssi":  -20.0,
					},
					map[string]interface{}{
						"bssid": "fc:f5:28:7b:07:e5",
						"rssi":  -80.0,
					},
				},
			},
			Ok: false,
		},
		{
			Name: "MissingBSSID",
			Input: map[string]interface{}{
				"access_points": []interface{}{
					map[string]interface{}{
						"rssi": -20.0,
					},
					map[string]interface{}{
						"bssid": "fc:f5:28:7b:07:e5",
						"rssi":  -80.0,
					},
				},
			},
			Ok: false,
		},
		{
			Name: "MissingRSSI",
			Input: map[string]interface{}{
				"access_points": []interface{}{
					map[string]interface{}{
						"bssid": "14:60:80:9a:19:58",
						"rssi":  -20.0,
					},
					map[string]interface{}{
						"bssid": "fc:f5:28:7b:07:e5",
					},
				},
			},
			Ok: false,
		},
		{
			Name: "Colon",
			Input: map[string]interface{}{
				"access_points": []interface{}{
					map[string]interface{}{
						"bssid": "14:60:80:9a:19:58",
						"rssi":  -20.0,
					},
					map[string]interface{}{
						"bssid": "fc:f5:28:7b:07:e5",
						"rssi":  -80.0,
					},
				},
			},
			Output: []apppayload.AccessPoint{
				{
					BSSID: [6]byte{0x14, 0x60, 0x80, 0x9a, 0x19, 0x58},
					RSSI:  -20.0,
				},
				{
					BSSID: [6]byte{0xfc, 0xf5, 0x28, 0x7b, 0x07, 0xe5},
					RSSI:  -80.0,
				},
			},
			Ok: true,
		},
		{
			Name: "Dash",
			Input: map[string]interface{}{
				"access_points": []interface{}{
					map[string]interface{}{
						"bssid": "15-61-81-9b-1a-59",
						"rssi":  -21.1,
					},
					map[string]interface{}{
						"bssid": "fd-f6-29-7c-09-e6",
						"rssi":  -81.1,
					},
				},
			},
			Output: []apppayload.AccessPoint{
				{
					BSSID: [6]byte{0x15, 0x61, 0x81, 0x9b, 0x1a, 0x59},
					RSSI:  -21.1,
				},
				{
					BSSID: [6]byte{0xfd, 0xf6, 0x29, 0x7c, 0x09, 0xe6},
					RSSI:  -81.1,
				},
			},
			Ok: true,
		},
		{
			Name: "FormatB",
			Input: map[string]interface{}{
				"wifi": []interface{}{
					map[string]interface{}{
						"mac":  "a0b3ccd358e6",
						"rssi": -92.0,
					},
				},
			},
			Output: []apppayload.AccessPoint{
				{
					BSSID: [6]byte{0xa0, 0xb3, 0xcc, 0xd3, 0x58, 0xe6},
					RSSI:  -92.0,
				},
			},
			Ok: true,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			points, ok := apppayload.InferWiFiAccessPoints(tc.Input)
			if tc.Ok != ok {
				t.Fatal("Expected access points not found")
			}
			if tc.Ok {
				if !reflect.DeepEqual(tc.Output, points) {
					t.Fatal("Output mismatch")
				}
			}
		})
	}
}
