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

	// type of request
	reqt := fnt.In(0)

	// need to call fn using fnv.Call
	fnv := reflect.ValueOf(fn)

	// reflect on endpoint
	et := reflect.TypeOf(e)
	respt := et.Out(0)
	//errt := et.Out(1)

	// TODO:
	// -> to examine if the input values type
	//    and return error if mismatch
	callFn := func(in []reflect.Value) []reflect.Value {

		// TODO: test if the input variable
		//       match the type of request

		// set the input variable to typed variable
		reqv := reflect.New(reqt)
		reqv.Elem().Set(in[1].Elem())
		out := fnv.Call([]reflect.Value{reqv.Elem()})

		// cast the output variable's address
		// into an empty interface
		respv := reflect.New(respt)
		respv.Elem().Set(out[0])

		return []reflect.Value{respv.Elem(), out[1]}
	}

	makeEndpoint := func(fin, fpout interface{}) {
		// fptr is a pointer to a function.
		// Obtain the function value itself (likely nil) as a reflect.Value
		// so that we can query its type and then set the value.
		fn := reflect.ValueOf(fpout).Elem()

		// Make a function of the right type.
		v := reflect.MakeFunc(fn.Type(), callFn)

		// Assign it to the value fn represents.
		fn.Set(v)
	}

	makeEndpoint(fn, &e)

	return
}
