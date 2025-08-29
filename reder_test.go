package csvons

import (
	"testing"
)

func TestRead(t *testing.T) {
	records := read_test()
	t.Log(records)
}
