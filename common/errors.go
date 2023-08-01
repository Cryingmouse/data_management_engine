package common

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate = validator.New()

type Module int

type Code int

type ExtendedCode int

// The error information for the external only.
type Error struct {
	Module       Module       `validate:"lt=100"`
	Code         Code         `validate:"lt=10000"`
	ExtendedCode ExtendedCode `validate:"lt=10000"`
	Params       []string     /* The parameter for the error string */

	// Binary representation of the error code
	ErrorCode string
}

// Return the ErrorCode string.
func (e *Error) Error() string {
	if e.ErrorCode != "" {
		return e.ErrorCode
	}
	err := validate.Struct(e)
	if err != nil {
		return ""
	}

	moduleStr := fmt.Sprintf("%02d", e.Module)
	codeStr := fmt.Sprintf("%04d", e.Code)
	extendedCodeStr := fmt.Sprintf("%04d", e.ExtendedCode)

	e.ErrorCode = "E" + moduleStr + codeStr + extendedCodeStr
	return e.ErrorCode
}

// Error Module definition
var (
	ErrHost      = Module(1)
	ErrDirectory = Module(2)
	ErrShare     = Module(3)
)

// Error Code definition
var (
	ErrRegister   = Code(1)
	ErrUnregister = Code(2)
	ErrCreate     = Code(3)
	ErrDelete     = Code(4)
)

// Error ExtendedCode definition
var (
	ErrUnkonwn           = ExtendedCode(0)
	ErrInvalideRequest   = ExtendedCode(1)
	ErrAlreadyRegistered = ExtendedCode(2)
)

func setError(module Module, code Code, extendedCode ExtendedCode, params []string) Error {
	return Error{
		Module:       module,
		Code:         code,
		ExtendedCode: extendedCode,
		Params:       params,
	}
}

var (
	ErrHostRegisterUnknown        = setError(ErrHost, ErrRegister, ErrUnkonwn, []string{""})           /* E0100010000 */
	ErrHostRegisterInvalidRequest = setError(ErrHost, ErrRegister, ErrInvalideRequest, []string{""})   /* E0100010001 */
	ErrHostAlreadyRegistered      = setError(ErrHost, ErrRegister, ErrAlreadyRegistered, []string{""}) /* E0100010002 */
	ErrDirectoryCreateUnknown     = setError(ErrDirectory, ErrCreate, ErrUnkonwn, []string{""})        /* E0200030000 */
	ErrShareCreateUnknown         = setError(ErrShare, ErrCreate, ErrUnkonwn, []string{""})            /* E0300030000 */
)
