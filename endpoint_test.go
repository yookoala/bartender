package bartender

import (
	"testing"

	"fmt"
	"reflect"

	"golang.org/x/net/context"
)

func f1(op string, i, j int) (k int, err error) {
	if op == "add" {
		k = i + j
		return
	}
	if op == "sub" {
		k = i - j
		return
	}
	err = fmt.Errorf("Unknown operation \"%s\" ", op)
	return
}

func f2(pow int) (fn func(i int) (j int, err error)) {

	return func(i int) (j int, err error) {
		if pow == 0 {
			return
		} else if pow < 0 {
			err = fmt.Errorf("power < 0")
		}

		j = i
		for n := pow; n >= 2; n-- {
			j *= i
		}

		return
	}

}

type testErr int

func (e *testErr) Error() string {
	return "test error"
}

func Test_validEndpoint(t *testing.T) {
	if _, err := validEndpoint(f1); err == nil {
		t.Error("Failed to catch casting error")
	} else {
		t.Log("validEndpoint catches casting error")
	}
}

func TestEndpoint(t *testing.T) {
	var ctx context.Context
	fp, err := Endpoint(f2(2))
	if err != nil {
		t.Errorf("Error running Endpoint(f2(2)). %s", err.Error())
		return
	}

	resp, err := fp(ctx, 3)
	if err != nil {
		t.Errorf("Error running f2(2)(3). %s", err.Error())
	} else if resp != 9 {
		t.Errorf("Endpoint test with f2(2)(3) failed. Expected %d but get %#v", 9, resp)
	} else {
		t.Log("Endpoint test with f2(2)(3) success.")
	}
}

func Test_temp(t *testing.T) {
	var x error
	t.Logf("type: %v", reflect.TypeOf(x))
}

func Test_isTypeError(t *testing.T) {

	var err1 error
	var err2 *testErr
	var i int

	if isTypeError(reflect.TypeOf(&err1).Elem()) {
		t.Log("Passed test 1")
	} else {
		t.Error("Failed to identify error type as error")
	}

	if !isTypeError(reflect.TypeOf(&err2).Elem()) {
		t.Log("Passed test 2")
	} else {
		t.Error("Failed to identify type implements error as not error")
	}

	if !isTypeError(reflect.TypeOf(&i).Elem()) {
		t.Log("Passed test 3")
	} else {
		t.Error("Failed to identify integer as not error")
	}

}

func TestEndpoint_OutputMismatch(t *testing.T) {
	_, err := Endpoint(func(in int) (out1, out2 int) {
		return 1, 2
	})
	if err != nil && err.Error() == "Second return parameter is not error" {
		t.Log("Detects second parameter is not error")
	} else {
		t.Errorf("Failed to detect input function problem")
	}
}
