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
	IsExternal bool
	// To present The external error
	Module       Module       `validate:"lt=100"`
	Code         Code         `validate:"lt=10000"`
	ExtendedCode ExtendedCode `validate:"lt=10000"`
	Params       []string     /* The parameter for the error string */
	errorCode    string       /* Binary representation of the error code */

	// To present internal error
	Err error
}

func (e *Error) Error() string {
	if e.IsExternal {
		if e.errorCode != "" {
			return e.errorCode
		}
		err := validate.Struct(e)
		if err != nil {
			return ""
		}

		moduleStr := fmt.Sprintf("%02d", e.Module)
		codeStr := fmt.Sprintf("%04d", e.Code)
		extendedCodeStr := fmt.Sprintf("%04d", e.ExtendedCode)

		e.errorCode = "E" + moduleStr + codeStr + extendedCodeStr
		return e.errorCode
	} else {
		return e.Err.Error()
	}
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
	ErrGet        = Code(5)
)

// Error ExtendedCode definition
var (
	ErrUnkonwn           = ExtendedCode(0)
	ErrInvalideRequest   = ExtendedCode(1)
	ErrAlreadyRegistered = ExtendedCode(2)
)

func setExternalError(module Module, code Code, extendedCode ExtendedCode) *Error {
	return &Error{
		IsExternal:   true,
		Module:       module,
		Code:         code,
		ExtendedCode: extendedCode,
	}
}

var (
	ErrRegisterHostUnknown             = setExternalError(ErrHost, ErrRegister, ErrUnkonwn)           /* E0100010000 */
	ErrRegisterHostInvalidRequest      = setExternalError(ErrHost, ErrRegister, ErrInvalideRequest)   /* E0100010001 */
	ErrHostAlreadyRegistered           = setExternalError(ErrHost, ErrRegister, ErrAlreadyRegistered) /* E0100010002 */
	ErrUnregisterHostUnknown           = setExternalError(ErrHost, ErrUnregister, ErrUnkonwn)         /* E0100020000 */
	ErrUnregisterHostInvalidRequest    = setExternalError(ErrHost, ErrUnregister, ErrInvalideRequest) /* E0100020001 */
	ErrGetRegisteredHost               = setExternalError(ErrHost, ErrGet, ErrUnkonwn)                /* E0100050000 */
	ErrGetRegisteredHostInvalidRequest = setExternalError(ErrHost, ErrGet, ErrInvalideRequest)        /* E0100050001 */
	ErrCreateDirectoryUnknown          = setExternalError(ErrDirectory, ErrCreate, ErrUnkonwn)        /* E0200030000 */
	ErrCreateShareUnknown              = setExternalError(ErrShare, ErrCreate, ErrUnkonwn)            /* E0300030000 */
)
