package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrepareImports(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		res := prepareImports(File{})
		require.Equal(t, 0, len(res))
	})
	t.Run("import fmt", func(t *testing.T) {
		res := prepareImports(File{
			Structs: []Struct{
				{
					Fields: []Field{{
						Type: FieldType{
							Type: TypeArray,
						},
					}},
				},
			},
		})
		require.ElementsMatch(t, res, []string{"fmt"})
	})
	t.Run("import regexp", func(t *testing.T) {
		res := prepareImports(File{
			Structs: []Struct{
				{
					Fields: []Field{{
						Validators: []Validator{
							{
								Type: ValidateTyperRgexp,
							},
						},
					}},
				},
			},
		})
		require.ElementsMatch(t, res, []string{"regexp"})
	})
}
