package csvons

type ConstrainsConfig struct {
	Exists   []Exists `json:"exists"`
	Metadata Metadata `json:"metadata"`
}

type Metadata struct {
	CSVFileFolder  string `json:"csv_file_folder"`
	FieldNameIndex int    `json:"field_name_index"`
	DataIndex      int    `json:"data_index"`
	Extension      string `json:"extension"`
}

type Exists struct {
	DstFileStem string `json:"dst_file_stem"`
	Fields      []struct {
		Src string `json:"src"`
		Dst string `json:"dst"`
	} `json:"fields"`
}
