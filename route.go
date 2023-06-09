package main

import (
	"reflect"
)

type Route struct {
	URL     string
	Method  string
	Handler *Handler
}

type Handler struct {
	InpType  reflect.Type
	OutType  reflect.Type
	Name     string
	Function interface{}
}

func (r *Route) setMethod(method string) *Route {
	// todo - check if the method is valid string
	r.Method = method
	return r
}
func (r *Route) setURL(url string) *Route {
	// todo - check for any invalid URL
	r.URL = url
	return r
}
