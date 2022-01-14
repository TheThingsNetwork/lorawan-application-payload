// Copyright © 2022 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package apppayload

import "encoding/hex"

// InferGNSS uses a predefined set of rules in order to infer
// a GNSS payload from the provided map.
//
// The following keys are checked for hexadecimal payloads:
// - nav
//
// If no GNSS payload can be inferred, this function returns false.
func InferGNSS(m map[string]interface{}) ([]byte, bool) {
	if len(m) == 0 {
		return nil, false
	}

	for _, key := range []string{
		"nav",
	} {
		s, ok := m[key].(string)
		if !ok {
			continue
		}
		b, err := hex.DecodeString(s)
		if err != nil {
			continue
		}
		return b, true
	}

	return nil, false
}
