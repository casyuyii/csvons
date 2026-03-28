package csvons

import (
	"log"
	"slices"
	"strings"
)

// FieldOccurrence captures a value emitted from a field expression together
// with the 1-based CSV row number it came from.
type FieldOccurrence struct {
	Row   int
	Value string
}

type fieldOccurrenceProvider interface {
	FieldOccurrences(fields []string, records [][]string) <-chan FieldOccurrence
}

func (p *PlainField) FieldOccurrences(fields []string, records [][]string) <-chan FieldOccurrence {
	fieldIndex := slices.Index(fields, p.fieldName)
	if fieldIndex == -1 {
		return nil
	}

	output := make(chan FieldOccurrence, 128)
	go func() {
		defer close(output)
		for i := p.metadata.DataIndex; i < len(records); i++ {
			record := records[i]
			if fieldIndex < len(record) {
				output <- FieldOccurrence{Row: i + 1, Value: record[fieldIndex]}
			} else {
				log.Printf("plain field [%s] not found in record [%d]", p.fieldName, i)
			}
		}
	}()

	return output
}

func (r *RepeatField) FieldOccurrences(fields []string, records [][]string) <-chan FieldOccurrence {
	fieldIndex := slices.Index(fields, r.fieldName)
	if fieldIndex == -1 {
		return nil
	}

	output := make(chan FieldOccurrence, 128)
	go func() {
		defer close(output)
		for i := r.metadata.DataIndex; i < len(records); i++ {
			record := records[i]
			if fieldIndex < len(record) {
				lev1Vals := strings.Split(record[fieldIndex], r.metadata.Lev1Separator)
				for _, lev1Val := range lev1Vals {
					output <- FieldOccurrence{Row: i + 1, Value: lev1Val}
				}
			} else {
				log.Printf("repeat field [%s] not found in record [%d]", r.fieldName, i)
			}
		}
	}()

	return output
}

func (n *NestedField) FieldOccurrences(fields []string, records [][]string) <-chan FieldOccurrence {
	fieldIndex := slices.Index(fields, n.fieldName)
	if fieldIndex == -1 {
		return nil
	}

	output := make(chan FieldOccurrence, 128)
	go func() {
		defer close(output)
		for i := n.metadata.DataIndex; i < len(records); i++ {
			record := records[i]
			if fieldIndex < len(record) {
				lev1Vals := strings.Split(record[fieldIndex], n.metadata.Lev1Separator)
				for _, lev1Val := range lev1Vals {
					lev2Vals := strings.Split(lev1Val, n.metadata.Lev2Separator)
					if n.index < len(lev2Vals) {
						output <- FieldOccurrence{Row: i + 1, Value: lev2Vals[n.index]}
					} else {
						log.Printf("nested field [%s] level 2 value [%s] length [%d] not found in record [%d]", n.fieldName, lev1Val, len(lev2Vals), i)
					}
				}
			} else {
				log.Printf("nested field [%s] not found in record [%d]", n.fieldName, i)
			}
		}
	}()

	return output
}

func (c *ComplexField) FieldOccurrences(fields []string, records [][]string) <-chan FieldOccurrence {
	fieldIndexes := make([]int, len(c.fieldNames))
	for i, fieldName := range c.fieldNames {
		fieldIndexes[i] = slices.Index(fields, fieldName)
		if fieldIndexes[i] == -1 {
			log.Printf("complex field [%s] not found in fields", fieldName)
			return nil
		}
	}

	output := make(chan FieldOccurrence, 128)
	go func() {
		defer close(output)
		for i := c.metadata.DataIndex; i < len(records); i++ {
			record := records[i]
			cpxStr := ""
			for idx, fieldIndex := range fieldIndexes {
				if fieldIndex < len(record) {
					cpxStr += record[fieldIndex] + c.metadata.FieldConnector
				} else {
					log.Printf("complex field [%s] not found in record [%d]", c.fieldNames[idx], i)
					return
				}
			}
			output <- FieldOccurrence{Row: i + 1, Value: cpxStr}
		}
	}()

	return output
}
