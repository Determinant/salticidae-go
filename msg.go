package salticidae

// #include <stdlib.h>
// #include "salticidae/msg.h"
import "C"
import "runtime"

// The C pointer type for a Msg object.
type CMsg = *C.msg_t
type msg struct {
	inner    CMsg
	autoFree bool
	freed    bool
}

//// begin Msg def

// Msg is a message sent by MsgNetwork.
type Msg = *msg

// MsgFromC converts an existing C pointer into a go pointer.
func MsgFromC(ptr CMsg) Msg { return &msg{inner: ptr} }

func msgSetFinalizer(res Msg, autoFree bool) {
	res.autoFree = autoFree
	if res.inner != nil && autoFree {
		runtime.SetFinalizer(res, func(self Msg) { self.Free() })
	}
}

// Free the underlying C pointer manually.
func (msg Msg) Free() {
	if doubleFreeWarn(&msg.freed) {
		return
	}
	C.msg_free(msg.inner)
	if msg.autoFree {
		runtime.SetFinalizer(msg, nil)
	}
}

//// end Msg def

//// begin Msg methods

// NewMsgMovedFromByteArray creates a message by taking out all data from src.
// Notice this is a zero-copy operation that consumes and invalidates the data
// in src ("move" semantics) so that no more operation should be done to src
// after this function call.
func NewMsgMovedFromByteArray(opcode Opcode, src ByteArray, autoFree bool) (res Msg) {
	res = MsgFromC(C.msg_new_moved_from_bytearray(C._opcode_t(opcode), src.inner))
	msgSetFinalizer(res, autoFree)
	runtime.KeepAlive(src)
	return
}

// GetPayloadByMove gets the message payload by taking out all data. Notice
// this is a zero-copy operation that consumes and invalidates the data in the
// payload ("move" semantics) so that no more operation should be done to the
// payload after this function call.
func (msg Msg) GetPayloadByMove() (res DataStream) {
	res = DataStreamFromC(C.msg_consume_payload(msg.inner))
	dataStreamSetFinalizer(res, true)
	runtime.KeepAlive(msg)
	return
}

// GetOpcode gets the opcode.
func (msg Msg) GetOpcode() (res Opcode) {
	res = Opcode(C.msg_get_opcode(msg.inner))
	runtime.KeepAlive(msg)
	return
}

// GetMagic gets the magic number.
func (msg Msg) GetMagic() (res uint32) {
	res = uint32(C.msg_get_magic(msg.inner))
	runtime.KeepAlive(msg)
	return
}

// SetMagic sets the magic number (the default value is 0x0).
func (msg Msg) SetMagic(magic uint32) {
	C.msg_set_magic(msg.inner, C.uint32_t(magic))
	runtime.KeepAlive(msg)
}

//// end Msg methods
