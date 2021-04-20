package bcbp

import (
	"fmt"
	"strconv"
	"strings"
)

// DecodeError implements error interface and represents an error decoding a
// Bar Coded Boarding Pass.
type DecodeError struct {
	Type         ErrorType
	BoardingPass string
	pos          int
	Item         string
	got          string
	Detail       string
}

// ErrorType represents the type of error that occurred.
type ErrorType string

const (
	// ErrInvalidDataFormat is used when processing of an item failed due to
	// the format of the data not matching the format specified by the item.
	ErrInvalidDataFormat ErrorType = "ErrInvalidDataFormat"

	// ErrInsufficientData is used when the Bar Coded Boarding Pass data
	// provided is less than 60 characters. The IATA 792 specification states
	// that the length of the mandatory section of a Bar Coded Boarding Pass
	// is 60 characters.
	ErrInsufficientData ErrorType = "ErrInsufficientData"

	// ErrNonASCII is used when the Bar Coded Boarding Pass data provided
	// contains non ASCII characters.
	ErrNonASCII ErrorType = "ErrNonASCII"

	// ErrUnsupportedBoardingPass is used when the Bar Coded Boarding Pass
	// data is not an "M" type boarding pass.
	ErrUnsupportedBoardingPass ErrorType = "ErrUnsupportedBoardingPass"

	// ErrUnexpectedEndOfInput is used when processing an item where the
	// length to process is greater than the remaining length of the Bar Coded
	// Boarding Pass data.
	ErrUnexpectedEndOfInput ErrorType = "ErrUnexpectedEndOfInput"

	// ErrMalformedSpec is used when a sub-section of a boarding pass is
	// expected to be processed, however, the item being processed is not the
	// correct type. This error is indicates an issue introduced by a code
	// change rather than an issue with usage.
	ErrMalformedSpec ErrorType = "ErrMalformedSpec"

	// ErrUnknownData is used when decoding of the boarding pass data has
	// successfully completed, however, there are remaining unprocessed
	// data. This error is more of a warning that indicates that there
	// is extra data in the Bar Coded Boarding Pass.
	ErrUnknownData ErrorType = "ErrUnknownData"
)

// String converts ErrorType into a human readable prettyPrint.
func (et ErrorType) String() string {
	switch et {
	case ErrInvalidDataFormat:
		return "Invalid data format"
	case ErrInsufficientData:
		return "Insufficient data"
	case ErrNonASCII:
		return "Non ASCII"
	case ErrUnsupportedBoardingPass:
		return "Unsupported boarding pass"
	case ErrUnexpectedEndOfInput:
		return "Unexpected end of input"
	case ErrMalformedSpec:
		return "Malformed spec"
	case ErrUnknownData:
		return "Unknown data"
	default:
		panic(fmt.Sprintf("unrecognized decode error type: %q", string(et)))
	}
}

var _ error = &DecodeError{}

var tmpl = `bcbp: %s:
  boarding pass data:
  | %q
  | %s
  |
  = reason: %s
`

// Error returns a pretty printed error report. It provides:
//   - the ErrorType
//   - position in the Bar Coded Boarding Pass data where error occurred
//   - expected value
//   - actual value
//   - detailed reason for error
func (de *DecodeError) Error() string {
	diff := fmt.Sprintf("%s^ got %s", whitespace(de.pos), de.got)
	return fmt.Sprintf(tmpl, de.Type, de.BoardingPass, diff, de.Detail)
}

// whitespace is a function used by DecodeError.Error() to help visualize
// which character generated an error in the Bar Coded Boarding Pass data.
func whitespace(num int) string {
	var sb strings.Builder
	for i := 0; i < num; i++ {
		sb.WriteString(" ")
	}
	return sb.String()
}

// InvalidDataFormat returns a *DecodeError indicating "invalid data format".
// This is used to report that the value for the given item does not match
// the data format as specified by the IATA 792 resolution.
func InvalidDataFormat(bp string, pos int, item item, value string) *DecodeError {
	return &DecodeError{
		Type:         ErrInvalidDataFormat,
		BoardingPass: bp,
		pos:          pos,
		Item:         item.description,
		got:          fmt.Sprintf("%q", value),
		Detail:       fmt.Sprintf("data for %q must be %s", item.description, item.format),
	}
}

