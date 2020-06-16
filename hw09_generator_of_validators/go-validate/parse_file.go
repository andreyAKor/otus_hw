package main

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

var (
	ErrDeclNotFound = errors.New("top-level declarations not found")
	ErrNameNotFound = errors.New("package name not found")
)

type Type string

const (
	TypeVar   Type = "var"
	TypeArray Type = "array"
)

// Type implements a structure type field data.
type FieldType struct {
	Name string
	Type Type
}

// Field implements a structure field data.
type Field struct {
	Names      []string
	Type       FieldType
	Validators []Validator
}

// Struct implements a structure data.
type Struct struct {
	Name       string
	NameLetter string
	Fields     []Field
}

// File implements a file validate structure.
type File struct {
	PackageName string
	Structs     []Struct
}

// Parsing file as AST structure.
func ParseFile(filename string) (file File, err error) {
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		return
	}

	file, err = prepareFile(*astFile)
	return
}

// Preparing file data.
func prepareFile(astFile ast.File) (file File, err error) {
	if astFile.Decls == nil {
		err = ErrDeclNotFound
		return
	}
	if astFile.Name == nil {
		err = ErrNameNotFound
		return
	}

	file.PackageName = astFile.Name.Name

	baseTypes := searchBaseTypes(astFile.Decls)

	for _, decl := range astFile.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			st, err := prepareStructs(baseTypes, *typeSpec)
			if err != nil {
				return File{}, nil
			}
			if st != nil {
				file.Structs = append(file.Structs, *st)
			}
		}
	}

	return file, nil
}

//nolint:funlen
// Preparing structs data.
func prepareStructs(
	baseTypes map[string]string,
	typeSpec ast.TypeSpec,
) (*Struct, error) {
	structType, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return nil, nil
	}

	if typeSpec.Name == nil {
		return nil, nil
	}

	fields := structType.Fields
	if fields == nil {
		return nil, nil
	}

	list := fields.List
	if list == nil {
		return nil, nil
	}

	structFields := []Field{}

	for _, item := range list {
		if item.Tag == nil || item.Names == nil {
			continue
		}

		fieldType := prepareFieldType(item.Type)
		if fieldType == nil {
			continue
		}

		baseType, ok := baseTypes[fieldType.Name]
		if !ok {
			continue
		}
		fieldType.Name = baseType

		names := prepareStructNames(item.Names)
		if len(names) == 0 {
			continue
		}

		validators, err := ParseTag(fieldType.Name, item.Tag.Value)
		if err != nil {
			return nil, err
		}
		if len(validators) == 0 {
			continue
		}

		structFields = append(structFields, Field{
			Names:      names,
			Type:       *fieldType,
			Validators: validators,
		})
	}

	if len(structFields) == 0 {
		return nil, nil
	}

	st := Struct{
		Name:       typeSpec.Name.Name,
		NameLetter: strings.ToLower(strings.Split(typeSpec.Name.Name, "")[0]),
		Fields:     structFields,
	}

	return &st, nil
}

// Preparing filed type structure data.
func prepareFieldType(tp ast.Expr) *FieldType {
	switch v := tp.(type) {
	case *ast.Ident:
		return &FieldType{
			Name: v.Name,
			Type: TypeVar,
		}
	case *ast.ArrayType:
		ident, ok := v.Elt.(*ast.Ident)
		if !ok {
			return nil
		}
		return &FieldType{
			Name: ident.Name,
			Type: TypeArray,
		}
	}

	return nil
}

// Search for base types for predefined types.
func searchBaseTypes(decls []ast.Decl) map[string]string {
	baseTypes := map[string]string{
		"int":    "int",
		"string": "string",
	}

	// Search defined types
	for _, decl := range decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			typeIdent, ok := typeSpec.Type.(*ast.Ident)
			if !ok {
				continue
			}

			if typeSpec.Name == nil {
				continue
			}

			if _, ok := baseTypes[typeSpec.Name.Name]; !ok {
				baseTypes[typeSpec.Name.Name] = typeIdent.Name
			}
		}
	}

	// Normalize all defined types to basic types
	for tp, base := range baseTypes {
		newBase, ok := baseTypes[base]
		if !ok {
			continue
		}

		baseTypes[tp] = newBase
	}

	return baseTypes
}

// Preparing list names of fields.
func prepareStructNames(itemNames []*ast.Ident) []string {
	names := []string{}
	for _, name := range itemNames {
		if name == nil {
			continue
		}

		names = append(names, name.Name)
	}

	return names
}
