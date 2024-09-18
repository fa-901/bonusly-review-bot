package main

import "testing"

func TestHello(t *testing.T) {
	got := Hello("yee")
	want := "Hello yee"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
