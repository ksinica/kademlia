package kademlia

import "testing"

func TestBitMask(t *testing.T) {
	for _, test := range []struct {
		input int
		want  string
	}{
		{input: 0, want: "0"},
		{input: 1, want: "1"},
		{input: 9, want: "1ff"},
		{input: 10, want: "3ff"},
		{input: 11, want: "7ff"},
		{input: 12, want: "fff"},
		{input: 13, want: "1fff"},
		{input: 14, want: "3fff"},
		{input: 15, want: "7fff"},
	} {
		if got := bitMask(test.input).Text(16); got != test.want {
			t.Errorf("expected %v, got %v", test.want, got)
		}
	}
}
