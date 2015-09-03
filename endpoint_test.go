package bartender

import (
	"fmt"
	"golang.org/x/net/context"
	"testing"
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
