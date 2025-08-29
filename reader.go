package csvons

import (
	"encoding/csv"
	"log/slog"
	"os"
)

func read_test() [][]string {
	f, err := os.Open("./testcsv/username.csv")
	if err != nil {
		slog.Error("error opening file", "error", err)
		return nil
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, _ := r.ReadAll()
	return records
}
