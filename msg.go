package salticidae

// #include <stdlib.h>
// #include "salticidae/msg.h"
import "C"

type Msg = *C.struct_msg_t

func NewMsg(opcode Opcode, _moved_payload ByteArray) Msg {
    return C.msg_new(C._opcode_t(opcode), _moved_payload)
}

func (self Msg) Free() { C.msg_free(self) }

func (self Msg) GetPayload() DataStream {
    return C.msg_get_payload(self)
}

func (self Msg) GetOpcode() Opcode {
    return Opcode(C.msg_get_opcode(self))
}

