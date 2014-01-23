package citadel

import (
	"testing"
)

func TestDial(t *testing.T) {
	cit, err := Dial("localhost")
	if err != nil {
		t.Errorf("Dial: %v\n", err)
	}
	cit.Close()
}
