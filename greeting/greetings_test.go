package greetings

import (
	"regexp"
	"testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestHelloName(t *testing.T) {
	name := "Jay"
	want := regexp.MustCompile(`\b` + name + `\b`)
	res, err := Hello(name)
	if !want.MatchString(res) || err != nil {
		t.Fatalf(`Hello("Jay") = %q,%v, want match for %#q, nil`, res, err, want)
	}
}

// TestHelloEmpty calls greetings.Hello with an empty string,
// checking for an error.
func TestHelloEmpty(t *testing.T) {
	res, err := Hello("")
	if res != "" || err == nil {
		t.Fatalf(`Hello("") = %q,%v, want match "", error`, res, err)
	}
}
