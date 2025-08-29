package csvons

import (
	"testing"
)

func TestRead(t *testing.T) {
	m := read_config_file("./testdata/ruler.json")

	metadata := get_metadata(m)
	t.Log("metadata", metadata)

	exists := get_exists_ruler(m)
	t.Log("exists", exists)

	exists_test(exists, metadata)
}
