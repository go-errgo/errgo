// Copyright 2014 Roger Peppe.
// See LICENCE file for details.

package errors_test

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"
	"testing"

	"gopkg.in/errgo.v2/fmt/errors"
)

func TestNewf(t *testing.T) {
	err := errors.Newf("foo %d", 99) //err TestNewf
	checkErr(t, err, nil, "foo 99", `[
	{$TestNewf$: foo 99}
]`, err)
}

func TestNotef(t *testing.T) {
	err0 := errors.Because(nil, someErr, "foo")  //err TestNotef#0
	err := errors.Notef(err0, nil, "bar %d", 99) //err TestNotef#1
	checkErr(t, err, err0, "bar 99: foo", `[
	{$TestNotef#1$: bar 99}
	{$TestNotef#0$: foo}
]`, err)
}

func TestBecausef(t *testing.T) {
	if got := errors.Cause(someErr); got != someErr {
		t.Fatalf("Cause(%#v); got %#v want %#v", someErr, got, someErr)
	}

	causeErr := errors.New("cause error")
	underlyingErr := errors.New("underlying error")               //err TestBecausef#1
	err := errors.Becausef(underlyingErr, causeErr, "foo %d", 99) //err TestBecausef#2
	if got := errors.Cause(err); got != causeErr {
		t.Fatalf("Cause(%#v); got %#v want %#v", err, got, causeErr)
	}

	checkErr(t, err, underlyingErr, "foo 99: underlying error", `[
	{$TestBecausef#2$: foo 99}
	{$TestBecausef#1$: underlying error}
]`, causeErr)

	err = customError{err}
	if got := errors.Cause(err); got != causeErr {
		t.Fatalf("Cause(%#v); got %#v want %#v", err, got, causeErr)
	}
}

// All the code from here on is identical to the code
// in ../../errors/errors_test.go.

func TestNew(t *testing.T) {
	err := errors.New("foo") //err TestNew
	checkErr(t, err, nil, "foo", `[
	{$TestNew$: foo}
]`, err)
}

var someErr = errors.New("some error") //err varSomeErr

func annotate1() error {
	err := errors.Note(someErr, nil, "annotate1") //err annotate1
	return err
}

func annotate2() error {
	err := annotate1()
	err = errors.Note(err, nil, "annotate2") //err annotate2
	return err
}

func TestNoteUsage(t *testing.T) {
	err0 := annotate2()
	err, ok := err0.(errors.Wrapper)
	if !ok {
		t.Fatalf("%#v does not implement errors.Wrapper", err0)
	}
	underlying := err.Underlying()
	checkErr(
		t, err0, underlying,
		"annotate2: annotate1: some error",
		`[
	{$annotate2$: annotate2}
	{$annotate1$: annotate1}
	{$varSomeErr$: some error}
]`,
		err0)
}

func TestWrap(t *testing.T) {
	err0 := errors.Because(nil, someErr, "foo") //err TestWrap#0
	err := errors.Wrap(err0)                    //err TestWrap#1
	checkErr(t, err, err0, "foo", `[
	{$TestWrap#1$: }
	{$TestWrap#0$: foo}
]`, err)

	err = errors.Wrap(nil)
	if err != nil {
		t.Fatalf("Wrap(nil); got %#v want nil", err)
	}
}

func TestNoteWithNilError(t *testing.T) {
	if got := errors.Note(nil, nil, "annotation"); got != nil {
		t.Fatalf("note of nil; got %#v want nil", got)
	}
}

func TestNote(t *testing.T) {
	err0 := errors.Because(nil, someErr, "foo") //err TestNote#0
	err := errors.Note(err0, nil, "bar")        //err TestNote#1
	checkErr(t, err, err0, "bar: foo", `[
	{$TestNote#1$: bar}
	{$TestNote#0$: foo}
]`, err)

	err = errors.Note(err0, errors.Is(someErr), "bar") //err TestNote#2
	checkErr(t, err, err0, "bar: foo", `[
	{$TestNote#2$: bar}
	{$TestNote#0$: foo}
]`, someErr)

	err = errors.Note(err0, func(error) bool { return false }, "") //err TestNote#3
	checkErr(t, err, err0, "foo", `[
	{$TestNote#3$: }
	{$TestNote#0$: foo}
]`, err)
}

func TestCause(t *testing.T) {
	if got := errors.Cause(someErr); got != someErr {
		t.Fatalf("Cause(%#v); got %#v want %#v", someErr, got, someErr)
	}

	causeErr := errors.New("cause error")
	underlyingErr := errors.New("underlying error")          //err TestCause#1
	err := errors.Because(underlyingErr, causeErr, "foo 99") //err TestCause#2
	if got := errors.Cause(err); got != causeErr {
		t.Fatalf("Cause(%#v); got %#v want %#v", err, got, causeErr)
	}

	checkErr(t, err, underlyingErr, "foo 99: underlying error", `[
	{$TestCause#2$: foo 99}
	{$TestCause#1$: underlying error}
]`, causeErr)

	err = customError{err}
	if got := errors.Cause(err); got != causeErr {
		t.Fatalf("Cause(%#v); got %#v want %#v", err, got, causeErr)
	}
}

