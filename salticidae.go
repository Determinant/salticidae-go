package salticidae

// #include "salticidae/util.h"
import "C"
import "unsafe"

type rawPtr = unsafe.Pointer

// Opcode is the opcode type.
type Opcode = uint8

// Error is SalticidaeError.
type Error = C.struct_SalticidaeCError

// GetCode returns the error code.
func (err *Error) GetCode() int { return int((*C.struct_SalticidaeCError)(err).code) }

// NewError creates an Error object.
func NewError() Error {
	return C.struct_SalticidaeCError{}
}

// StrError converts the error code into a human readable string.
func StrError(code int) string {
	return C.GoString(C.salticidae_strerror(C.int(code)))
}
