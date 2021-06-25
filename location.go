// Copyright © 2021 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package apppayload

import (
	"fmt"
	"regexp"
)

// Location is a geographical location.
type Location struct {
	Latitude,
	Longitude,
	Altitude,
	Accuracy float64
}

// GoString returns the textual representation of location.
func (l Location) GoString() string {
	return fmt.Sprintf("%f %f %fm ±%fm", l.Latitude, l.Longitude, l.Altitude, l.Accuracy)
}

// gpsKeyRegexp matches keys in the form of gps_#. These are commonly used by CayenneLPP in The Things Stack.
var gpsKeyRegexp = regexp.MustCompile(`^gps_\d+$`)

// InferLocation uses a set of predefined rules to determine the location from the given map.
//
// If the map contains a key with the format gps_#, the latitude, longitude and altitude is assumed to be value map.
//
// Otherwise, this function checks for the following key combinations for latitude, longitude and altitude:
//   - lat, lon, alt
//   - lat, lng, alt
//   - lat, long, alt
//   - latitude, longitude, altitude
//   - Latitude, Longitude, Altitude
//   - latitudeDeg, longitudeDeg, altitude
//   - latitudeDeg, longitudeDeg, height
//   - gps_lat, gps_lng, gps_alt
//   - gps_lat, gps_lng, gpsalt
// And the following keys for accuracy:
//   - acc
//   - accuracy
//   - hacc
//   - hdop
//   - gps_hdop
//
// All numeric values are assumed to be float64.
//
// If no location can be inferred, this function returns false.
func InferLocation(m map[string]interface{}) (res Location, ok bool) {
	if len(m) == 0 {
		ok = false
		return
	}

	// Check GPS location (field gps_#).
	for k := range m {
		if gpsKeyRegexp.MatchString(k) {
			if vm, ok := m[k].(map[string]interface{}); ok {
				lat, _ := vm["latitude"].(float64)
				lon, _ := vm["longitude"].(float64)
				alt, _ := vm["altitude"].(float64)
				return Location{
					Latitude:  lat,
					Longitude: lon,
					Altitude:  alt,
				}, true
			}
		}
	}

	// Check location in fields.
	for _, pp := range []struct {
		latKey, lonKey, altKey string
	}{
		{"lat", "lon", "alt"},
		{"lat", "lng", "alt"},
		{"lat", "long", "alt"},
		{"latitude", "longitude", "altitude"},
		{"Latitude", "Longitude", "Altitude"},
		{"latitudeDeg", "longitudeDeg", "altitude"},
		{"latitudeDeg", "longitudeDeg", "height"},
		{"gps_lat", "gps_lng", "gps_alt"},
		{"gps_lat", "gps_lng", "gpsalt"},
	} {
		lat, hasLat := m[pp.latKey].(float64)
		lon, hasLon := m[pp.lonKey].(float64)
		if !hasLat || !hasLon {
			continue
		}
		alt, _ := m[pp.altKey].(float64)
		if lat != 0 || lon != 0 {
			res, ok = Location{
				Latitude:  lat,
				Longitude: lon,
				Altitude:  alt,
			}, true
			break
		}
	}

	if ok {
		// Check accuracy in fields. Horizontal dilution of precision is considered accuracy.
		for _, ap := range []string{
			"acc",
			"accuracy",
			"hacc",
			"hdop",
			"gps_hdop",
		} {
			acc, hasAcc := m[ap].(float64)
			if hasAcc {
				res.Accuracy = acc
				break
			}
		}
	}

	return
}
