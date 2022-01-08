package morestrings

import "testing"

func TestReverseRunes(t *testing.T) {
	data := []struct {
		in, want string
	}{
		{"oG olleH", "Hello Go"},
		{"Hello, 世界", "界世 ,olleH"},
		{"", ""},
	}
	for _, d := range data {
		got := ReverseRunes(d.in)
		if got != d.want {
			t.Errorf("ReverseRunes(%q) == %q, wanted: %q\n", d.in, got, d.want)
		}
	}
}
