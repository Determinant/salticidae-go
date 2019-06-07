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

func (self NetAddr) Free() { C.netaddr_free(self) }

type EventContext = *C.struct_eventcontext_t

func NewEventContext() EventContext { return C.eventcontext_new() }
func (self EventContext) Free() { C.eventcontext_free(self) }
func (self EventContext) Dispatch() { C.eventcontext_dispatch(self) }
func (self EventContext) Stop() { C.eventcontext_stop(self) }

type Opcode = uint8

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

type MsgNetworkInner = *C.struct_msgnetwork_t

type MsgNetwork struct {
    inner MsgNetworkInner
}

type MsgNetworkConn = *C.struct_msgnetwork_conn_t

type MsgNetworkConnMode = C.msgnetwork_conn_mode_t


func (self MsgNetworkConn) GetNet() MsgNetworkInner {
    return C.msgnetwork_conn_get_net(self)
}

var (
    CONN_MODE_ACTIVE = MsgNetworkConnMode(C.CONN_MODE_ACTIVE)
    CONN_MODE_PASSIVE = MsgNetworkConnMode(C.CONN_MODE_PASSIVE)
    CONN_MODE_DEAD = MsgNetworkConnMode(C.CONN_MODE_DEAD)
)

func (self MsgNetworkConn) GetMode() MsgNetworkConnMode {
    return C.msgnetwork_conn_get_mode(self)
}

func (self MsgNetworkConn) GetAddr() NetAddr {
    return C.msgnetwork_conn_get_addr(self)
}

type MsgNetworkConfig = *C.struct_msgnetwork_config_t

func NewMsgNetworkConfig() MsgNetworkConfig { return C.msgnetwork_config_new() }

func (self MsgNetworkConfig) Free() { C.msgnetwork_config_free(self) }

func NewMsgNetwork(ec EventContext, config MsgNetworkConfig) MsgNetwork {
    return MsgNetwork {
        inner: C.msgnetwork_new(ec, config),
    }
}

func (self MsgNetwork) Free() { C.msgnetwork_free(self.inner) }
func (self MsgNetwork) SendMsg(msg Msg, conn MsgNetworkConn) { self.inner.SendMsg(msg, conn) }
func (self MsgNetwork) Connect(addr NetAddr) { self.inner.Connect(addr) }
func (self MsgNetwork) Listen(addr NetAddr) { C.msgnetwork_listen(self.inner, addr) }
func (self MsgNetwork) Start() { C.msgnetwork_start(self.inner) }
func (self MsgNetwork) GetInner() MsgNetworkInner { return self.inner }

func (self MsgNetworkInner) SendMsg(msg Msg, conn MsgNetworkConn) {
    C.msgnetwork_send_msg(self, msg, conn)
}

func (self MsgNetworkInner) Connect(addr NetAddr) {
    C.msgnetwork_connect(self, addr)
}

type SigEvent = *C.sigev_t
type SigEventCallback = C.sigev_callback_t
var SIGTERM = C.SIGTERM
var SIGINT = C.SIGINT

func NewSigEvent(ec EventContext, cb SigEventCallback) SigEvent {
    return C.sigev_new(ec, cb)
}

func (self SigEvent) Add(sig int) { C.sigev_add(self, C.int(sig)) }
func (self SigEvent) Free() { C.sigev_free(self) }
