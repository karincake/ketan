// Go Struct Validator Default Validator
package gosv

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// register default tag validator
func init() {
	tagValidator = make(map[string]validator)
	RegisterValidator("required", requiredTagValidator)
	RegisterValidator("eq", eqTagValidator)
	RegisterValidator("similar", eqTagValidator)
	RegisterValidator("min", minTagValidator)
	RegisterValidator("max", maxTagValidator)
	RegisterValidator("minLength", minLengthTagValidator)
	RegisterValidator("maxLength", maxLengthTagValidator)
}

//// for the default validator we have: required + comparisons
func requiredTagValidator(val reflect.Value, exptVal string) error {
	if (val.Kind() == reflect.String && val.String() == "") || (val.Kind() == reflect.Ptr && val.IsNil()) {
		val.Interface()
		return errors.New("field is required")
	}
	return nil
}

func eqTagValidator(val reflect.Value, exptVal string) error {
	return nil
}

func minTagValidator(val reflect.Value, exptVal string) error {
	if err := valLimiter(val, exptVal, "<"); err != nil {
		return err
	}
	return nil
}

func maxTagValidator(val reflect.Value, exptVal string) error {
	if err := valLimiter(val, exptVal, ">"); err != nil {
		return err
	}
	return nil
}

func minLengthTagValidator(val reflect.Value, exptVal string) error {
	exptValInt, err := strconv.Atoi(exptVal)
	if err != nil {
		return errors.New("input must be numeric")
	}

	valC := valStringer(val) // value converted
	if len(valC) < exptValInt {
		return fmt.Errorf("the minimum length is %v", exptVal)
	}
	return nil
}

func maxLengthTagValidator(val reflect.Value, exptVal string) error {
	exptValInt, err := strconv.Atoi(exptVal)
	if err != nil {
		return errors.New("input must be numeric")
	}

	valC := valStringer(val) // value converted
	if len(valC) > exptValInt {
		return fmt.Errorf("the maximum length is %v", exptVal)
	}
	return nil
}

///// some helper for default validator func
func valLimiter(val reflect.Value, exptVal string, mode string) error {
	exptValFloat, err := strconv.ParseFloat(exptVal, 64)
	if err != nil {
		return err
	}

	valC := 0.0 // converted value
	valK := val.Kind()
	if valK == reflect.String {
		valCT, err := strconv.ParseFloat(val.String(), 64)
		if err != nil {
			return errors.New("field must be numeric")
		}
		valC = valCT
	} else if valK >= reflect.Int && valK <= reflect.Uint64 {
		valC = float64(val.Int())
	} else if valK <= reflect.Float32 && valK <= reflect.Float64 {
		valC = val.Float()
	}

	if mode == "<" {
		if exptValFloat > valC {
			return fmt.Errorf("minimum value is %v", exptVal)
		}
	} else {
		if exptValFloat < valC {
			return fmt.Errorf(fmt.Sprintf("maximum value is %v", exptVal))
		}
	}
	return nil
}

func valStringer(val reflect.Value) string {
	valK := val.Kind()
	var valC string
	if valK == reflect.String {
		valC = val.String()
	} else if valK >= reflect.Int && valK < reflect.Uint64 {
		valC = strconv.Itoa(int(val.Int()))
	} else if valK >= reflect.Float32 && valK < reflect.Float64 {
		valC = fmt.Sprintf("%v", val.Float())
	}
	return valC
}
