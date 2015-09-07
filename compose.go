package bartender

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/go-kit/kit/endpoint"

	"log"
)

var errNumArgMismatch = errors.New("number of arguement(s) mismatch")

// validAssign check if type *from can be assigned to *to
// and detect possible error
func validAssign(from, to *reflect.Type, pos int) (err error) {
	if (*to).Kind() == reflect.Interface {
		// if to is an interface, check if from implements to
		if !(*from).Implements(*to) {
			err = fmt.Errorf("cannot use type %s as type %s "+
				"in argument %d:\n%s does not implement %s",
				*from, *to, pos,
				*from, *to)
		}
		return
	} else if *from != *to {
		// simply check if the types are the same
		err = fmt.Errorf("cannot use type %s as type %s "+
			"in argument %d", *from, *to, pos)
		return
	}
	return
}

// validConnect check if output of function from
// could be connected to the function to.
// reports error involved
func validConnect(from, to interface{}) (err error) {

	// assumes from and to are function
	// and will not check
	fn1 := reflect.TypeOf(from)
	fn2 := reflect.TypeOf(to)

	// check number of arguments
	if fn1.NumOut() != fn2.NumIn() {
		err = errNumArgMismatch
		return
	}

	// check each arguments type
	for i := 0; i < fn1.NumOut(); i++ {
		vt1 := fn1.Out(i)
		vt2 := fn2.In(i)
		if err = validAssign(&vt1, &vt2, i+1); err != nil {
			return
		}
	}

	return
}

// canConnect check if output of function from
// could be connected to the function to
func canConnect(from, to interface{}) bool {
	if err := validConnect(from, to); err != nil {
		return false
	}
	return true
}

// assignType assign a given value into a
// new value of given type
func assignType(in reflect.Value, t reflect.Type) (out reflect.Value) {
	out = reflect.New(t)
	out.Elem().Set(in)
	return
}

// Compose takes multiple functions and compose them into a single
// endpoint. The functions will be called in the order of input.
// The rule to valid compose:
//
// 1. The output variables of any functions should match the
//    types and position of arguments of the next function.
//
// 2. The first function should take a context.Context and
//    1 other argument.
//
// 3. The last function should output 1 variable with 1 error
//
func Compose(fns ...interface{}) (e endpoint.Endpoint, err error) {

	// examine all functions and see if they can be composed
	l := len(fns)
	var last interface{}
	for i, fn := range fns {

		if !isFunc(fn) {
			err = fmt.Errorf("arugment %d is not a function", i+1)
			return
		}

		// test the first and the last one
		fnt := reflect.TypeOf(fn)
		if i == 0 {
			if errInner := validEndpointIn(&fnt); errInner != nil {
				err = fmt.Errorf(
					"begining function is not valid:\n%s",
					errInner.Error())
				return
			}
		} else {

			// test the connection
			if errInner := validConnect(last, fn); errInner != nil {
				err = fmt.Errorf(
					"failed to connect function %d and %d:\n%s",
					i, i+1, errInner.Error())
				return
			}

			if i == l-1 {
				if errInner := validEndpointOut(&fnt); errInner != nil {
					err = fmt.Errorf(
						"ending function is not valid:\n%s",
						errInner.Error())
					return
				}
			}

		}

		last = fn
	}

	// reflect all function values
	fnvs := make([]reflect.Value, len(fns))
	for i, fn := range fns {
		fnvs[i] = reflect.ValueOf(fn)
	}

	// handle first fn request
	fn1t := reflect.TypeOf(fns[0])
	numIn := fn1t.NumIn()
	reqt := fn1t.In(numIn - 1)

	// handle last fn response
	et := reflect.TypeOf(e)
	respt := et.Out(0)

	// use MakeFunc to compose an endpoint.Endpoint
	endpointFn := func(in []reflect.Value) []reflect.Value {

		var work []reflect.Value

		// recast request as type of first function request
		in[1] = assignType(in[1].Elem(), reqt).Elem()
		work = in

		for i := range fnvs {
			log.Printf("run here %d", i)

			work = fnvs[i].Call(work)
		}

		// recast last response as interface{}
		work[0] = assignType(work[0], respt).Elem()

		return work
	}

	// generate the endpoint.Endpoint it
	fout := reflect.ValueOf(&e).Elem()
	v := reflect.MakeFunc(fout.Type(), endpointFn)
	fout.Set(v)

	return
}
