package salticidae


// #cgo CFLAGS: -I${SRCDIR}/salticidae/include/
// #cgo LDFLAGS: ${SRCDIR}/salticidae/libsalticidae.so -Wl,-rpath=${SRCDIR}/salticidae/
// #include "salticidae/util.h"
import "C"
import "unsafe"

type rawptr_t = unsafe.Pointer
type Opcode = uint8
type Error = C.struct_SalticidaeCError

func (self *Error) GetCode() int { return int((*C.struct_SalticidaeCError)(self).code) }

func NewError() Error {
    return C.struct_SalticidaeCError {}
}

func StrError(code int) string {
    return C.GoString(C.salticidae_strerror(C.int(code)))
}
