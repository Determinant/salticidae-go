package salticidae

// #include <stdlib.h>
// #include "salticidae/netaddr.h"
import "C"
import "runtime"

//// begin NetAddr def

// CNetAddr is the C pointer type for a NetAddr object
type CNetAddr = *C.netaddr_t
type netAddr struct {
	inner    CNetAddr
	autoFree bool
}

// NetAddr is a network address object.
type NetAddr = *netAddr

// NetAddrFromC converts an existing C pointer into a go object. Notice that
// when the go object does *not* own the resource of the C pointer, so it is
// only valid to the extent in which the given C pointer is valid. The C memory
// will not be deallocated when the go object is finalized by GC. This applies
// to all other "FromC" functions.
func NetAddrFromC(ptr CNetAddr) NetAddr {
	return &netAddr{inner: ptr}
}

func netAddrSetFinalizer(res NetAddr, autoFree bool) {
	res.autoFree = autoFree
	if res != nil && autoFree {
		runtime.SetFinalizer(res, func(self NetAddr) { self.Free() })
	}
}

// Free the underlying C pointer manually.
func (na NetAddr) Free() {
	C.netaddr_free(na.inner)
	if na.autoFree {
		runtime.SetFinalizer(na, nil)
	}
}

//// end NetAddr def

//// begin NetAddrArray def

// CNetAddrArray is the C pointer type for a NetAddrArray object.
type CNetAddrArray = *C.netaddr_array_t
type netAddrArray struct {
	inner    CNetAddrArray
	autoFree bool
}

// NetAddrArray is an array of network address.
type NetAddrArray = *netAddrArray

func NetAddrArrayFromC(ptr CNetAddrArray) NetAddrArray {
	return &netAddrArray{inner: ptr}
}

func netAddrArraySetFinalizer(res NetAddrArray, autoFree bool) {
	res.autoFree = autoFree
	if res != nil && autoFree {
		runtime.SetFinalizer(res, func(self NetAddrArray) { self.Free() })
	}
}

// Free the underlying C pointer manually.
func (naa NetAddrArray) Free() {
	C.netaddr_array_free(naa.inner)
	if naa.autoFree {
		runtime.SetFinalizer(naa, nil)
	}
}

//// end NetAddrArray def

//// begin NetAddr methods

// NewNetAddrFromIPPortString creates NetAddr from a TCP socket format string
// (e.g. 127.0.0.1:8888).
func NewNetAddrFromIPPortString(addr string, autoFree bool, err *Error) (res NetAddr) {
	cStr := C.CString(addr)
	res = NetAddrFromC(C.netaddr_new_from_sipport(cStr, err))
	C.free(rawPtr(cStr))
	netAddrSetFinalizer(res, autoFree)
	return
}

// IsEq checks if two addresses are the same.
func (na NetAddr) IsEq(other NetAddr) (res bool) {
	res = bool(C.netaddr_is_eq(na.inner, other.inner))
	runtime.KeepAlive(na)
	runtime.KeepAlive(other)
	return
}

// IsNull checks the address is empty.
func (na NetAddr) IsNull() (res bool) {
	res = bool(C.netaddr_is_null(na.inner))
	runtime.KeepAlive(na)
	return
}

// GetIP gets the 32-bit IP representation.
func (na NetAddr) GetIP() (res uint32) {
	res = uint32(C.netaddr_get_ip(na.inner))
	runtime.KeepAlive(na)
	return
}

// GetPort gets the 16-bit port number (in UNIX network byte order, so need to
// apply ntohs(), for example, to convert the returned integer to the local
// endianness).
func (na NetAddr) GetPort() (res uint16) {
	res = uint16(C.netaddr_get_port(na.inner))
	runtime.KeepAlive(na)
	return
}

// Copy the object. This is required if you want to keep the NetAddr returned
// (or passed as a callback parameter) by other salticidae methods (such like
// MsgNetwork/PeerNetwork), unless those method return a moved object.
func (na NetAddr) Copy(autoFree bool) (res NetAddr) {
	res = NetAddrFromC(C.netaddr_copy(na.inner))
	netAddrSetFinalizer(res, autoFree)
	runtime.KeepAlive(na)
	return
}

//// end NetAddr methods

//// begin NetAddrArray methods

// NewNetAddrArrayFromAddrs converts a Go slice of net addresses to NetAddrArray.
func NewNetAddrArrayFromAddrs(arr []NetAddr, autoFree bool) (res NetAddrArray) {
	size := len(arr)
	_arr := make([]CNetAddr, size)
	for i, v := range arr {
		_arr[i] = v.inner
		runtime.KeepAlive(v)
	}
	if size > 0 {
		// FIXME: here we assume struct of a single pointer has the same memory
		// footprint the pointer
		base := (**C.netaddr_t)(rawPtr(&_arr[0]))
		res = NetAddrArrayFromC(C.netaddr_array_new_from_addrs(base, C.size_t(size)))
	} else {
		res = NetAddrArrayFromC(C.netaddr_array_new())
	}
	runtime.KeepAlive(_arr)
	netAddrArraySetFinalizer(res, autoFree)
	return
}

//// end NetAddrArray methods
