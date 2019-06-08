package salticidae

// #include <stdlib.h>
// #include "salticidae/netaddr.h"
import "C"

type NetAddr = *C.struct_netaddr_t

func NewAddrFromIPPortString(addr string) (res NetAddr) {
    c_str := C.CString(addr)
    res = C.netaddr_new_from_sipport(c_str)
    C.free(rawptr_t(c_str))
    return
}

func (self NetAddr) Free() { C.netaddr_free(self) }
