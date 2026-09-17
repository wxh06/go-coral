// Copyright (c) 2026 Xinhe Wang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package contract

import "testing"

// The numerals go on the wire as they are, as hash_method to a remote signer and as the method
// argument to an in-process one, so they have to be the service's own. A constant inserted into
// the block would renumber the ones after it and fail here rather than at the server.
func TestHashMethodNumerals(t *testing.T) {
	for _, tc := range []struct {
		m    HashMethod
		want int
	}{
		{HashMethodNone, 0},
		{HashMethodAuth, 1},
		{HashMethodWebService, 2},
	} {
		if int(tc.m) != tc.want {
			t.Errorf("%v = %d, want %d", tc.m, int(tc.m), tc.want)
		}
	}
}
