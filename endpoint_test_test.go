package bartender

import (
	"testing"
)

func Test_f1(t *testing.T) {

	if i, err := f1("add", 3, 2); i != 5 {
		t.Errorf("Failed to add. Expected %d but get %d", 5, i)
		if err != nil {
			t.Error("No error message")
		} else {
			t.Errorf("Error message: %s", err.Error())
		}
	}

	if i, err := f1("sub", 3, 2); i != 1 {
		t.Errorf("Failed to subtract. Expected %d but get %d", 1, i)
		if err != nil {
			t.Error("No error message")
		} else {
			t.Errorf("Error message: %s", err.Error())
		}
	}

	if _, err := f1("wat", 3, 2); err == nil {
		t.Error("Expect getting error but not")
	}
}

func Test_f2(t *testing.T) {

	if i, err := f2(5)(2); i != 32 {
		t.Errorf("Failed to f2(5)(2) %d but get %d", 32, i)
		if err != nil {
			t.Error("No error message")
		} else {
			t.Errorf("Error message: %s", err.Error())
		}
	}

	if i, err := f2(2)(3); i != 9 {
		t.Errorf("Failed to f2(2)(3). Expected %d but get %d", 9, i)
		if err != nil {
			t.Error("No error message")
		} else {
			t.Errorf("Error message: %s", err.Error())
		}
	}

	if _, err := f2(-1)(3); err == nil {
		t.Error("Expect getting error but not")
	}

}
