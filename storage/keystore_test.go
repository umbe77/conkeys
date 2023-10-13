package storage

import (
	"testing"
)

func TestCheckString_ok(t *testing.T) {
	testValue := Value{T: String, V: "Hello!"}

	check := testValue.CheckType()
	if !check {
		t.Errorf("%v is not a string", testValue.V)
	}
}

func TestCheckString_ko(t *testing.T) {
	testValue := Value{T: String, V: 12}

	check := testValue.CheckType()
	if check {
		t.Errorf("%v should not be a string", testValue.V)
	}
}

func TestCheckBool_ok(t *testing.T) {
	testValue := Value{T: Boolean, V: true}

	check := testValue.CheckType()
	if !check {
		t.Errorf("%v is not a bool", testValue.V)
	}
}

func TestCheckBool_ko(t *testing.T) {
	testValue := Value{T: Boolean, V: "true"}

	check := testValue.CheckType()
	if check {
		t.Errorf("%v should not be a bool", testValue.V)
	}
}

func TestCheckInteger_ok(t *testing.T) {
	testValue := Value{T: Integer, V: float64(180)}

	check := testValue.CheckType()
	if !check {
		t.Errorf("%v is not an integer", testValue.V)
	}
}

func TestCheckInteger_ko(t *testing.T) {
	testValue := Value{T: Integer, V: true}

	check := testValue.CheckType()
	if check {
		t.Errorf("%v should not be an integer", testValue.V)
	}
}

func TestCheckDecimal_ok(t *testing.T) {
	testValue := Value{T: Decimal, V: 12.234}

	check := testValue.CheckType()
	if !check {
		t.Errorf("%v is not a decimal", testValue.V)
	}
}

func TestCheckDecimal_ko(t *testing.T) {
	testValue := Value{T: Decimal, V: 12}

	check := testValue.CheckType()
	if check {
		t.Errorf("%v should not be a decimal", testValue.V)
	}
}
