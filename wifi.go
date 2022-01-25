// Copyright © 2022 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package apppayload

import (
	"encoding/hex"
	"strings"
)

// AccessPoint is the signal description of a particular WiFi
// access point.
type AccessPoint struct {
	BSSID [6]byte
	RSSI  float64
}

// parseBSSID parses the provided string form hexadecimal BSSID
// into a [6]byte. Separators such as - and : are removed before
// the conversion takes place, and the conversion is case insensitive.
func parseBSSID(s string) ([6]byte, bool) {
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, ":", "")
	bs, err := hex.DecodeString(s)
	if err != nil {
		return [6]byte{}, false
	}
	if len(bs) != 6 {
		return [6]byte{}, false
	}
	return *(*[6]byte)(bs), true
}

// InferWiFiAccessPoints uses a predefined set of rules in order to infer the WiFi
// Access Point information from the provided map.
//
// The function will attempt to parse the following structures:
// - An entry in the map called `access_points` which is an array of objects.
// The objects are expected to have the AP BSSID inside the `bssid` key
// and the RSSI in the `rssi` key.
// - An entry in the map called `wifi` which is an array of objects.
// The objects are expected to have the AP BSSID inside the `mac` key
// and the RSSI inside the `rssi` key.
//
// The BSSIDs are expected to be in hexadecimal format. Separators such as
// - and : are stripped before the BSSID is parsed.
//
// All numeric values are assumed to be float64.
//
// If no access points can be inferred, this function returns false.
func InferWiFiAccessPoints(m map[string]interface{}) ([]AccessPoint, bool) {
	if len(m) == 0 {
		return nil, false
	}

outer:
	for apKey, format := range map[string]struct {
		bssidKey string
		rssiKey  string
	}{
		"access_points": {
			bssidKey: "bssid",
			rssiKey:  "rssi",
		},
		"wifi": {
			bssidKey: "mac",
			rssiKey:  "rssi",
		},
	} {
		accessPoints, ok := m[apKey].([]interface{})
		if !ok {
			continue
		}
		points := []AccessPoint{}
		for _, ap := range accessPoints {
			ap, ok := ap.(map[string]interface{})
			if !ok {
				continue outer
			}
			bssidStr, ok := ap[format.bssidKey].(string)
			if !ok {
				continue outer
			}
			bssid, ok := parseBSSID(bssidStr)
			if !ok {
				continue outer
			}
			rssi, ok := ap[format.rssiKey].(float64)
			if !ok {
				continue outer
			}
			points = append(points, AccessPoint{
				BSSID: bssid,
				RSSI:  rssi,
			})
		}
		return points, true
	}

	return nil, false
}
