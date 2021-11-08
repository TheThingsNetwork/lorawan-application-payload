// Copyright © 2021 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package apppayload

import (
	"fmt"
	"math"
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

// Valid returns if the location is valid.
// A location is considered valid if all fields are well defined (not infinity or NaN),
// latitude is between -90 and 90 degrees, and longitude is between -180 and 180 degrees.
func (l Location) Valid() bool {
	undefined := func(x float64) bool { return math.IsInf(x, 0) || math.IsNaN(x) }
	if undefined(l.Latitude) || undefined(l.Longitude) || undefined(l.Altitude) || undefined(l.Accuracy) {
		return false
	}
	bounded := func(x, low, high float64) bool { return low <= x && x <= high }
	return bounded(l.Latitude, -90, 90) && bounded(l.Longitude, -180, 180)
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
// And the following keys for accuracy in metres:
//   - acc
//   - accuracy
//   - hacc
//
// HDOP and satellite count are currently not used.
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
		if !gpsKeyRegexp.MatchString(k) {
			continue
		}
		vm, ok := m[k].(map[string]interface{})
		if !ok {
			continue
		}
		lat, _ := vm["latitude"].(float64)
		lon, _ := vm["longitude"].(float64)
		alt, _ := vm["altitude"].(float64)
		l := Location{
			Latitude:  lat,
			Longitude: lon,
			Altitude:  alt,
		}
		if l.Valid() {
			return l, true
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
		alt, _ := m[pp.altKey].(float64)
		if !hasLat || !hasLon || (lat == 0 && lon == 0) {
			continue
		}
		res, ok = Location{
			Latitude:  lat,
			Longitude: lon,
			Altitude:  alt,
		}, true
		break
	}

	if !ok {
		return
	}

	// Check accuracy in fields.
	for _, ap := range []string{
		"acc",
		"accuracy",
		"hacc", // horizontal accuracy
	} {
		if acc, hasAcc := m[ap].(float64); hasAcc {
			res.Accuracy = acc
			break
		}
	}

	if !res.Valid() {
		return Location{}, false
	}

	return
}
