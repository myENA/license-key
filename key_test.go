package lk_test

import (
	"testing"

	lk "github.com/myENA/license-key"
)

const goodKey = "2af1fe29-9d2771ad-6fb1a3b8-39b276bf-4db2ef1d-4a053327"
const badKey = "2af1fe39-9d2771ad-6fb1a3b8-39b276bf-4db2ef1d-4a153327"

// TestNew tests key creation
func TestNew(t *testing.T) {
	k, err := lk.New()
	t.Logf("key: %s", k)
	if err != nil {
		t.Errorf("failed generating key: %s", err)
		t.Fail()
	}
}

// TestParse tests key parsing and validation
func TestParse(t *testing.T) {
		tests := []struct{
			k string
			b bool
		} {
			{ "2af1fe29-9d2771ad-6fb1a3b8-39b276bf-4db2ef1d-4a053327", true },
			{ "3af1fe29-9d2771ad-6fb1a3b8-39b276bf-4db2ef1d-4a053327", false },
			{ "2af1fe29-9d2771ad-6fb1a3b8-39b276bf-4db2ef1d-4a053328", false },
			{ "2af1fe29-9d2771ad-6fb1a3b8-39b276bf-4db2ef1d4a053327", false },
			{ "2af1fe299d2771ad6fb1a3b839b276bf4db2ef1d4a053327", false },
			{"", false},
		}
	    for _, test := range tests {
	    	k, err := lk.Parse(test.k)
			if test.b && (k == nil || err != nil) {
				t.Errorf("failed to test valid key: %s", err)
			}
	    	if !test.b && (k!= nil || err == nil) {
	    		t.Errorf("parsed invalid key without error: %s", test.k)
			}
		}
}
