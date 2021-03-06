package bartender

import (
	"fmt"
	"reflect"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"

	"log"
)

var errorType reflect.Type

var ctxType reflect.Type

func init() {
	var err error
	var ctx context.Context
	errorType = reflect.TypeOf(&err).Elem()
	ctxType = reflect.TypeOf(&ctx).Elem()
}

func isFunc(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Func
}

func isTypeError(t reflect.Type) bool {
	return t.String() == "error"
}

func validEndpointIn(fnt *reflect.Type) (err error) {
	n := (*fnt).NumIn()
	if n == 0 || n > 2 {
		err = fmt.Errorf("Given function acceptEndpoint functions accepts 1 argument or context + 1 argument.\n"+
			"Given function has %d", n)
		return
	}

	// if there is 2 argument, the first one must be context
	if n == 2 {
		if t := (*fnt).In(0); t != ctxType {
			log.Printf("t = %#v", t)
			log.Printf("ctxType = %#v", ctxType)
			err = fmt.Errorf("first argument of input must be context.Context. Not %s", t)
		}
	}
	return
}

func validEndpointOut(fnt *reflect.Type) (err error) {
	if n := (*fnt).NumOut(); n != 2 {
		err = fmt.Errorf("Given function returns too many parameters (%d). Should be 2", n)
		// TODO: might accept lesser return parameter (e.g. no error or not output)
		return
	}

	if out2 := (*fnt).Out(1); !isTypeError(out2) {
		err = fmt.Errorf("Second return parameter is not error")
		return
	}
	return
}

// validEndpoint returns the reflect.Type of the given fn
// and error if the fn is not valid for casting as Endpoint
func validEndpoint(fn interface{}) (fnt reflect.Type, err error) {

	if !isFunc(fn) {
		err = fmt.Errorf("input is not a function: %#v", fn)
		return
	}
	fnt = reflect.TypeOf(fn)

	if err = validEndpointIn(&fnt); err != nil {
		return
	}

	if err = validEndpointOut(&fnt); err != nil {
		return
	}

	return
}

// CanEndpoint test if a function can be cast as Endpoint
func CanEndpoint(fn interface{}) bool {
	if _, err := validEndpoint(fn); err != nil {
		return false
	}
	return true
}

// Endpoint wraps a provided function and returns
// a valid go-kit Endpoint
func Endpoint(fn interface{}) (e endpoint.Endpoint, err error) {

	var fnt reflect.Type
	if fnt, err = validEndpoint(fn); err != nil {
		return
	}

	// type of request
	numIn := fnt.NumIn()
	reqt := fnt.In(numIn - 1)

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
		var out []reflect.Value

		if numIn == 1 {
			// single argument
			out = fnv.Call([]reflect.Value{reqv.Elem()})
		} else {
			// context + 1 argument
			out = fnv.Call([]reflect.Value{in[0], reqv.Elem()})
		}

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
