package bartender

import (
	"fmt"
	"golang.org/x/net/context"

	"testing"

	"log"
)

// testing type 1
type t1 int

func (v *t1) Get() string {
	return fmt.Sprintf("%d", *v)
}

// testing type 2
type t2 string

func (v *t2) Get() string {
	return string(*v)
}

// testing type 3
type t3 interface {
	Get() string
}

// testing type 4
type t4 int

func Test_canConnect_SameType(t *testing.T) {
	f1 := func() (j t1) {
		return
	}
	f2 := func(i t1) {
		return
	}

	if !canConnect(f1, f2) {
		t.Error("f1 and f2 should be able to connect")
	}

}

func Test_canConnect_SameKindDiffType(t *testing.T) {
	f1 := func() (j t1) {
		return
	}
	f2 := func(i t4) {
		return
	}

	exp := "cannot use type bartender.t1 as type bartender.t4 in argument 1"
	if canConnect(f1, f2) {
		t.Error("f1 and f2 should not be able to connect")
	}

	if err := validConnect(f1, f2); err.Error() != exp {
		t.Errorf("Error message incorrect\nget:\t%s\nexpect:\t%s",
			err.Error(), exp)
	}

}

func Test_canConnect_CorrectInterface(t *testing.T) {
	f1a := func() (j *t1) {
		return
	}
	f1b := func() (j *t2) {
		return
	}
	f2 := func(i t3) {
		return
	}

	if !canConnect(f1a, f2) {
		t.Error("f1a and f2 should be able to connect")
	}

	if !canConnect(f1b, f2) {
		t.Error("f1b and f2 should be able to connect")
	}
}

func Test_canConnect_IncorrectInterface(t *testing.T) {
	f1 := func() (j t2) {
		return
	}
	f2 := func(i t3) {
		return
	}

	// TODO: In the future should add more detailed error message
	//       about the problem implementation details
	exp := "cannot use type bartender.t2 as type bartender.t3 in argument 1:\n" +
		"bartender.t2 does not implement bartender.t3"
	if canConnect(f1, f2) {
		t.Error("f1 and f2 should not be able to connect")
	}

	if err := validConnect(f1, f2); err.Error() != exp {
		t.Errorf("Error message incorrect\nget:\t%s\nexpect:\t%s",
			err.Error(), exp)
	}

}

func TestCompose_normal(t *testing.T) {
	f1 := func(ctx context.Context, i int) string {
		return fmt.Sprintf("%d", i)
	}
	f2 := func(s string) string {
		return s + " is a number"
	}
	f3 := func(s string) (string, error) {
		log.Printf("f3!!!")
		return s, nil
	}

	ep, err := Compose(f1, f2, f3)
	if err != nil {
		t.Errorf("Endpoint compose error: %s", err.Error())
	}

	_, err = ep(nil, 101)

	/*
		resp, err := ep(nil, 101)
		if err != nil {
			t.Errorf("Endpoint execution error: %s", err.Error())
		}

		if resp != "101 is a number" {
			t.Errorf("Output is not expected: %#v", resp)
		}
	*/
}
