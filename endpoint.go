package bartender

import (
	"fmt"
	"reflect"

	"github.com/go-kit/kit/endpoint"
)

func isFunc(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Func
}

// Endpoint wraps a provided function and returns
// a valid go-kit Endpoint
func Endpoint(fn interface{}) (e endpoint.Endpoint, err error) {
	if !isFunc(fn) {
		err = fmt.Errorf("input is not a function: %#v", fn)
		return
	}
	fnt := reflect.TypeOf(fn)

	if n := fnt.NumIn(); n != 1 {
		err = fmt.Errorf("Given function acceptEndpoint functions accepts only 1 argument. Given function has %d", n)
		// TODO: should allow to take context.Context as input, too
	}

	if n := fnt.NumOut(); n != 2 {
		err = fmt.Errorf("Given function returns too many parameters (%d). Should be 2", n)
		// TODO: might accept lesser return parameter (e.g. no error or not output)
	}

	// TODO: rewrite `pass` into `call`
	// -> to examine if the input values type
	//    and return error if mismatch
	// -> to call fn with input values,
	//    then return the output to the output values
	pass := func(in []reflect.Value) []reflect.Value {

		// temp: construct an empty error variable
		// and return its value
		var errVar error
		err := reflect.ValueOf(&errVar).Elem()

		return []reflect.Value{in[1], err} // pass the request to response
	}

	makeEndpoint := func(fin, fpout interface{}) {
		// fptr is a pointer to a function.
		// Obtain the function value itself (likely nil) as a reflect.Value
		// so that we can query its type and then set the value.
		fn := reflect.ValueOf(fpout).Elem()

		// Make a function of the right type.
		v := reflect.MakeFunc(fn.Type(), pass)

		// Assign it to the value fn represents.
		fn.Set(v)
	}

	makeEndpoint(fn, &e)

	return
}
