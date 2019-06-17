package salticidae

// #include <stdlib.h>
// #include "salticidae/msg.h"
import "C"
import "runtime"

type CMsg = *C.struct_msg_t
type msg struct { inner CMsg }
type Msg = *msg

func MsgFromC(ptr *C.struct_msg_t) Msg { return &msg{ inner: ptr } }

func NewMsgMovedFromByteArray(opcode Opcode, _moved_payload ByteArray) Msg {
    res := &msg{ inner: C.msg_new_moved_from_bytearray(C._opcode_t(opcode), _moved_payload.inner) }
    runtime.SetFinalizer(res, func(self Msg) { self.free() })
    return res
}

func (self Msg) free() { C.msg_free(self.inner) }

func (self Msg) ConsumePayload() DataStream {
    res := &dataStream{ inner: C.msg_consume_payload(self.inner) }
    runtime.SetFinalizer(res, func(self DataStream) { self.free() })
    return res
}

func (self Msg) GetOpcode() Opcode {
    return Opcode(C.msg_get_opcode(self.inner))
}
