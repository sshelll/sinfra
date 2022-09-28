package ast

import (
	"go/ast"
	"go/token"
	"strings"
)

func IsGoTestFunc(fn *Func, testingPkgAlias *string) bool {

	if fn == nil || fn.AstDecl == nil {
		return false
	}

	if !strings.HasPrefix(fn.Name, "Test") {
		return false
	}

	fDecl, ok := fn.AstDecl.(*ast.FuncDecl)
	if !ok {
		return false
	}

	argList := fDecl.Type.Params
	if argList == nil || len(argList.List) != 1 || len(argList.List[0].Names) != 1 {
		return false
	}

	starExpr, ok := argList.List[0].Type.(*ast.StarExpr)
	if !ok {
		return false
	}

	selectorExpr, ok := starExpr.X.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	testingPkgName := "testing"
	if testingPkgAlias != nil {
		testingPkgName = *testingPkgAlias
	}
	idt, ok := selectorExpr.X.(*ast.Ident)
	if !ok || idt.Name != testingPkgName {
		return false
	}

	switch selectorExpr.Sel.Name {
	case "T", "B", "M", "F", "TB", "PB":
		return true
	}

	return false

}

func IsTestifySuiteEntryFunc(fn *Func, testingPkgAlias, testifyPkgAlias *string) (suiteName string, ok bool) {

	if !IsGoTestFunc(fn, testingPkgAlias) {
		return
	}

	fDecl := fn.AstDecl.(*ast.FuncDecl)
	if fDecl.Body == nil {
		return
	}

	testifyPkgName := "suite"
	if testifyPkgAlias != nil {
		testifyPkgName = *testifyPkgAlias
	}
	testingT := fDecl.Type.Params.List[0].Names[0].Name

	for _, stmt := range fDecl.Body.List {
		exprStmt, ok := stmt.(*ast.ExprStmt)
		if !ok {
			continue
		}
		if exprStmt.X == nil {
			continue
		}
		callExpr, ok := exprStmt.X.(*ast.CallExpr)
		if !ok {
			continue
		}
		if callExpr.Fun == nil || len(callExpr.Args) != 2 {
			continue
		}
		selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		if selectorExpr.X == nil {
			continue
		}
		if idt, ok := selectorExpr.X.(*ast.Ident); !ok || idt.Name != testifyPkgName {
			continue
		}
		if selectorExpr.Sel == nil || selectorExpr.Sel.Name != "Run" {
			continue
		}
		if len(callExpr.Args) != 2 {
			continue
		}
		idt, ok := callExpr.Args[0].(*ast.Ident)
		if !ok {
			continue
		}
		if idt.Name != testingT {
			continue
		}
		// parse new(suite)
		if callNewExpr, ok := callExpr.Args[1].(*ast.CallExpr); ok {
			idt, ok := callNewExpr.Fun.(*ast.Ident)
			if ok && idt.Name == "new" {
				structIdt, ok := callNewExpr.Args[0].(*ast.Ident)
				if ok {
					return structIdt.Name, true
				}
			}
		}
		// parse &suite{}
		if unaryExpr, ok := callExpr.Args[1].(*ast.UnaryExpr); ok {
			if unaryExpr.Op == token.AND {
				compositeLit, ok := unaryExpr.X.(*ast.CompositeLit)
				if ok {
					idt, isIdt := compositeLit.Type.(*ast.Ident)
					if isIdt {
						return idt.Name, true
					}
				}
			}
		}
	}

	return
}
