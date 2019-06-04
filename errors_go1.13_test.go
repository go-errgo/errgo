// +build go1.13

package errgo_test

import (
	"errors"
	"testing"

	qt "github.com/frankban/quicktest"

	"gopkg.in/errgo.v1"
)

func TestUnwrap(t *testing.T) {
	c := qt.New(t)

	causeErr := errgo.New("cause error")
	underlyingErr := errgo.New("underlying error")       //err TestCause#1
	err := errgo.WithCausef(underlyingErr, causeErr, "") //err TestCause#2

	c.Assert(errors.Unwrap(err), qt.Equals, causeErr)
	c.Assert(errors.Unwrap(causeErr), qt.Equals, nil)
}
