# csvons

Use a JSON file (ruler.json) to configure constraints between CSV files

## Features

- [x] Validate that values in a colum exists in a specified column of another file.
- [ ] Ensure that values in a column are unique.
- [ ] Check type and range of cloumn values.

## How to write ruler.json

Apart from csvons_metadata, each key in the ruler.json file represents the stem (base name) of a CSV file, and its value defines the rules (constraints) for that file.

## Structure of metadata

- **csv_file_folder** : The folder that contains the CSV files.
- **name_index**: The row index where the column names are defined in the CSV file.
- **data_index**: The row index where the actual data starts in the CSV file.
- **extension**: The file extension (should be ".csv").

## Structure of `ruler`

- **exists**: An array of rules that specify that the values in a column of this CSV file must also exist in a specified column of another file.
  - **dst_file_stem**: The stem (base name) of the target CSV file.
  - **fields**: A pair of column names to be compared.
    - **src**: The column name in the source file.
    - **dst**: The column name in the target file.

## Cautions

- Use Go's default [CSV library](https://pkg.go.dev/encoding/csv#pkg-overview); it supports only the [RFC4180](https://www.rfc-editor.org/rfc/rfc4180.html) specification.
- The priorities of this library are correctness first, features second, and performance third.
