package main

import (
	"go/ast"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrepareFile(t *testing.T) {
	t.Run("declarations not found", func(t *testing.T) {
		_, err := prepareFile(ast.File{
			Decls: nil,
		})
		require.EqualError(t, ErrDeclNotFound, err.Error())
	})
	t.Run("package name not found", func(t *testing.T) {
		_, err := prepareFile(ast.File{
			Decls: []ast.Decl{&ast.GenDecl{}},
			Name:  nil,
		})
		require.EqualError(t, ErrNameNotFound, err.Error())
	})
}

func TestPrepareFieldType(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		res := prepareFieldType(nil)
		require.Nil(t, res)
	})
	t.Run("filed type is var", func(t *testing.T) {
		name := "string"
		res := prepareFieldType(&ast.Ident{
			Name: name,
		})
		require.NotNil(t, res)
		require.EqualValues(t, FieldType{
			Name: name,
			Type: TypeVar,
		}, *res)
	})
	t.Run("filed type is array", func(t *testing.T) {
		name := "string"
		res := prepareFieldType(&ast.ArrayType{
			Elt: &ast.Ident{
				Name: name,
			},
		})
		require.NotNil(t, res)
		require.EqualValues(t, FieldType{
			Name: name,
			Type: TypeArray,
		}, *res)
	})
}

func TestSearchBaseTypes(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		res := searchBaseTypes([]ast.Decl{
			&ast.BadDecl{},
		})
		require.EqualValues(t, map[string]string{
			"int":    "int",
			"string": "string",
		}, res)
	})
	t.Run("check search base types", func(t *testing.T) {
		res := searchBaseTypes([]ast.Decl{
			&ast.GenDecl{
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Type: &ast.Ident{
							Name: "string",
						},
						Name: &ast.Ident{
							Name: "Lala",
						},
					},
					&ast.TypeSpec{
						Type: &ast.Ident{
							Name: "Lala",
						},
						Name: &ast.Ident{
							Name: "Lala2",
						},
					},
				},
			},
		})
		require.EqualValues(t, map[string]string{
			"Lala":   "string",
			"Lala2":  "string",
			"int":    "int",
			"string": "string",
		}, res)
	})
}

func TestPrepareStructNames(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		res := prepareStructNames([]*ast.Ident{nil, nil})
		require.ElementsMatch(t, res, []string{})
	})
	t.Run("empty", func(t *testing.T) {
		names := []string{"Lala", "Bebe"}
		res := prepareStructNames([]*ast.Ident{
			&ast.Ident{
				Name: names[0],
			},
			&ast.Ident{
				Name: names[1],
			},
		})
		require.ElementsMatch(t, res, names)
	})
}
