// Package errors is the same as gopkg.in/errgo.v2/errors except that
// it adds convenience functions that use the fmt package to format
// error messages.
package errors
import "fmt"

// Newf is like New except it formats the message with a fmt
// format specifier.
func Newf(format string, a ...interface{}) error {
	return New(fmt.Sprintf(format, a...))
}

// Notef is like Note except it formats the message with a fmt
// format specifier.
func Notef(err error, shouldPreserveCause func(error) bool, format string, a ...interface{}) error {
	return Note(err, shouldPreserveCause, fmt.Sprintf(format, a...))
}

// Becausef is like Because except it formats the message with a fmt
// format specifier.
func Becausef(err, cause error, format string, a ...interface{}) error {
	return Because(err, cause, fmt.Sprintf(format, a...))
}
