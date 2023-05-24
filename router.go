package main

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"
)

type HandlerFunction func(interface{}, context.Context) (interface{}, error)

type Handler struct {
	inpType  reflect.Type
	outType  reflect.Type
	name     string
	function HandlerFunction
}
type Route struct {
	pattern string
	handler *Handler
}

type Router struct {
	routes []*Route
	// will it have nested Routers ? For Sub-routing ?
}

func wrapper(fn interface{}) func(interface{}, context.Context) (interface{}, error) {
	fnValue := reflect.ValueOf(fn)

	return func(i interface{}, ctx context.Context) (interface{}, error) {
		args := []reflect.Value{
			reflect.ValueOf(i),
			reflect.ValueOf(ctx),
		}

		result := fnValue.Call(args)

		if len(result) == 2 {
			res := result[0].Interface()
			err := result[1].Interface()

			if err != nil {
				return nil, err.(error)
			}

			return res, nil
		}

		return nil, fmt.Errorf("Invalid function signature: expected (struct, context.Context) (struct, error)")
	}
}

func (r *Router) Register(pattern string, fn interface{}) error {

	inpType, outType, name, err := validateAndRetrieveHandlerParamType(fn)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// convert to a wrapper
	newFunc := wrapper(fn)

	handler := &Handler{inpType, outType, name, newFunc}
	// add the logic to validate the function, along with its parameters
	route := &Route{pattern, handler}
	r.routes = append(r.routes, route)
	fmt.Println("Added 1 route")
	return nil
}

func validateAndRetrieveHandlerParamType(fn interface{}) (reflect.Type, reflect.Type, string, error) {

	fnValue := reflect.ValueOf(fn)

	// check if the type if Func
	if fnValue.Kind() != reflect.Func {
		return nil, nil, "", errors.New("Provided handler is not a function")
	}

	// get the name of the function
	name := runtime.FuncForPC(fnValue.Pointer()).Name()

	// get the function's signature
	fnType := fnValue.Type()

	inpType, err := validateAndRetrieveInputParamType(fnType)
	if err != nil {
		fmt.Println(err)
		return nil, nil, "", err
	}

	outType, err := validateAndRetrieveOutputParamType(fnType)
	if err != nil {
		return nil, nil, "", err
	}
	return inpType, outType, name, nil
}
func validateAndRetrieveInputParamType(fnType reflect.Type) (reflect.Type, error) {

	noOfInputParam := fnType.NumIn()

	if noOfInputParam == 2 {
		if reflect.TypeOf((*context.Context)(nil)).Elem() == fnType.In(1) {
			return fnType.In(0), nil
		} else {
			return nil, errors.New("Second input parameter of handler function should be of type context")
		}
	} else {
		return nil, errors.New(fmt.Sprintf("Invalid number of input arguments - Expected : 2, Actual : %v", noOfInputParam))
	}

}

func validateAndRetrieveOutputParamType(fnType reflect.Type) (reflect.Type, error) {

	noOfOutputParam := fnType.NumOut()

	if noOfOutputParam == 2 {
		if reflect.TypeOf((*error)(nil)).Elem() == fnType.Out(1) {
			return fnType.Out(0), nil
		}
		return nil, errors.New("Second output parameter of handler function should be of type error")
	}
	return nil, errors.New(fmt.Sprintf("Invalid number of output arguments of handler function - Expected : 2, Actual : %v", noOfOutputParam))
}
