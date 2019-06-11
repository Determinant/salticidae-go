package salticidae

// #include "salticidae/network.h"
import "C"

type MsgNetwork = *C.struct_msgnetwork_t

type MsgNetworkConn = *C.struct_msgnetwork_conn_t

type MsgNetworkConnMode = C.msgnetwork_conn_mode_t


func (self MsgNetworkConn) GetNet() MsgNetwork {
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
    return C.msgnetwork_new(ec, config)
}

func (self MsgNetwork) Free() { C.msgnetwork_free(self) }
func (self MsgNetwork) Listen(addr NetAddr) { C.msgnetwork_listen(self, addr) }
func (self MsgNetwork) Start() { C.msgnetwork_start(self) }

func (self MsgNetwork) SendMsg(msg Msg, conn MsgNetworkConn) { C.msgnetwork_send_msg(self, msg, conn) }
func (self MsgNetwork) Connect(addr NetAddr) { C.msgnetwork_connect(self, addr) }
func (self MsgNetwork) Terminate(conn MsgNetworkConn) { C.msgnetwork_terminate(self, conn) }

type MsgNetworkMsgCallback = C.msgnetwork_msg_callback_t
type MsgNetworkConnCallback = C.msgnetwork_conn_callback_t

func (self MsgNetwork) RegHandler(opcode Opcode, callback MsgNetworkMsgCallback, userdata rawptr_t) {
    C.msgnetwork_reg_handler(self, C._opcode_t(opcode), callback, userdata)
}

func (self MsgNetwork) RegConnHandler(callback MsgNetworkConnCallback, userdata rawptr_t) {
    C.msgnetwork_reg_conn_handler(self, callback, userdata)
}

type PeerNetwork = *C.struct_peernetwork_t

type PeerNetworkConn = *C.struct_peernetwork_conn_t

type PeerNetworkConfig = *C.struct_peernetwork_config_t

func NewPeerNetworkConfig() PeerNetworkConfig { return C.peernetwork_config_new() }

func NewPeerNetwork(ec EventContext, config PeerNetworkConfig) PeerNetwork {
    return C.peernetwork_new(ec, config)
}

func (self PeerNetwork) Free() { C.peernetwork_free(self) }

func (self PeerNetwork) AddPeer(paddr NetAddr) { C.peernetwork_add_peer(self, paddr) }

func (self PeerNetwork) HasPeer(paddr NetAddr) bool { return bool(C.peernetwork_has_peer(self, paddr)) }

func (self PeerNetwork) GetPeerConn(paddr NetAddr) PeerNetworkConn { return C.peernetwork_get_peer_conn(self, paddr) }

func (self PeerNetwork) AsMsgNetwork() MsgNetwork { return C.peernetwork_as_msgnetwork(self) }

func NewMsgNetworkConnFromPeerNetWorkConn(conn PeerNetworkConn) MsgNetworkConn { return C.msgnetwork_conn_new_from_peernetwork_conn(conn) }

func (self PeerNetwork) SendMsg(_moved_msg Msg, paddr NetAddr) { C.peernetwork_send_msg(self, _moved_msg, paddr) }

func (self PeerNetwork) MulticastMsg(_moved_msg Msg, paddrs []NetAddr) {
    base := uintptr(rawptr_t(&paddrs[0]))
    C.peernetwork_multicast_msg(self, _moved_msg, (*C.struct_netaddr_t)(rawptr_t(base)), C.size_t(len(paddrs)))
}

func (self PeerNetwork) Listen(listenAddr NetAddr) { C.peernetwork_listen(self, listenAddr) }
