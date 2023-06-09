package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
)

// todo: Change when middleware support to be added
type Middleware struct {
}

type Router struct {
	Router     *Router
	Routes     []*Route
	Middleware []*Middleware
}

func (r *Router) Query(url string, fn interface{}) *Route {
	// todo: check if any reserved url like openapi is being used
	// todo: should I return error?
	inpType, outType, name, err := validateAndRetrieveHandlerParamType(fn)
	if err != nil {
		panic(err)
	}

	handler := &Handler{inpType, outType, name, fn}
	// add the logic to validate the function, along with its parameters
	route := &Route{url, http.MethodGet, handler}
	r.Routes = append(r.Routes, route)
	return route
}

func (r *Router) Mutation(url string, fn interface{}) error {
	// todo: check if any reserved url like openapi is being used
	inpType, outType, name, err := validateAndRetrieveHandlerParamType(fn)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// convert to a wrapper
	// newFunc := wrapper(fn)

	handler := &Handler{inpType, outType, name, fn}
	// add the logic to validate the function, along with its parameters
	route := &Route{url, http.MethodPost, handler}
	r.Routes = append(r.Routes, route)
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
		if reflect.TypeOf((*context.Context)(nil)).Elem() == fnType.In(0) {
			return fnType.In(1), nil
		} else {
			return nil, errors.New("Second input parameter of handler function should be of type context")
		}
	} else {
		return nil, errors.New(fmt.Sprintf("Invalid number of input arguments - Expected : 2, Actual : %v", noOfInputParam))
	}

}

func validateAndRetrieveOutputParamType(fnType reflect.Type) (reflect.Type, error) {

	noOfOutputParam := fnType.NumOut()

	if noOfOutputParam == 1 {
		if reflect.TypeOf((*error)(nil)).Elem() == fnType.Out(0) {
			return fnType.Out(0), nil
		}
		return nil, errors.New("If the output of the function has 1 paramteres, it should of type error")
	} else if noOfOutputParam == 2 {
		if reflect.TypeOf((*error)(nil)).Elem() == fnType.Out(1) {
			return fnType.Out(0), nil
		}
		return nil, errors.New("Second output parameter of handler function should be of type error")
	}
	return nil, errors.New(fmt.Sprintf("Invalid number of output arguments of handler function - Expected : 2, Actual : %v", noOfOutputParam))
}

// For now, this wrapper is not needed
// func wrapper(fn interface{}) func(context.Context, interface{}) []interface{} {
// 	fnValue := reflect.ValueOf(fn)

// 	return func(ctx context.Context, i interface{}) []interface{} {
// 		args := []reflect.Value{
// 			reflect.ValueOf(ctx),
// 			reflect.ValueOf(i),
// 		}

// 		result := fnValue.Call(args)

// 		if len(result) == 1 {
// 			err := result[0].Interface()
// 			return []interface{}{
// 				err,
// 			}
// 		} else if len(result) == 2 {
// 			res := result[0].Interface()
// 			err := result[1].Interface()

// 			if err != nil {
// 				return []interface{}{
// 					nil,
// 					err,
// 				}
// 			}

// 			return []interface{}{
// 				res,
// 				nil,
// 			}
// 		}

// 		return []interface{}{
// 			fmt.Errorf("Invalid function signature: expected (context.Context, struct) (struct, error)"),
// 		}
// 	}
// }
