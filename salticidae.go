package salticidae


// #cgo CFLAGS: -I${SRCDIR}/salticidae/include/
// #cgo LDFLAGS: ${SRCDIR}/salticidae/libsalticidae.so -Wl,-rpath=${SRCDIR}/salticidae/
import "C"
import "unsafe"

type rawptr_t = unsafe.Pointer
type Opcode = uint8
