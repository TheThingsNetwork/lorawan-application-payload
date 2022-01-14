// Copyright © 2021 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package apppayload_test

import (
	"reflect"
	"testing"

	apppayload "go.thethings.network/lorawan-application-payload"
)

func TestGNSSInference(t *testing.T) {
	for _, tc := range []struct {
		Name   string
		Input  map[string]interface{}
		Output []byte
		Ok     bool
	}{
		{
			Name: "ValidPayload",
			Input: map[string]interface{}{
				"nav": "aabbccdd",
			},
			Output: []byte{0xaa, 0xbb, 0xcc, 0xdd},
			Ok:     true,
		},
		{
			Name: "InvalidPayload",
			Input: map[string]interface{}{
				"nav": "aabbccddgg",
			},
			Ok: false,
		},
		{
			Name: "NoKey",
			Input: map[string]interface{}{
				"gnss": "aabbccdd",
			},
			Ok: false,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			points, ok := apppayload.InferGNSS(tc.Input)
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
