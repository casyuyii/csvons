package csvons

import (
	"strings"
	"testing"
)

type nilChanExpr struct{}

func (n *nilChanExpr) FieldValue(fields []string, records [][]string) <-chan string { return nil }
func (n *nilChanExpr) typeString() string                                           { return "nil" }
func (n *nilChanExpr) Init(metadata *Metadata, expr string)                         {}

func TestRequiredFieldValuesPanicsWhenExprNil(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected panic")
		}
		if !strings.Contains(r.(error).Error(), "field expression [Nope] is nil") {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()

	requiredFieldValues(nil, "Nope", nil, nil)
}

func TestRequiredFieldValuesPanicsWhenChannelNil(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected panic")
		}
		if !strings.Contains(r.(error).Error(), "field expression [Nope] cannot resolve values") {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()

	requiredFieldValues(&nilChanExpr{}, "Nope", nil, nil)
}
