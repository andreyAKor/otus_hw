package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseValidator(t *testing.T) {
	t.Run("validator syntax mismatch", func(t *testing.T) {
		t.Run("empty", func(t *testing.T) {
			_, err := parseValidator("", "")
			require.EqualError(t, err, ErrValidatorSyntaxMismatch.Error())
		})
		t.Run("len < 2", func(t *testing.T) {
			_, err := parseValidator("", "1")
			require.EqualError(t, err, ErrValidatorSyntaxMismatch.Error())
		})
		t.Run("len > 2", func(t *testing.T) {
			_, err := parseValidator("", "123")
			require.EqualError(t, err, ErrValidatorSyntaxMismatch.Error())
		})
	})
	t.Run("unknow validator type", func(t *testing.T) {
		_, err := parseValidator("", "lala:123")
		require.EqualError(t, err, ErrUnknowValidatorType.Error())

		t.Run("empty validator", func(t *testing.T) {
			_, err := parseValidator("", ":")
			require.EqualError(t, err, ErrUnknowValidatorType.Error())
		})
		t.Run("empty validator type", func(t *testing.T) {
			_, err := parseValidator("", ":123")
			require.EqualError(t, err, ErrUnknowValidatorType.Error())
		})
	})
	t.Run("normal validator syntax", func(t *testing.T) {
		t.Run("validator is `len:32`", func(t *testing.T) {
			res, err := parseValidator("string", "len:32")
			require.NoError(t, err)
			require.EqualValues(t, Validator{
				Type:  ValidateTypeLen,
				Value: 32,
			}, res)
		})
		t.Run("validator is `regexp:\\d+`", func(t *testing.T) {
			res, err := parseValidator("string", "regexp:\\d+")
			require.NoError(t, err)
			require.EqualValues(t, Validator{
				Type:  ValidateTyperRgexp,
				Value: "\\d+",
			}, res)
		})
		t.Run("validator is `min:10`", func(t *testing.T) {
			res, err := parseValidator("int", "min:10")
			require.NoError(t, err)
			require.EqualValues(t, Validator{
				Type:  ValidateTypeMin,
				Value: 10,
			}, res)
		})
		t.Run("validator is `max:20`", func(t *testing.T) {
			res, err := parseValidator("int", "max:20")
			require.NoError(t, err)
			require.EqualValues(t, Validator{
				Type:  ValidateTypeMax,
				Value: 20,
			}, res)
		})
		t.Run("validator is `in:foo,bar`", func(t *testing.T) {
			res, err := parseValidator("string", "in:foo,bar")
			require.NoError(t, err)
			require.EqualValues(t, Validator{
				Type:  ValidateTypeIn,
				Value: []string{"foo", "bar"},
			}, res)
		})
		t.Run("validator is `in:256,1024` for int", func(t *testing.T) {
			res, err := parseValidator("int", "in:256,1024")
			require.NoError(t, err)
			require.EqualValues(t, Validator{
				Type:  ValidateTypeIn,
				Value: []int{256, 1024},
			}, res)
		})
		t.Run("validator is `in:256,1024` for string", func(t *testing.T) {
			res, err := parseValidator("string", "in:256,1024")
			require.NoError(t, err)
			require.EqualValues(t, Validator{
				Type:  ValidateTypeIn,
				Value: []string{"256", "1024"},
			}, res)
		})
	})

	t.Run("bad validator syntax", func(t *testing.T) {
		t.Run("validator value as space", func(t *testing.T) {
			_, err := parseValidator("string", "len: ")
			require.Error(t, err)
		})
		t.Run("validator is `len:3.2`", func(t *testing.T) {
			_, err := parseValidator("string", "len:3.2")
			require.Error(t, err)
		})
		t.Run("validator is `len: 32`", func(t *testing.T) {
			_, err := parseValidator("string", "len: 32")
			require.Error(t, err)
		})
		t.Run("validator is `in:foo,bar`", func(t *testing.T) {
			_, err := parseValidator("int", "in:foo,bar")
			require.Error(t, err)
		})
	})
}

func TestParseTag(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		_, err := ParseTag("", "")
		require.NoError(t, err)

		t.Run("tag", func(t *testing.T) {
			_, err := ParseTag("", "``")
			require.NoError(t, err)
		})
		t.Run("validate tag", func(t *testing.T) {
			_, err := ParseTag("", "`validate:\"\"`")
			require.NoError(t, err)

		})
		t.Run("multi tags", func(t *testing.T) {
			_, err := ParseTag("", "`validate:\"\" json:\"\"`")
			require.NoError(t, err)
		})
	})

	t.Run("parse tags", func(t *testing.T) {
		t.Run("type string", func(t *testing.T) {
			t.Run("validate is `len:32`", func(t *testing.T) {
				res, err := ParseTag("string", "`validate:\"len:32\"`")
				require.NoError(t, err)
				require.EqualValues(t, []Validator{
					{
						Type:  ValidateTypeLen,
						Value: 32,
					},
				}, res)
			})
			t.Run("validate is `regexp:\\d+`", func(t *testing.T) {
				res, err := ParseTag("string", "`validate:\"regexp:\\d+\"`")
				require.NoError(t, err)
				require.EqualValues(t, []Validator{
					{
						Type:  ValidateTyperRgexp,
						Value: "\\d+",
					},
				}, res)
			})
			t.Run("validate is `in:foo,bar`", func(t *testing.T) {
				res, err := ParseTag("string", "`validate:\"in:foo,bar\"`")
				require.NoError(t, err)
				require.EqualValues(t, []Validator{
					{
						Type:  ValidateTypeIn,
						Value: []string{"foo", "bar"},
					},
				}, res)
			})
			t.Run("validate is `regexp:\\d+|len:20`", func(t *testing.T) {
				res, err := ParseTag("string", "`validate:\"regexp:\\d+|len:20\"`")
				require.NoError(t, err)
				require.EqualValues(t, []Validator{
					{
						Type:  ValidateTyperRgexp,
						Value: "\\d+",
					},
					{
						Type:  ValidateTypeLen,
						Value: 20,
					},
				}, res)
			})
		})
		t.Run("type int", func(t *testing.T) {
			t.Run("validate is `min:10`", func(t *testing.T) {
				res, err := ParseTag("int", "`validate:\"min:10\"`")
				require.NoError(t, err)
				require.EqualValues(t, []Validator{
					{
						Type:  ValidateTypeMin,
						Value: 10,
					},
				}, res)
			})
			t.Run("validate is `max:20`", func(t *testing.T) {
				res, err := ParseTag("int", "`validate:\"max:20\"`")
				require.NoError(t, err)
				require.EqualValues(t, []Validator{
					{
						Type:  ValidateTypeMax,
						Value: 20,
					},
				}, res)
			})
			t.Run("validate is `in:256,1024`", func(t *testing.T) {
				res, err := ParseTag("int", "`validate:\"in:256,1024\"`")
				require.NoError(t, err)
				require.EqualValues(t, []Validator{
					{
						Type:  ValidateTypeIn,
						Value: []int{256, 1024},
					},
				}, res)
			})
			t.Run("validate is `min:0|max:10`", func(t *testing.T) {
				res, err := ParseTag("int", "`validate:\"min:0|max:10\"`")
				require.NoError(t, err)
				require.EqualValues(t, []Validator{
					{
						Type:  ValidateTypeMin,
						Value: 0,
					},
					{
						Type:  ValidateTypeMax,
						Value: 10,
					},
				}, res)
			})
		})
	})
}
