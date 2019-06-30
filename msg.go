package salticidae

// #include <stdlib.h>
// #include "salticidae/msg.h"
import "C"
import "runtime"

// The C pointer type for a Msg object
type CMsg = *C.msg_t
type msg struct {
    inner CMsg
    autoFree bool
}
// Message sent by MsgNetwork
type Msg = *msg

// Convert an existing C pointer into a go object. Notice that when the go
// object does *not* own the resource of the C pointer, so it is only valid to
// the extent in which the given C pointer is valid. The C memory will not be
// deallocated when the go object is finalized by GC. This applies to all other
// "FromC" functions.
func MsgFromC(ptr *C.msg_t) Msg { return &msg{ inner: ptr } }

func msgSetFinalizer(res Msg, autoFree bool) {
    res.autoFree = autoFree
    if res != nil && autoFree {
        runtime.SetFinalizer(res, func(self Msg) { self.Free() })
    }
}

// Create a message by taking out all data from src. Notice this is a zero-copy
// operation that consumes and invalidates the data in src ("move" semantics)
// so that no more operation should be done to src after this function call.
func NewMsgMovedFromByteArray(opcode Opcode, src ByteArray, autoFree bool) Msg {
    res := MsgFromC(C.msg_new_moved_from_bytearray(C._opcode_t(opcode), src.inner))
    msgSetFinalizer(res, autoFree)
    return res
}

func (self Msg) Free() {
    C.msg_free(self.inner)
    if self.autoFree {
        runtime.SetFinalizer(self, nil)
    }
}

// Get the message payload by taking out all data. Notice this is a zero-copy
// operation that consumes and invalidates the data in the payload ("move"
// semantics) so that no more operation should be done to the payload after
// this function call.
func (self Msg) GetPayloadByMove() DataStream {
    res := DataStreamFromC(C.msg_consume_payload(self.inner))
    runtime.SetFinalizer(res, func(self DataStream) { self.Free() })
    return res
}

// Get the opcode.
func (self Msg) GetOpcode() Opcode {
    return Opcode(C.msg_get_opcode(self.inner))
}

// Get the magic number.
func (self Msg) GetMagic() uint32 {
    return uint32(C.msg_get_magic(self.inner))
}

// Set the magic number (the default value is 0x0).
func (self Msg) SetMagic(magic uint32) {
    C.msg_set_magic(self.inner, C.uint32_t(magic))
}
