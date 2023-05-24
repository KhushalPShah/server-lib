package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	router *Router
}

func (server *Server) Start() {
	// todo: check if router was nil
	// todo: check if the routes was empty
	// todo: any additional checks needed ?
	// todo: check if the server was already started ?

	// create a new internal go-chi router
	chi_router := chi.NewRouter()

	// read all the routes from the router
	// format of the Post function - Post(pattern string, h http.HandlerFunc)

	// func(ResponseWriter, *Request)

	routes := server.router.routes
	fmt.Printf("%v", routes)
	for _, route := range routes {
		newFunc := func(w http.ResponseWriter, r *http.Request) {
			//todo: Check for the JSON Schema

			// Decode the JSON request body into the struct instance
			structPtr := reflect.New(route.handler.inpType)
			err := json.NewDecoder(r.Body).Decode(structPtr.Interface())
			if err != nil {
				http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
				return
			}
			data := structPtr.Elem().Interface()

			// Fetching the function from the saved Handler
			fn := route.handler.function

			function := reflect.ValueOf(fn)
			args := []reflect.Value{
				reflect.ValueOf(data),
				reflect.ValueOf(context.TODO()),
			}
			result := function.Call(args)
			if len(result) == 2 {
				res := result[0].Interface()
				err := result[1].Interface()

				responseJSON, err := json.Marshal(res)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "Error marshaling response: %v", err)
					return
				}

				// Set the appropriate headers
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)

				// Send the JSON response
				w.Write(responseJSON)

			}
		}
		chi_router.Post(route.pattern, newFunc)
	}

	http.ListenAndServe(":8080", chi_router)
}
