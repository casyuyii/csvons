package csvons

import (
	"log/slog"
	"testing"
)

func TestRead(t *testing.T) {
	configFileName := "./testdata/ruler.json"
	m := readConfigFile(configFileName)
	if m == nil {
		slog.Error("read config file error", "file_name", configFileName)
		return
	}

	metadata := getMetadata(m)
	if metadata == nil {
		slog.Error("get metadata error", "file_name", configFileName)
		return
	}

	exists := getExistsRuler(m)
	if exists == nil {
		slog.Error("get exists ruler error", "file_name", configFileName)
		return
	}

	slog.Info("", "metadata", metadata, "exists", exists)

	existsTest(exists, metadata)
}
