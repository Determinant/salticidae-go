package salticidae

// #include <stdlib.h>
// #include "salticidae/netaddr.h"
import "C"
import "runtime"

type netAddr struct {
    inner *C.netaddr_t
}

type NetAddr = *netAddr

type netAddrArray struct {
    inner *C.netaddr_array_t
}

type NetAddrArray = *netAddrArray

func NewAddrFromIPPortString(addr string) (res NetAddr) {
    c_str := C.CString(addr)
    res = &netAddr{ inner: C.netaddr_new_from_sipport(c_str) }
    C.free(rawptr_t(c_str))
    runtime.SetFinalizer(res, func(self NetAddr) { self.free() })
    return
}

func NewAddrArrayFromAddrs(arr []NetAddr) (res NetAddrArray) {
    size := len(arr)
    if size > 0 {
        // FIXME: here we assume struct of a single pointer has the same memory
        // footprint the pointer
        base := (**C.netaddr_t)(rawptr_t(&arr[0]))
        res = &netAddrArray{ inner: C.netaddr_array_new_from_addrs(base, C.size_t(size)) }
    } else {
        res = &netAddrArray{ inner: C.netaddr_array_new() }
    }
    runtime.SetFinalizer(res, func(self NetAddrArray) { self.free() })
    return
}

func (self NetAddr) free() { C.netaddr_free(self.inner) }

func (self NetAddr) IsEq(other NetAddr) bool { return bool(C.netaddr_is_eq(self.inner, other.inner)) }

func (self NetAddr) IsNull() bool { return bool(C.netaddr_is_null(self.inner)) }

func (self NetAddr) GetIP() uint32 { return uint32(C.netaddr_get_ip(self.inner)) }

func (self NetAddr) GetPort() uint16 { return uint16(C.netaddr_get_port(self.inner)) }

func (self NetAddrArray) free() { C.netaddr_array_free(self.inner) }
