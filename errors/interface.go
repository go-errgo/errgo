package errors

// Causer is the type of an error that may provide an error cause for
// error diagnosis. Cause may return nil if there is no cause (for
// example because the cause has been masked).
type Causer interface {
	Cause() error
}

// Wrapper is the type of an error that wraps another error. It is
// exposed so that external types may implement it, but should in
// general not be used otherwise.
type Wrapper interface {
	// Message returns the top level error message,
	// not including the message from the underlying
	// error.
	Message() string

	// Underlying returns the underlying error, or nil
	// if there is none.
	Underlying() error
}

// Locationer can be implemented by any error type that wants to expose
// the source location of an error.
type Locationer interface {
	// Location returns the name of the file and the line number
	// associated with an error.
	Location() (file string, line int)
}
