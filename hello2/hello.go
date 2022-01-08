package main

import (
	"example/user/hello/morestrings"
	"fmt"

	"github.com/google/go-cmp/cmp"
)

func main() {
	fmt.Println("Hello World")
	s := "oG olleH"
	fmt.Println(morestrings.ReverseRunes(s))
	want, got := "Hello World", "Hello Go"
	if diff := cmp.Diff(want, got); diff != "" {
		fmt.Printf("%q[-want] <---> %q[+got]:\n %s",
			want, got, diff)
	}
}
