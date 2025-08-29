package csvons

type ConstrainsConfig struct {
	Exists   []Exists `json:"exists"`
	Metadata Metadata `json:"metadata"`
}

type Metadata struct {
	FieldNameIndex int `json:"field_name_index"`
}

type Exists struct {
	Src struct {
		FileName  string `json:"file_name"`
		FieldName string `json:"field_name"`
	} `json:"src"`
	Dst struct {
		FileName  string `json:"file_name"`
		FieldName string `json:"field_name"`
	} `json:"dst"`
}
