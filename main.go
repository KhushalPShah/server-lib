package main

import (
	"context"
	"fmt"

	"github.com/invopop/jsonschema"
)

type RegistrationFunctionSchema struct {
	InputSchema  *jsonschema.Schema
	OutputSchema *jsonschema.Schema
	FnName       string
}

type RegistrationFunction func(MyRequest, context.Context) (MyResponse, error)

type MyRequest struct {
	A int `json:"A"`
	B int `json:"B"`
}

type MyResponse struct {
	A int
}

func main() {
	r := &Router{routes: []*Route{}}
	r.Register("/hello", add)

	s := &Server{router: r}
	s.Start()
}

func add(req MyRequest, ctx context.Context) (MyResponse, error) {
	fmt.Println("Inside the function")
	return MyResponse{req.A + req.B}, nil
}

// todo: add this functionality as a part of the handler registration process
// func register(fn interface{}) (*RegistrationFunctionSchema, error) {

// 	fnValue := reflect.ValueOf(fn)

// 	// Check if the provided function is a valid function
// 	if fnValue.Kind() != reflect.Func {
// 		return nil, errors.New("Invalid function provided")
// 	}

// 	fnName := runtime.FuncForPC(fnValue.Pointer()).Name()

// 	// Get the function's signature
// 	fnType := fnValue.Type()

// 	// Validate the input and get input struct
// 	inputStruct, err := getInputStruct(fnType)
// 	if err != nil {
// 		return nil, err
// 	}
// 	inputSchema := jsonschema.ReflectFromType(inputStruct)

// 	// validate the output and get output struct
// 	outputStruct, err := getOutputStruct(fnType)
// 	if err != nil {
// 		return nil, err
// 	}
// 	outputSchema := jsonschema.ReflectFromType(outputStruct)

// 	registrationFuncSchema := RegistrationFunctionSchema{inputSchema, outputSchema, fnName}

// 	// print the computed value
// 	data, err := json.MarshalIndent(registrationFuncSchema, "", " ")
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	fmt.Println(string(data))

// 	return &registrationFuncSchema, nil

// }
