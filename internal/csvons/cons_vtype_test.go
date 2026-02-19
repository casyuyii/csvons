package csvons

import (
	"encoding/json"
	"path/filepath"
	"testing"
)

// TestVType validates the "vtype" constraint using the original ruler.json.
func TestVType(t *testing.T) {
	root := projectRoot()
	configFileName := filepath.Join(root, "ruler.json")

	rules, metadata := ReadConfigFile(configFileName)
	if rules == nil || metadata == nil {
		t.Fatalf("read config file error: file_name=%s", configFileName)
	}

	metadata.CSVFileFolder = filepath.Join(root, metadata.CSVFileFolder)

	for stem, v := range rules {
		rulers := map[string]json.RawMessage{}
		if err := json.Unmarshal(v, &rulers); err != nil {
			t.Fatalf("error unmarshalling rulers for %s: %v", stem, err)
		}

		for k, v := range rulers {
			if k == "vtype" {
				var vtype []VType
				if err := json.Unmarshal(v, &vtype); err != nil {
					t.Fatalf("error unmarshalling vtype: %v", err)
				}
				VTypeTest(stem, vtype, metadata)
			}
		}
	}
}

// TestVTypeProducts validates vtype constraints using the products test data.
// Tests float64 (Price), int (Stock), and bool (Available) types.
func TestVTypeProducts(t *testing.T) {
	root := projectRoot()
	configFileName := filepath.Join(root, "testdata", "ruler_products.json")

	rules, metadata := ReadConfigFile(configFileName)
	if rules == nil || metadata == nil {
		t.Fatalf("read config file error: file_name=%s", configFileName)
	}

	metadata.CSVFileFolder = filepath.Join(root, metadata.CSVFileFolder)

	for stem, v := range rules {
		rulers := map[string]json.RawMessage{}
		if err := json.Unmarshal(v, &rulers); err != nil {
			t.Fatalf("error unmarshalling rulers for %s: %v", stem, err)
		}
		for k, v := range rulers {
			if k == "vtype" {
				var vtype []VType
				if err := json.Unmarshal(v, &vtype); err != nil {
					t.Fatalf("error unmarshalling vtype: %v", err)
				}
				VTypeTest(stem, vtype, metadata)
			}
		}
	}
}

// TestVTypeOrders validates vtype constraints using the orders test data.
// Tests nested field expressions (Scores{0}, Scores{1}) with int type.
func TestVTypeOrders(t *testing.T) {
	root := projectRoot()
	configFileName := filepath.Join(root, "testdata", "ruler_orders.json")

	rules, metadata := ReadConfigFile(configFileName)
	if rules == nil || metadata == nil {
		t.Fatalf("read config file error: file_name=%s", configFileName)
	}

	metadata.CSVFileFolder = filepath.Join(root, metadata.CSVFileFolder)

	for stem, v := range rules {
		rulers := map[string]json.RawMessage{}
		if err := json.Unmarshal(v, &rulers); err != nil {
			t.Fatalf("error unmarshalling rulers for %s: %v", stem, err)
		}
		for k, v := range rulers {
			if k == "vtype" {
				var vtype []VType
				if err := json.Unmarshal(v, &vtype); err != nil {
					t.Fatalf("error unmarshalling vtype: %v", err)
				}
				VTypeTest(stem, vtype, metadata)
			}
		}
	}
}

// TestVTypeEmployees validates vtype constraints using the employees test data.
// Tests int (Salary) with range and bool (Active) types.
func TestVTypeEmployees(t *testing.T) {
	root := projectRoot()
	configFileName := filepath.Join(root, "testdata", "ruler_employees.json")

	rules, metadata := ReadConfigFile(configFileName)
	if rules == nil || metadata == nil {
		t.Fatalf("read config file error: file_name=%s", configFileName)
	}

	metadata.CSVFileFolder = filepath.Join(root, metadata.CSVFileFolder)

	for stem, v := range rules {
		rulers := map[string]json.RawMessage{}
		if err := json.Unmarshal(v, &rulers); err != nil {
			t.Fatalf("error unmarshalling rulers for %s: %v", stem, err)
		}
		for k, v := range rulers {
			if k == "vtype" {
				var vtype []VType
				if err := json.Unmarshal(v, &vtype); err != nil {
					t.Fatalf("error unmarshalling vtype: %v", err)
				}
				VTypeTest(stem, vtype, metadata)
			}
		}
	}
}
