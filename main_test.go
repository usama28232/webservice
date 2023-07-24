package main

import "testing"

func TestAdd(t *testing.T) {
	i, j := 1, 2
	expect := 3

	got := add(i, j)

	if expect != got {
		t.Errorf("calculation error: expected %v, got %v", expect, got)
	} else {
		t.Logf("calculation result: expected %v, got %v", expect, got)
	}
}
