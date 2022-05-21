# go struct validator
A simple validator for struct.

# usage
Import the package then call the function, ie. validate(myStruct). See the example at `./test/main.go`

# func return
The validation return map[string]ValidationError where each key of the map represents a struct field name containing error.

# type validationError
Struct with the following fields:
1. 	Error      error // the error message
2.	Code       string // code of the rule given to the field
3.	ExptdValue string // expected value for the field
4.  GivenValue interface{} // given value for the field
