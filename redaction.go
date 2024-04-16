package redaction

import (
	"github.com/pkg/errors"
	"github.com/weisbartb/rcache"
	"github.com/weisbartb/redact/internal"
	"reflect"
)

var ErrMustBeStruct = errors.New("must be struct or map/slice of structs")

var Methods = map[string]internal.RawMethod{
	"zero":   internal.MethodZero,
	"remove": internal.MethodRemove,
	"star":   internal.MethodStar,
	"redact": internal.MethodRedact,
}

type redactionInstruction struct {
	eval internal.Evaluator
}

func (r redactionInstruction) FieldName(tag string) string {
	// This tells the cache to just use the default field name, this simplifies the tag signature.
	return ""
}

func (r redactionInstruction) TagNamespace() string {
	return "redact"
}

func (r redactionInstruction) Skip(tag string) bool {
	// We are not interested in fields that are not tagged, we do not want to create copies of them.
	if len(tag) == 0 {
		return true
	}
	return false
}

func (r redactionInstruction) GetMetadata(fieldType reflect.Type, tag string) rcache.InstructionSet {
	var resp redactionInstruction
	ris := internal.NewInstructionScanner(tag)
	ris.Scan()
	resp.eval, _ = ris.GetEvaluator(Methods)
	return resp
}

var instructions = rcache.NewCache(redactionInstruction{})

// RedactRecord redacts a given record of T and a list of groups the current context belongs to.
// It will then recursively redact all data that matches the conditions.
// Please note:
// This method creates a new copy of T in most cases (this is intentional rather than using pointer re-assignment).
// This is an intentional decision to not accidentally wipe out non-exported values from the original value.
// Reflection can't interact with those values, and they may be needed by some call later on (such as open handlers).
// The response for this is only intended to be used for output encoding.
func RedactRecord[T any](record T, groups ...string) (T, error) {
	vOf := reflect.ValueOf(record)
	out, err := redactRecord(vOf, groups...)
	if err != nil {
		return record, err
	}
	return out.Interface().(T), nil
}
func redactRecord(vOf reflect.Value, groups ...string) (reflect.Value, error) {
	// Resolve any pointer or interface wrappings to get the underlying type
	tOf := reflect.Indirect(vOf).Type()
	// Get the cached instruction recordset
	var cachedRecord = instructions.GetTypeDataFor(tOf)
	switch tOf.Kind() {
	case reflect.Interface:
		// Interfaces need to be unwrapped to their underlying types
		var out = reflect.New(vOf.Type()).Elem()
		var returnPtr bool
		if vOf.Kind() == reflect.Pointer {
			returnPtr = true
			vOf = vOf.Elem()
		}
		item, err := redactRecord(vOf.Elem(), groups...)
		if err != nil {
			return vOf, err
		}
		out.Set(item)
		if returnPtr {
			return out.Addr(), nil
		}
		return out, nil
	case reflect.Slice:
		var out reflect.Value
		var addr bool
		var addressableSlice reflect.Value
		if vOf.Kind() == reflect.Pointer {
			out = reflect.New(reflect.SliceOf(tOf.Elem())).Elem()
			addr = true
			addressableSlice = out
			vOf = vOf.Elem()
		} else {
			out = reflect.MakeSlice(vOf.Type(), 0, 0)
		}
		for i := 0; i < vOf.Len(); i++ {
			item, err := redactRecord(vOf.Index(i), groups...)
			if err != nil {
				return vOf, err
			}
			out = reflect.Append(out, item)
		}
		if addr {
			addressableSlice.Set(out)
			return addressableSlice.Addr(), nil
		}
		return out, nil
	case reflect.Array:
		var out reflect.Value
		var addr bool
		out = reflect.New(reflect.ArrayOf(vOf.Len(), tOf.Elem())).Elem()
		if vOf.Kind() == reflect.Pointer {
			addr = true
			vOf = vOf.Elem()
		}
		for i := 0; i < vOf.Len(); i++ {
			item, err := redactRecord(vOf.Index(i), groups...)
			if err != nil {
				return vOf, err
			}
			out.Index(i).Set(item)
		}
		if addr {
			return out.Addr(), nil
		}
		return out, nil
	case reflect.Map:
		var out reflect.Value
		var addr bool
		if vOf.Kind() == reflect.Pointer {
			out = reflect.New(reflect.MapOf(tOf.Key(), tOf.Elem())).Elem()
			addr = true
			vOf = vOf.Elem()
			out.Set(reflect.MakeMap(vOf.Type()))
		} else {
			out = reflect.MakeMap(vOf.Type())
		}

		for _, key := range vOf.MapKeys() {
			item, err := redactRecord(vOf.MapIndex(key), groups...)
			if err != nil {
				return vOf, err
			}
			out.SetMapIndex(key, item)
		}
		if addr {
			return out.Addr(), nil
		}
		return out, nil
	case reflect.Struct:
		// pass down
	default:
		return vOf, ErrMustBeStruct
	}
	var addr = vOf.Kind() == reflect.Pointer
	var out = reflect.New(tOf).Elem()
	if addr {
		out.Set(vOf.Elem())
	} else {
		out.Set(vOf)
	}
	for _, field := range cachedRecord.Fields() {
		fieldV := out.Field(field.Idx)
		var typeStack []reflect.Type
		for fieldV.Kind() == reflect.Ptr || fieldV.Kind() == reflect.Interface {
			ogVal := fieldV.Elem()
			typeStack = append(typeStack, fieldV.Type())
			fieldV = reflect.New(fieldV.Elem().Type()).Elem()
			fieldV.Set(ogVal)
		}
		if len(field.Fields()) > 0 {
			item, err := redactRecord(fieldV)
			if err != nil {
				return vOf, err
			}
			fieldV.Set(item)
			continue
		}
		id := field.InstructionData()
		if id.eval != nil {
			if _, err := id.eval(fieldV, groups...); err != nil {
				return vOf, err
			}
			if len(typeStack) > 0 {
				for i := len(typeStack) - 1; i >= 0; i-- {
					tmp := reflect.New(typeStack[i]).Elem()
					if typeStack[i].Kind() == reflect.Pointer {
						tmp.Set(fieldV.Addr())
					} else {
						tmp.Set(fieldV)
					}
					fieldV = tmp
				}
				out.Field(field.Idx).Set(fieldV)
			}
		}
	}
	if addr {
		return out.Addr(), nil
	}
	return out, nil
}