// InsufficientData returns a *DecodeError indicating "insufficient data".
// This is used to report that the Bar Coded Boarding Pass data is invalid because
// the IATA 792 resolution states that there must be at least 60 characters.
func InsufficientData(bp string, length int) *DecodeError {
	return &DecodeError{
		Type:         ErrInsufficientData,
		BoardingPass: bp,
		pos:          length,
		got:          fmt.Sprintf("%q character(s)", strconv.Itoa(length)),
		Detail:       "boarding pass data must have at least 60 characters",
	}
}

// NonASCII returns a *DecodeError indicating "non ASCII". This is used to
// report that the Bar Coded Boarding Pass data contains a non ASCII character. This
// is similar to InvalidDataFormat but the error is returned before processing
// any data.
func NonASCII(bp string, pos int, val rune) *DecodeError {
	return &DecodeError{
		Type:         ErrNonASCII,
		BoardingPass: bp,
		pos:          pos,
		got:          fmt.Sprintf("%c", val),
		Detail:       "boarding pass data must contain only ASCII characters",
	}
}

// UnsupportedBoardingPass returns a *DecodeError indicating
// "unsupported boarding pass". This is used to report that the boarding pass
// is invalid because it is not an "M" type boarding pass.
//
// "M" refers to multi-leg boarding pass.
func UnsupportedBoardingPass(bp string, value string) *DecodeError {
	return &DecodeError{
		Type:         ErrUnsupportedBoardingPass,
		BoardingPass: bp,
		pos:          1,
		got:          fmt.Sprintf("%q", value),
		Detail:       `boarding pass must be a "M" type`,
	}
}

// UnexpectedEndOfInput returns a *DecodeError indicating
// "unexpected end of input". This is used to report that an error occurred
// while processing item where the length required to process item is
// greater than the length of the remainder of the Bar Coded Boarding Pass data.
func UnexpectedEndOfInput(bp string, pos int, item item, value string, length int) *DecodeError {
	return &DecodeError{
		Type:         ErrUnexpectedEndOfInput,
		BoardingPass: bp,
		pos:          pos,
		got:          fmt.Sprintf("%q character(s)", strconv.Itoa(len(value))),
		Detail: fmt.Sprintf(
			"%q must have at least %d character(s)",
			item.description,
			length),
	}
}

// MalformedSpec returns a *DecodeError indicating "malformed spec". This is
// used to report that a sub-section of a boarding pass is expected to be
// processed, however, the item being processed is not the correct type.
//
// Sub-sections in a boarding pass are denoted by the following items:
//   Field Size of Variable Size Field
//   Field Size of Following Structured Message (Unique)
//   Field Size of Following Structured Message (Repeated)
//   Length of Security data
func MalformedSpec(bp string, pos int, item item) *DecodeError {
	return &DecodeError{
		Type:         ErrMalformedSpec,
		BoardingPass: bp,
		pos:          pos,
		got:          fmt.Sprintf("%q item defines sub-section", item.description),
		Detail: fmt.Sprintf(
			"only following items can define sub-sections:"+
				"\n\t\t- %q\n\t\t- %q\n\t\t- %q\n\t\t- %q",
			"Field Size of variable size field",
			"Field Size of following structured message - unique",
			"Field Size of following structured message - repeated",
			"Length of Security data"),
	}
}

// UnknownData returns a *DecodeError indicating "unknown data". This is used
// to report that decoding of the boarding pass has successfully completed,
// however, there are remaining unprocessed data. This occurs when the length
// specified by "Length of Security data" is shorter than the remaining
// unprocessed boarding pass data.
func UnknownData(bp string, pos int, value string) *DecodeError {
	return &DecodeError{
		Type:         ErrUnknownData,
		BoardingPass: bp,
		pos:          pos,
		got:          fmt.Sprintf("%q character(s)", strconv.Itoa(len(value))),
		Detail:       fmt.Sprintf("boarding pass successfully decoded but %q is unknown and has not been processed", value),
	}
}
