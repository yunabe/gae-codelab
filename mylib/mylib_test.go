package mylib

import (
	"testing"
)

func TestHelloMessage(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{input: "World", want: "Hello World"},
		{input: "John", want: "Hello John"},
	}
	for _, test := range tests {
		got := GetHelloMessage(test.input)
		if got != test.want {
			t.Errorf("GetHelloMessage(%q) = %q, want %q", test.input, got, test.want)
		}
	}
}
