package salticidae

// #include <stdlib.h>
// #include "salticidae/netaddr.h"
import "C"

type NetAddr = *C.netaddr_t
type NetAddrArray = *C.netaddr_array_t

func NewAddrFromIPPortString(addr string) (res NetAddr) {
    c_str := C.CString(addr)
    res = C.netaddr_new_from_sipport(c_str)
    C.free(rawptr_t(c_str))
    return
}

func (self NetAddr) Free() { C.netaddr_free(self) }

func (self NetAddr) IsEq(other NetAddr) bool { return bool(C.netaddr_is_eq(self, other)) }

func (self NetAddr) IsNull() bool { return bool(C.netaddr_is_null(self)) }

func (self NetAddr) GetIP() uint32 { return uint32(C.netaddr_get_ip(self)) }

func (self NetAddr) GetPort() uint16 { return uint16(C.netaddr_get_port(self)) }

func NewAddrArrayFromAddrs(arr []NetAddr) NetAddrArray {
    size := len(arr)
    if size > 0 {
        base := (**C.netaddr_t)(&arr[0])
        return C.netaddr_array_new_from_addrs(base, C.size_t(size))
    } else {
        return C.netaddr_array_new()
    }
}
