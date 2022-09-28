package ast

import (
	"fmt"
	"go/ast"

	"github.com/SCU-SJL/sinfra/conv"
)

type File struct {
	AstFile     *ast.File
	Path        string
	ImportList  []*Import
	StructList  []*Struct
	TypeDefList []*TypeDef
	FuncList    []*Func
}

type Import struct {
	AstDecl ast.Decl
	Alias   string
	Pkg     string
}

type TypeDef struct {
	AstDecl    ast.Decl
	Name       string
	ReferName  string
	MethodList []*Method
}

type Struct struct {
	AstDecl    ast.Decl
	Name       string
	FieldList  []*Field
	MethodList []*Method
}

type Field struct {
	AstField *ast.Field
	NameList []string
	TypeName string
}

type Method struct {
	AstDecl   ast.Decl
	Name      string
	TypeName  string
	IsPtrRecv bool
}

type Func struct {
	AstDecl ast.Decl
	Name    string
}

func (f *File) Print() {
	fmt.Println("-----------------------------")
	fmt.Printf("\n********* File Path *********\n\n%s\n", f.Path)
	fmt.Printf("\n********* Import List *********\n\n")
	for _, imp := range f.ImportList {
		imp.Print()
	}
	fmt.Printf("\n********* Func List *********\n\n")
	for _, fn := range f.FuncList {
		fn.Print()
	}
	fmt.Printf("\n********* TypeDef List *********\n\n")
	for _, t := range f.TypeDefList {
		t.Print()
	}
	fmt.Printf("\n******** Struct List ********\n\n")
	for _, s := range f.StructList {
		s.Print()
		fmt.Println()
	}
	fmt.Println("-----------------------------")
}

func (i *Import) Print() {
	if i == nil {
		return
	}
	fmt.Printf("Alias: '%s', Pkg: %s\n", i.Alias, i.Pkg)
}

func (t *TypeDef) Print() {
	if t == nil {
		return
	}
	fmt.Println("TypeDef Name:", t.Name, "Refer to:", t.ReferName)
	fmt.Println("[Methods]:")
	for _, m := range t.MethodList {
		m.Print()
	}
}

func (s *Struct) Print() {
	if s == nil {
		return
	}
	fmt.Println("Struct Name:", s.Name)
	fmt.Println("[Fields]:")
	for _, f := range s.FieldList {
		f.Print()
	}
	fmt.Println("[Methods]:")
	for _, m := range s.MethodList {
		m.Print()
	}
}

func (f *Field) Print() {
	if f == nil {
		return
	}
	fmt.Printf("Field: name: '%v' type: %s\n", conv.StrConcat(",", f.NameList...), f.TypeName)
}

func (m *Method) Print() {
	if m == nil {
		return
	}
	if m.IsPtrRecv {
		fmt.Printf("Method: (recv *%s) %s(...)\n", m.TypeName, m.Name)
	} else {
		fmt.Printf("Method: (recv %s) %s(...)\n", m.TypeName, m.Name)
	}
}

func (fn *Func) Print() {
	if fn == nil {
		return
	}
	fmt.Printf("Func: %s(...) - Normal\n", fn.Name)
}
