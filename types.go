package csvons

// ConstrainsConfig is the configuration for the constraints
// @param Exists the constraints for the exists
// @param Metadata the metadata for the CSV files
type ConstrainsConfig struct {
	Exists   []Exists `json:"exists"`
	Unique   Unique   `json:"unique"`
	Metadata Metadata `json:"csvons_metadata"`
}

// Metadata is the metadata for the CSV files
// @param CSVFileFolder the folder that contains the CSV files
// @param NameIndex the row index where the column names are defined in the CSV file
// @param DataIndex the row index where the actual data starts in the CSV file
// @param Extension the file extension (should be ".csv")
// @example Metadata{CSVFileFolder: "testdata", NameIndex: 0, DataIndex: 1, Extension: ".csv"}
type Metadata struct {
	CSVFileFolder string `json:"csv_file_folder"`
	NameIndex     int    `json:"name_index"`
	DataIndex     int    `json:"data_index"`
	Extension     string `json:"extension"`
}

// Exists is the constraints for the exists
// @param DstFileStem the stem (base name) of the target CSV file
// @param Fields the fields to be compared
// @example Exists{DstFileStem: "username-d1", Fields: []struct {Src string; Dst string}{Src: "Username", Dst: "Username"}}
type Exists struct {
	DstFileStem string `json:"dst_file_stem"`
	Fields      []struct {
		Src string `json:"src"`
		Dst string `json:"dst"`
	} `json:"fields"`
}

// Unique is the constraints for the unique
// @param Fields the fields to be compared
// @example Unique{Fields: []string{"Username"}}
type Unique struct {
	Fields []string `json:"fields"`
}
