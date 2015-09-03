package bartender

import (
	"fmt"
	"reflect"

	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
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

	// TODO: parse the fn for later `req` checking

	e = func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		resp = 8
		return
	}
	return
}
