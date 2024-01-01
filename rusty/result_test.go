package rusty

import (
	"errors"
	"testing"
)

func TestOk(t *testing.T) {
	res := OK(1)
	if 1 != res.Unwrap() {
		t.Fatalf("expected %d, got %d", 1, res.Unwrap())
	}
	if res.Ok().IsNone() {
		t.Fatalf("expected %t, got %t", true, res.Ok().IsNone())
	}
	if 1 != res.Ok().Unwrap() {
		t.Fatalf("expected %d, got %d", 1, res.Ok().Unwrap())
	}
	if !res.Err().IsNone() {
		t.Fatalf("expected %t, got %t", true, res.Err().IsNone())
	}
}

func TestErr(t *testing.T) {
	res := Err[int](errors.New("error"))
	if res.IsOk() {
		t.Fatalf("expected %t, got %t", false, res.IsOk())
	}
	if res.Err().IsNone() {
		t.Fatalf("expected %t, got %t", true, res.Err().IsNone())
	}
	if "error" != res.Err().Unwrap().Error() {
		t.Fatalf("expected %s, got %s", "error", res.Err().Unwrap().Error())
	}
	if !res.Ok().IsNone() {
		t.Fatalf("expected %t, got %t", true, res.Ok().IsNone())
	}
}
