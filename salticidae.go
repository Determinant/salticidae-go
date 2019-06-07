package salticidae


// #cgo CFLAGS: -I${SRCDIR}/salticidae/include/
// #cgo LDFLAGS: ${SRCDIR}/salticidae/libsalticidae.so -Wl,-rpath=${SRCDIR}/salticidae/
// #include <stdlib.h>
// #include <signal.h>
// #include "salticidae/netaddr.h"
// #include "salticidae/network.h"
// #include "salticidae/event.h"
import "C"
import "unsafe"

type rawptr_t = unsafe.Pointer

type NetAddr = *C.struct_netaddr_t

func NewAddrFromIPPortString(addr string) (res NetAddr) {
    c_str := C.CString(addr)
    res = C.netaddr_new_from_sipport(c_str)
    C.free(rawptr_t(c_str))
    return
}

type MsgNetwork = *C.struct_msgnetwork_t

type EventContext = *C.struct_eventcontext_t

func NewEventContext() EventContext { return C.eventcontext_new() }
func (self EventContext) Dispatch() { C.eventcontext_dispatch(self) }
func (self EventContext) Stop() { C.eventcontext_stop(self) }

type MsgNetworkConfig = *C.struct_msgnetwork_config_t

func NewMsgNetworkConfig() MsgNetworkConfig { return C.msgnetwork_config_new() }

func NewMsgNetwork(ec EventContext, config MsgNetworkConfig) MsgNetwork {
    return C.msgnetwork_new(ec, config)
}

func (self MsgNetwork) Start() { C.msgnetwork_start(self) }
func (self MsgNetwork) Listen(addr NetAddr) { C.msgnetwork_listen(self, addr) }
func (self MsgNetwork) Connect(addr NetAddr) { C.msgnetwork_connect(self, addr) }
func (self MsgNetwork) Free() { C.msgnetwork_free(self) }

type SigEvent = *C.sigev_t
type SigEventCallback = C.sigev_callback_t
var SIGTERM = C.SIGTERM
var SIGINT = C.SIGINT

func NewSigEvent(ec EventContext, cb SigEventCallback) SigEvent {
    return C.sigev_new(ec, cb)
}

func (self SigEvent) Add(sig int) { C.sigev_add(self, C.int(sig)) }
func (self SigEvent) Free() { C.sigev_free(self) }
