package main

import (
	"fmt"

	"github.com/karincake/gosv"
)

type myS struct {
	Field1 int    `validate:"min=5,max=10"`
	Field2 string `validate:"required,minLength=5,maxLength=10"`
	Field3 *myS2
	Field4 interface{}
}

type myS2 struct {
	Field1 int    `validate:"min=10,max=20"`
	Field2 string `validate:"required,min=10,max=20"`
	Field3 interface{}
}

func main() {
	myVar1 := myS{
		Field1: 4,
		Field2: "Lorem",
		Field3: &myS2{
			Field1: 1,
			Field2: "Lorem Ipsum",
			Field3: struct {
				Field1 int
				Field2 string
			}{10, "Train 10"}},
		Field4: 124,
	}
	myVar2 := &myVar1
	myVar3 := *myVar2

	myVar3.Field1 = 7
	myVar3.Field3.Field1 = 15
	myVar3.Field3.Field2 = "15"

	fmt.Println("Validating myVar1")
	if err := gosv.Validate(myVar1); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Validating myVar2")
	if err := gosv.Validate(myVar2); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Validating myVar3")
	if err := gosv.Validate(myVar3); err != nil {
		fmt.Println(err)
	}

}
