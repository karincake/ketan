/*---------------------------------------------------------------------------*\
IMPORTANT NOTES

This package is simply for my exercise only to improve my skill in using go.
But of course I will try my best to make it perfect.

Also, I am not fan of OOP things, or in go struct-func based thing. Therefore
I made the package exports the main function directly instead of making it part
of a struct.

Regards
\*---------------------------------------------------------------------------*/

// Go Struct Validator
package gosv

import (
	"reflect"
	"strings"

	"golang.org/x/exp/maps"
)

// error format, exported for casting purpose
type ValidationError struct {
	Error      error
	Code       string
	ExptdValue string
	GivenValue interface{}
}

// just key val for the tag
type keyVal struct {
	Key string
	Val string
}

// viladator func interface?
// param: reflect value, string
type validator func(reflect.Value, string) error

// list of validator for the given key from tag
var tagValidator map[string]validator

// tag name to validate
const tagName = "validate"

// parse tag for key - val
func parseTag(tag string) []keyVal {
	kvList := []keyVal{}
	for _, item := range strings.Split(tag, ",") {
		pair := strings.SplitN(strings.TrimSpace(item), "=", 2)
		if len(pair) == 0 {
			continue
		}
		if len(pair) == 1 {
			kvList = append(kvList, keyVal{pair[0], ""})
		}
		if len(pair) == 2 {
			kvList = append(kvList, keyVal{pair[0], pair[1]})
		}
	}
	return kvList
}

// main func, validation of each elem
func Validate(input interface{}, nameSpaces ...string) map[string]ValidationError {
	// identiy value and loop if its pointer until reaches non pointer
	inputV := reflect.ValueOf(input)
	for inputV.Kind() == reflect.Pointer {
		inputV = inputV.Elem()
	}

	// non struct cant be validated
	if inputV.Kind() != reflect.Struct {
		return nil
	}

	// namespace will be availble if it is sub validation
	nameSpace := ""
	if len(nameSpaces) > 0 {
		nameSpace += nameSpaces[0] + "."
	}

	// check each field
	inputT := reflect.TypeOf(inputV.Interface())
	errList := map[string]ValidationError{}
	for i := 0; i < inputV.NumField(); i++ {
		// identify type and value of the field
		fieldT := inputT.Field(i)
		fieldV := inputV.Field(i)
		for fieldV.Kind() == reflect.Pointer {
			fieldV = fieldV.Elem()
		}

		// if current field is struct, validate again
		if fieldT.Type.Kind() == reflect.Struct || fieldT.Anonymous || fieldT.Type.Name() == "" {
			maps.Copy(errList, Validate(fieldV.Interface(), fieldT.Name))
			continue
		}

		tag := fieldT.Tag.Get(tagName)
		if tag != "" {
			parsedTag := parseTag(tag)
			for _, kv := range parsedTag {
				if _, ok := tagValidator[kv.Key]; ok {
					err := tagValidator[kv.Key](fieldV, kv.Val)
					if err != nil {
						errList[nameSpace+fieldT.Name] = ValidationError{err, kv.Key, kv.Val, fieldV.Interface()}
						break // 1 err is enough, break from error check of the current field
					}
				}
			}
		}
	}

	if len(errList) > 0 {
		return errList
	}
	return nil
}

// simply add or remove tag validator
func RegisterValidator(tag string, validatorF validator) {
	tagValidator[tag] = validatorF
}

func UnregisterValidator(tag string) {
	delete(tagValidator, tag)
}
