package internal

import (
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

type MemoizedMethod func(value any) error
type RawMethod func(arguments ...Arg) (MemoizedMethod, error)

var ErrValueCanOnlyBeString = errors.New("value can only be a string")
var ErrValueMustBeReference = errors.New("value must be a reference")

func getStringValueOf(value any) (reflect.Value, error) {
	vOf, ok := value.(reflect.Value)
	if !ok {
		vOf = reflect.ValueOf(value)
		if vOf.Kind() != reflect.Ptr {
			return reflect.Value{}, ErrValueMustBeReference
		}
		vOf = vOf.Elem()
		if vOf.Kind() != reflect.String {
			return reflect.Value{}, ErrValueCanOnlyBeString
		}
	}
	return vOf, nil
}

// MethodStar will asterisk (*) characters.
// Takes a single argument (integer) that if positive, it will star the first X characters.
// If negative, it will star the last X characters
// Example 555-555-5555 with remove(-4) would be ********5555
func MethodStar(arguments ...Arg) (MemoizedMethod, error) {
	var offset int
	if len(arguments) > 0 {
		offset = arguments[0].Int()
	}
	return func(value any) error {
		vOf, err := getStringValueOf(value)
		if err != nil {
			return errors.Wrap(err, "in redaction method star")
		}
		if offset >= 0 {
			var i int
			vOf.SetString(strings.Map(func(r rune) rune {
				i++
				if i >= offset+1 {
					return '*'
				}
				return r
			}, vOf.String()))
		} else {
			offset *= -1
			var i = len(vOf.String())
			vOf.SetString(strings.Map(func(r rune) rune {
				i--
				if i-offset < 0 {
					return r
				}
				return '*'
			}, vOf.String()))
		}
		return nil
	}, nil
}

// MethodRemove will remove characters from a string.
// It takes a single argument that if positive, it will remove characters after the offset.
// If the argument is negative, it will remove everything but the last X characters.
// Example 555-555-5555 with remove(-4) would be 5555
func MethodRemove(arguments ...Arg) (MemoizedMethod, error) {
	var offset int
	if len(arguments) > 0 {
		offset = arguments[0].Int()
	}
	return func(value any) error {
		vOf, err := getStringValueOf(value)
		if err != nil {
			return errors.Wrap(err, "in redaction method remove")
		}
		if offset >= 0 {
			vOf.SetString(vOf.String()[0:offset])
		} else {
			offset *= -1
			var i = len(vOf.String())
			vOf.SetString(vOf.String()[i-offset : i])
		}
		return nil
	}, nil
}

// MethodZero zeros out the entry
// Example: 66 with zero() will be 0, "66" with zero() will be ""
func MethodZero(arguments ...Arg) (MemoizedMethod, error) {
	return func(value any) error {
		vOf, ok := value.(reflect.Value)
		if !ok {
			vOf := reflect.ValueOf(value)
			if vOf.Kind() != reflect.Ptr {
				return errors.Wrap(ErrValueMustBeReference, "in redaction method zero")
			}
			vOf = vOf.Elem()
		}
		vOf.Set(reflect.New(vOf.Type()).Elem())
		return nil
	}, nil
}

// MethodRedact allows the redaction of all characters.
// Optionally takes a redaction character as the first argument, defaults to *.
// The second argument is a list of allowed characters (that won't be redacted)
// Example, 555-555-555 using redact("*","-") would result in ***-***-****
func MethodRedact(arguments ...Arg) (MemoizedMethod, error) {
	var allowedChars []byte
	var redactionChar byte
	if len(arguments) == 1 {
		redactionChar = []byte(arguments[0].String())[0]
	} else if len(arguments) >= 2 {
		redactionChar = []byte(arguments[0].String())[0]
		allowedChars = []byte(arguments[1].String())
	}
	if redactionChar == 0 {
		redactionChar = '*'
	}
	return func(value any) error {
		vOf, err := getStringValueOf(value)
		if err != nil {
			return errors.Wrap(err, "in redaction method redact")
		}
		var out = []byte(vOf.String())
	replace:
		for k, v := range out {
			for _, allowedChar := range allowedChars {
				if v == allowedChar {
					continue replace
				}
			}
			out[k] = redactionChar
		}
		vOf.SetString(string(out))
		return nil
	}, nil
}
