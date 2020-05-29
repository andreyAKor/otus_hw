package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrUnknowValidatorType     = errors.New("unknow validator type")
	ErrValidatorSyntaxMismatch = errors.New("validator syntax mismatch")
)

var (
	reParseTag = regexp.MustCompile(`[\s]?validate:"([^\"]+)"[\s]?`)
)

// Validators types.
const (
	ValidateTypeLen    string = "len"
	ValidateTyperRgexp string = "regexp"
	ValidateTypeIn     string = "in"
	ValidateTypeMin    string = "min"
	ValidateTypeMax    string = "max"
)

// Validator implements a validator.
type Validator struct {
	Type  string
	Value interface{}
}

// Parser for struct field tag.
func ParseTag(tp, tag string) ([]Validator, error) {
	list := reParseTag.FindStringSubmatch(tag)
	if len(list) == 0 || len(list[1]) == 0 {
		return []Validator{}, nil
	}

	validators := []Validator{}
	for _, validString := range strings.Split(list[1], "|") {
		validator, err := parseValidator(tp, validString)
		if err != nil {
			return []Validator{}, fmt.Errorf("%s in `%s`", err, validString)
		}
		validators = append(validators, validator)
	}

	return validators, nil
}

// Parsing validator rules.
func parseValidator(tp, validString string) (validator Validator, err error) {
	list := strings.SplitN(validString, ":", 2)
	if len(list) < 2 {
		err = ErrValidatorSyntaxMismatch
		return
	}

	validType := list[0]
	validValue := list[1]

	switch validType {
	case "len":
		validator.Type = ValidateTypeLen
		if validator.Value, err = strconv.Atoi(validValue); err != nil {
			return
		}
	case "regexp":
		validator.Type = ValidateTyperRgexp
		validator.Value = validValue
	case "in":
		validator.Type = ValidateTypeIn
		values := strings.Split(validValue, ",")

		if tp == "int" {
			var (
				ints []int
				i    int
			)
			for _, v := range values {
				i, err = strconv.Atoi(v)
				if err != nil {
					return
				}
				ints = append(ints, i)
			}

			validator.Value = ints
		} else {
			validator.Value = values
		}
	case "min":
		validator.Type = ValidateTypeMin
		if validator.Value, err = strconv.Atoi(validValue); err != nil {
			return
		}
	case "max":
		validator.Type = ValidateTypeMax
		if validator.Value, err = strconv.Atoi(validValue); err != nil {
			return
		}
	default:
		err = ErrUnknowValidatorType
		return
	}

	return
}
