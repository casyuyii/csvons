package csvons

import (
	"log"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type FieldExpr interface {
	FieldValue(fields []string, records [][]string) <-chan string
	typeString() string
	Init(metadata *Metadata, expr string)
}

// -----------------------------
// example: field1, field12
type PlainField struct {
	metadata  *Metadata
	fieldName string
}

func (p *PlainField) FieldValue(fields []string, records [][]string) <-chan string {
	fieldIndex := slices.Index(fields, p.fieldName)
	if fieldIndex == -1 {
		return nil
	}

	output := make(chan string, 128)
	go func() {
		defer close(output)
		for i := p.metadata.DataIndex; i < len(records); i++ {
			record := records[i]
			if fieldIndex < len(record) {
				output <- record[fieldIndex]
			} else {
				log.Printf("plain field [%s] not found in record [%d]", p.fieldName, i)
			}
		}
	}()

	return output
}

func (p *PlainField) typeString() string {
	return "plain"
}

func (p *PlainField) Init(metadata *Metadata, expr string) {
	p.metadata = metadata
	p.fieldName = expr
}

// -----------------------------
// example: field1[], field12[]
type RepeatField struct {
	metadata  *Metadata
	fieldName string
}

func (r *RepeatField) FieldValue(fields []string, records [][]string) <-chan string {
	fieldIndex := slices.Index(fields, r.fieldName)
	if fieldIndex == -1 {
		return nil
	}

	output := make(chan string, 128)
	go func() {
		defer close(output)
		for i := r.metadata.DataIndex; i < len(records); i++ {
			record := records[i]
			if fieldIndex < len(record) {
				lev1Vals := strings.Split(record[fieldIndex], r.metadata.Lev1Separator)
				for _, lev1Val := range lev1Vals {
					output <- lev1Val
				}
			} else {
				log.Printf("repeat field [%s] not found in record [%d]", r.fieldName, i)
			}
		}
	}()

	return output
}

func (r *RepeatField) typeString() string {
	return "repeat"
}

func (r *RepeatField) Init(metadata *Metadata, expr string) {
	r.metadata = metadata
	r.fieldName = expr[:len(expr)-2]
}

// -----------------------------
// example: field1{0}, field12{1}
type NestedField struct {
	metadata  *Metadata
	fieldName string
	index     int
}

func (n *NestedField) FieldValue(fields []string, records [][]string) <-chan string {
	fieldIndex := slices.Index(fields, n.fieldName)
	if fieldIndex == -1 {
		return nil
	}

	output := make(chan string, 128)
	go func() {
		defer close(output)
		for i := n.metadata.DataIndex; i < len(records); i++ {
			record := records[i]
			if fieldIndex < len(record) {
				lev1Vals := strings.Split(record[fieldIndex], n.metadata.Lev1Separator)
				for _, lev1Val := range lev1Vals {
					lev2Vals := strings.Split(lev1Val, n.metadata.Lev2Separator)
					if n.index < len(lev2Vals) {
						output <- lev2Vals[n.index]
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

func (n *NestedField) typeString() string {
	return "nested"
}

func (n *NestedField) Init(metadata *Metadata, expr string) {
	n.metadata = metadata
	matches := regexp.MustCompile(`([a-zA-Z0-9]+)\{(\d+)\}?`).FindStringSubmatch(expr)
	if len(matches) != 3 {
		return
	}

	n.fieldName = matches[1]
	n.index, _ = strconv.Atoi(matches[2])
}

// -----------------------------
// example: {field1}, {field1}{field2}
type ComplexField struct {
	metadata   *Metadata
	fieldNames []string
}

func (c *ComplexField) FieldValue(fields []string, records [][]string) <-chan string {
	fieldIndexes := make([]int, len(c.fieldNames))
	for i, fieldName := range c.fieldNames {
		fieldIndexes[i] = slices.Index(fields, fieldName)
		if fieldIndexes[i] == -1 {
			log.Printf("complex field [%s] not found in fields", fieldName)
			return nil
		}
	}

	output := make(chan string, 128)
	go func() {
		defer close(output)
		for i := c.metadata.DataIndex; i < len(records); i++ {
			cpxStr := ""
			for _, fieldIndex := range fieldIndexes {
				record := records[i]
				if fieldIndex < len(record) {
					cpxStr += record[fieldIndex] + c.metadata.FieldConnector
				} else {
					log.Printf("complex field [%s] not found in record [%d]", c.fieldNames[fieldIndex], i)
					return
				}
			}
			output <- cpxStr
		}
	}()

	return output
}

func (c *ComplexField) typeString() string {
	return "complex"
}

func (c *ComplexField) Init(metadata *Metadata, expr string) {
	c.metadata = metadata
	c.fieldNames = regexp.MustCompile(`([a-zA-Z0-9]+)`).FindAllString(expr, -1)
}

// -----------------------------

var fieldExprMap = map[string]func(string) FieldExpr{
	`^[a-zA-Z0-9]+$`:        func(string) FieldExpr { return &PlainField{} },
	`^[a-zA-Z0-9]+\[\]$`:    func(string) FieldExpr { return &RepeatField{} },
	`^[a-zA-Z0-9]+\{\d+\}$`: func(string) FieldExpr { return &NestedField{} },
	`^\{[a-zA-Z0-9]+\}+$`:   func(string) FieldExpr { return &ComplexField{} },
}

func GenerateFieldExpr(metadata *Metadata, fieldExpr string) FieldExpr {
	if metadata == nil {
		log.Fatal("metadata is nil")
		return nil
	}

	for pattern, makeFunc := range fieldExprMap {
		if match, _ := regexp.MatchString(pattern, fieldExpr); match {
			t := makeFunc(fieldExpr)
			t.Init(metadata, fieldExpr)
			return t
		}
	}

	return nil
}
