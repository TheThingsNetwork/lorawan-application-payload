// Copyright © 2021 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package apppayload_test

import (
	"strconv"
	"testing"

	apppayload "go.thethings.network/lorawan-application-payload"
)

func TestInferLocation(t *testing.T) {
	for i, tc := range []struct {
		m   map[string]interface{}
		loc apppayload.Location
		ok  bool
	}{
		{
			m: map[string]interface{}{
				"gps_5": map[string]interface{}{
					"latitude":  float64(1),
					"longitude": float64(2),
					"altitude":  float64(3),
				},
			},
			loc: apppayload.Location{
				Latitude:  1,
				Longitude: 2,
				Altitude:  3,
			},
			ok: true,
		},
		{
			m: map[string]interface{}{
				"lat":          float64(1),
				"longitudeDeg": float64(2),
			},
			ok: false, // invalid pair
		},
		{
			m: map[string]interface{}{
				"lat": 1,
				"lon": 2,
			},
			ok: false, // invalid numeric type
		},
		{
			m: map[string]interface{}{
				"lat": float64(1),
				"lon": float64(2),
			},
			loc: apppayload.Location{
				Latitude:  1,
				Longitude: 2,
			},
			ok: true,
		},
		{
			m: map[string]interface{}{
				"latitude":  float64(1),
				"longitude": float64(2),
				"altitude":  float64(3),
				"accuracy":  float64(4),
			},
			loc: apppayload.Location{
				Latitude:  1,
				Longitude: 2,
				Altitude:  3,
				Accuracy:  4,
			},
			ok: true,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			loc, ok := apppayload.InferLocation(tc.m)
			if ok != tc.ok {
				t.Fatalf("Expected location to return `%v` but it was `%v`", tc.ok, ok)
			}
			if loc != tc.loc {
				t.Fatalf("Expected location to be `%v` but it was `%v`", tc.loc, loc)
			}
		})
	}
}
