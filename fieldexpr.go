package csvons

import (
	"regexp"
)

type Fieldexpr interface {
	Check(srcRecords [][]string, dstRecords [][]string) bool
	typeString() string
	Init(expr string)
}

// -----------------------------

type PlainField struct {
	SrcFields []string
	DstFields []string
}

func (p *PlainField) Check(srcRecords [][]string, dstRecords [][]string) bool {
	return true
}

func (p *PlainField) typeString() string {
	return "plain"
}

func (p *PlainField) Init(expr string) {
}

type NestedField struct {
	SrcFields []string
	DstFields []string
}

func (n *NestedField) Check(srcRecords [][]string, dstRecords [][]string) bool {
	return true
}

func (n *NestedField) typeString() string {
	return "nested"
}

func (n *NestedField) Init(expr string) {
}

type ComplexField struct {
	SrcFields []string
	DstFields []string
}

func (c *ComplexField) Check(srcRecords [][]string, dstRecords [][]string) bool {
	return true
}

func (c *ComplexField) typeString() string {
	return "complex"
}

func (c *ComplexField) Init(expr string) {
}

// -----------------------------

var fieldexprMap = map[string]func(string) Fieldexpr{
	`^([a-zA-Z0-9]+)$`:      func(string) Fieldexpr { return &PlainField{} },
	`^([a-zA-Z0-9]+)\[\]$`:  func(string) Fieldexpr { return &NestedField{} },
	`^(\{[a-zA-Z0-9]+\})+$`: func(string) Fieldexpr { return &ComplexField{} },
}

func GenerateFieldexpr(fieldExpr string) Fieldexpr {
	for pattern, typ := range fieldexprMap {
		if match, _ := regexp.MatchString(pattern, fieldExpr); match {
			t := typ(fieldExpr)
			t.Init(fieldExpr)
			return t
		}
	}
	return nil
}
