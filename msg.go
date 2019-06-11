package salticidae

// #include <stdlib.h>
// #include "salticidae/msg.h"
import "C"

type Msg = *C.struct_msg_t

func NewMsgMovedFromByteArray(opcode Opcode, _moved_payload ByteArray) Msg {
    return C.msg_new_moved_from_bytearray(C._opcode_t(opcode), _moved_payload)
}

func (self Msg) Free() { C.msg_free(self) }

func (self Msg) ConsumePayload() DataStream {
    return C.msg_consume_payload(self)
}

func (self Msg) GetOpcode() Opcode {
    return Opcode(C.msg_get_opcode(self))
}

