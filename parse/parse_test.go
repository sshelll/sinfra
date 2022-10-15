package parse

import (
	"go/ast"
	"go/parser"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSimpleExpr(t *testing.T) {
	expr := `a == "3" && b == 0`
	printParsedResult(t, expr)
}

func TestParseNestedFieldExpr(t *testing.T) {
	expr := `person.Age == 22 && job.City == "New York" && a == 0.4`
	printParsedResult(t, expr)
}

func TestParseFuncCallExpr(t *testing.T) {
	expr := `a == "3" && b == "0" && in_array(c, []string{"900","1100"})`
	printParsedResult(t, expr)
}

func printParsedResult(t *testing.T, expr string) {
	res, err := parser.ParseExpr(expr)
	assert.Nil(t, err)
	ast.Print(nil, res)
}
