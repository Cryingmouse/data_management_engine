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
	ErrLocalUser = Module(4)
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
	ErrDirectoryExisted  = ExtendedCode(3)
	ErrConnectedError    = ExtendedCode(4)
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
	ErrRegisterHostUnknown        = setExternalError(ErrHost, ErrRegister, ErrUnkonwn)           /* E0100010000 */
	ErrRegisterHostInvalidRequest = setExternalError(ErrHost, ErrRegister, ErrInvalideRequest)   /* E0100010001 */
	ErrHostAlreadyRegistered      = setExternalError(ErrHost, ErrRegister, ErrAlreadyRegistered) /* E0100010002 */
	ErrRegisterHostConnectedError = setExternalError(ErrHost, ErrRegister, ErrConnectedError)    /* E0100010004 */

	ErrUnregisterHostUnknown           = setExternalError(ErrHost, ErrUnregister, ErrUnkonwn)          /* E0100020000 */
	ErrUnregisterHostNotExisted        = setExternalError(ErrHost, ErrUnregister, ErrInvalideRequest)  /* E0100020001 */
	ErrUnregisterHostDirectoryExisted  = setExternalError(ErrHost, ErrUnregister, ErrDirectoryExisted) /* E0100020003 */
	ErrGetRegisteredHost               = setExternalError(ErrHost, ErrGet, ErrUnkonwn)                 /* E0100050000 */
	ErrGetRegisteredHostInvalidRequest = setExternalError(ErrHost, ErrGet, ErrInvalideRequest)         /* E0100050001 */
	ErrCreateDirectoryUnknown          = setExternalError(ErrDirectory, ErrCreate, ErrUnkonwn)         /* E0200030000 */
	ErrCreateShareUnknown              = setExternalError(ErrShare, ErrCreate, ErrUnkonwn)             /* E0300030000 */
	ErrCreateLocalUserUnknown          = setExternalError(ErrLocalUser, ErrCreate, ErrUnkonwn)         /* E0400030000 */
	ErrCreateLocalUserInvalidRequest   = setExternalError(ErrLocalUser, ErrCreate, ErrInvalideRequest) /* E0400030001 */
	ErrDeleteLocalUserUnknown          = setExternalError(ErrLocalUser, ErrDelete, ErrUnkonwn)         /* E0400040000 */
	ErrDeleteLocalUserInvalidRequest   = setExternalError(ErrLocalUser, ErrDelete, ErrInvalideRequest) /* E0400040001 */
	ErrGetLocalUser                    = setExternalError(ErrLocalUser, ErrGet, ErrUnkonwn)            /* E0400050000 */
	ErrGetLocalUserInvalidRequest      = setExternalError(ErrLocalUser, ErrGet, ErrInvalideRequest)    /* E0400050001 */
)