func TestBecauseWithNoMessage(t *testing.T) {
	cause := errors.New("cause")
	err := errors.Because(nil, cause, "")
	if err == nil || err.Error() != "cause" {
		t.Fatalf(`unexpected error; want "cause" got %q`, err)
	}
	if got := errors.Cause(err); got != cause {
		t.Fatalf("Cause(%#v); got %#v want %#v", err, got, cause)
	}
}

func TestBecauseWithUnderlyingButNoMessage(t *testing.T) {
	err := errors.New("something")
	cause := errors.New("cause")
	err = errors.Because(err, cause, "")
	if err == nil || err.Error() != "something" {
		t.Fatalf(`unexpected error; want "cause" got %q`, err)
	}
	if got := errors.Cause(err); got != cause {
		t.Fatalf("Cause(%#v); got %#v want %#v", err, got, cause)
	}
}

func TestBecauseWithAllZeroArgs(t *testing.T) {
	err := errors.Because(nil, nil, "")
	if err != nil {
		t.Fatalf("Because with all zero args; got %#v want nil", err)
	}
}

func TestDetails(t *testing.T) {
	if got, want := errors.Details(nil), "[]"; got != want {
		t.Fatalf("errors.Details(nil); got %q want %q", got, want)
	}

	otherErr := fmt.Errorf("other")
	checkErr(t, otherErr, nil, "other", `[
	{other}
]`, otherErr)

	err0 := customError{errors.New("foo")} //err TestDetails#0
	checkErr(t, err0, nil, "foo", `[
	{$TestDetails#0$: foo}
]`, err0)

	err1 := customError{errors.Note(err0, nil, "bar")} //err TestDetails#1
	checkErr(t, err1, err0, "bar: foo", `[
	{$TestDetails#1$: bar}
	{$TestDetails#0$: foo}
]`, err1)

	err2 := errors.Wrap(err1) //err TestDetails#2
	checkErr(t, err2, err1, "bar: foo", `[
	{$TestDetails#2$: }
	{$TestDetails#1$: bar}
	{$TestDetails#0$: foo}
]`, err2)
}

func TestSetLocation(t *testing.T) {
	err := customNewError() //err TestSetLocation#0
	checkErr(t, err, nil, "custom", `[
	{$TestSetLocation#0$: custom}
]`, err)
}

func customNewError() error {
	err := errors.New("custom")
	errors.SetLocation(err, 1)
	return err
}

func checkErr(t *testing.T, err, underlying error, msg string, details string, cause error) {
	t.Helper()
	if err == nil {
		t.Fatalf("error is unexpectedly nil")
	}
	if got, want := err.Error(), msg; got != want {
		t.Fatalf("unexpected message; got %q want %q", got, want)
	}
	if err, ok := err.(errors.Wrapper); ok {
		if got, want := err.Underlying(), underlying; got != want {
			t.Fatalf("unexpected underlying error; got %#v want %#v", got, want)
		}
	} else {
		if underlying != nil {
			t.Fatalf("underlying error should be nil; got %#v", underlying)
		}
	}
	if got := errors.Cause(err); got != cause {
		t.Fatalf("Cause(%#v); got %#v want %#v", err, got, cause)
	}
	wantDetails := replaceLocations(details)
	if got := errors.Details(err); got != wantDetails {
		t.Fatalf("unexpected details; got %q want %q", got, wantDetails)
	}
}

func replaceLocations(s string) string {
	t := ""
	for {
		i := strings.Index(s, "$")
		if i == -1 {
			break
		}
		t += s[0:i]
		s = s[i+1:]
		i = strings.Index(s, "$")
		if i == -1 {
			panic("no second $")
		}
		file, line := location(s[0:i])
		t += fmt.Sprintf("%s:%d", file, line)
		s = s[i+1:]
	}
	t += s
	return t
}

func location(tag string) (string, int) {
	line, ok := tagToLine[tag]
	if !ok {
		panic(fmt.Errorf("tag %q not found", tag))
	}
	return filename, line
}

var tagToLine = make(map[string]int)
var filename string

func init() {
	data, err := ioutil.ReadFile("errors_test.go")
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		if j := strings.Index(line, "//err "); j >= 0 {
			tagToLine[line[j+len("//err "):]] = i + 1
		}
	}
	_, filename, _, _ = runtime.Caller(0)
}

type customError struct {
	error
}

func (e customError) Location() (string, int) {
	if err, ok := e.error.(errors.Locator); ok {
		return err.Location()
	}
	return "", 0
}

func (e customError) Underlying() error {
	if err, ok := e.error.(errors.Wrapper); ok {
		return err.Underlying()
	}
	return nil
}

func (e customError) Message() string {
	if err, ok := e.error.(errors.Wrapper); ok {
		return err.Message()
	}
	return ""
}

func (e customError) Cause() error {
	if err, ok := e.error.(errors.Causer); ok {
		return err.Cause()
	}
	return nil
}
