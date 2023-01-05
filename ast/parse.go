package ast

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strconv"

	serr "github.com/sshelll/sinfra/error"
)

func Parse(path string) (*File, error) {

	astFile, absPath, err := parseAST(path)
	if err != nil {
		return nil, err
	}

	file := &File{
		Path: absPath,
	}

	typeMethodMap := make(map[string][]*Method)
	importDecls, typeDecls, funcDecls := extractDecls(astFile)

	// build imports
	for _, decl := range importDecls {
		file.ImportList = append(file.ImportList, TryConvDeclToImports(decl)...)
	}

	// build func / methods
	for _, decl := range funcDecls {
		m, f := TryConvDeclToFuncOrMethod(decl)
		if m != nil {
			typeMethodMap[m.TypeName] = append(typeMethodMap[m.TypeName], m)
		}
		if f != nil {
			file.FuncList = append(file.FuncList, f)
		}
	}

	// build struct / types
	for _, decl := range typeDecls {
		t, s := TryConvDeclToTypeOrStruct(decl)
		if t != nil {
			t.MethodList = typeMethodMap[t.Name]
			file.TypeDefList = append(file.TypeDefList, t)
		}
		if s != nil {
			s.MethodList = typeMethodMap[s.Name]
			file.StructList = append(file.StructList, s)
		}
	}

	return file, nil

}

func ExtractImportDecls(astFile *ast.File) (decls []ast.Decl) {
	walkDecls(astFile, func(decl ast.Decl) {
		gDecl, ok := decl.(*ast.GenDecl)
		if ok && gDecl.Tok == token.IMPORT {
			decls = append(decls, decl)
		}
	})
	return
}

func ExtractTypeDecls(astFile *ast.File) (decls []ast.Decl) {
	walkDecls(astFile, func(decl ast.Decl) {
		gDecl, ok := decl.(*ast.GenDecl)
		if ok && gDecl.Tok == token.TYPE {
			decls = append(decls, decl)
		}
	})
	return
}

func ExtractFuncDecls(astFile *ast.File) (decls []ast.Decl) {
	walkDecls(astFile, func(decl ast.Decl) {
		if _, ok := decl.(*ast.FuncDecl); ok {
			decls = append(decls, decl)
		}
	})
	return
}

func TryConvDeclToImports(decl ast.Decl) (importList []*Import) {

	specList := decl.(*ast.GenDecl).Specs
	if len(specList) == 0 {
		return
	}

	importList = make([]*Import, 0, len(specList))

	for i := range specList {

		spec := specList[i]

		impSpec, ok := spec.(*ast.ImportSpec)
		if !ok {
			continue
		}

		imp := &Import{
			AstDecl: decl,
			Pkg:     serr.Drop1(strconv.Unquote(impSpec.Path.Value)).(string),
		}

		if name := impSpec.Name; name != nil {
			imp.Alias = name.Name
		}

		importList = append(importList, imp)

	}

	return

}

func TryConvDeclToFuncOrMethod(decl ast.Decl) (method *Method, fn *Func) {

	fnDecl, ok := decl.(*ast.FuncDecl)
	if !ok {
		return
	}

	// is func
	if fnDecl.Recv == nil {
		fn = &Func{
			AstDecl: fnDecl,
			Name:    fnDecl.Name.Name,
		}
		return
	}

	// is method
	method = &Method{
		AstDecl: fnDecl,
		Name:    fnDecl.Name.Name,
	}
	t := fnDecl.Recv.List[0].Type
	if starExpr, ok := t.(*ast.StarExpr); ok {
		method.IsPtrRecv = true
		method.TypeName = starExpr.X.(*ast.Ident).Name
	} else {
		method.IsPtrRecv = false
		method.TypeName = t.(*ast.Ident).Name
	}

	return

}

func TryConvDeclToTypeOrStruct(decl ast.Decl) (typeDef *TypeDef, struct_ *Struct) {

	gDecl, ok := decl.(*ast.GenDecl)
	if !ok {
		return
	}

	if gDecl.Tok != token.TYPE {
		return
	}

	spec := gDecl.Specs[0].(*ast.TypeSpec)

	if struct_ = TryConvSpecToStruct(spec); struct_ != nil {
		struct_.AstDecl = decl
		return
	}

	if typeDef = TryConvSpecToTypeDef(spec); typeDef != nil {
		typeDef.AstDecl = decl
	}

	return

}

func TryConvSpecToStruct(spec ast.Spec) (struct_ *Struct) {

	tspec, ok := spec.(*ast.TypeSpec)
	if !ok {
		return
	}

	stt, ok := tspec.Type.(*ast.StructType)
	if !ok {
		return
	}

	struct_ = &Struct{
		Name: tspec.Name.Name,
	}

	for i := range stt.Fields.List {

		field := stt.Fields.List[i]

		field_ := &Field{
			AstField: field,
		}

		// extract field name
		for _, name := range field.Names {
			field_.NameList = append(field_.NameList, name.Name)
		}

		// extract field type
		if idt, ok := field.Type.(*ast.Ident); ok { // basic type, such as 'string'
			field_.TypeName = idt.Name
		}

		// pkg type, such as 'context.Context'
		if expr, ok := field.Type.(*ast.SelectorExpr); ok {
			pkg := expr.X.(*ast.Ident).Name
			clz := expr.Sel.Name
			field_.TypeName = pkg + "." + clz
		}

		struct_.FieldList = append(struct_.FieldList, field_)

	}

	return

}

func TryConvSpecToTypeDef(spec ast.Spec) (typeDef *TypeDef) {

	tspec, ok := spec.(*ast.TypeSpec)
	if !ok {
		return
	}

	idt, ok := tspec.Type.(*ast.Ident)
	if !ok {
		return
	}

	typeDef = &TypeDef{
		Name:      tspec.Name.Name,
		ReferName: idt.Name,
	}

	return

}

func parseAST(path string) (astFile *ast.File, absPath string, err error) {
	fset := token.NewFileSet()
	absPath, err = filepath.Abs(path)
	if err != nil {
		return
	}
	astFile, err = parser.ParseFile(fset, absPath, nil, parser.AllErrors)
	return
}

func extractDecls(astFile *ast.File) (importDecls, typeDecls, funcDecls []ast.Decl) {
	walkDecls(astFile, func(decl ast.Decl) {
		gDecl, ok := decl.(*ast.GenDecl)
		if ok && gDecl.Tok == token.TYPE {
			typeDecls = append(typeDecls, decl)
		}
		if ok && gDecl.Tok == token.IMPORT {
			importDecls = append(importDecls, decl)
		}
		if _, ok := decl.(*ast.FuncDecl); ok {
			funcDecls = append(funcDecls, decl)
		}
	})
	return
}

func walkDecls(astFile *ast.File, fn func(ast.Decl)) {
	if astFile == nil {
		return
	}
	for i := range astFile.Decls {
		decl := astFile.Decls[i]
		fn(decl)
	}
}
