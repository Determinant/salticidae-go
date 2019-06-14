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

func (self MsgNetworkConfig) BurstSize(size int) {
    C.msgnetwork_config_burst_size(self, C.size_t(size))
}

func (self MsgNetworkConfig) MaxListenBacklog(backlog int) {
    C.msgnetwork_config_max_listen_backlog(self, C.int(backlog))
}

func (self MsgNetworkConfig) ConnServerTimeout(timeout float64) {
    C.msgnetwork_config_conn_server_timeout(self, C.double(timeout))
}

func (self MsgNetworkConfig) SegBuffSize(size int) {
    C.msgnetwork_config_seg_buff_size(self, C.size_t(size))
}

func (self MsgNetworkConfig) NWorker(nworker int) {
    C.msgnetwork_config_nworker(self, C.size_t(nworker))
}

func (self MsgNetworkConfig) QueueCapacity(capacity int) {
    C.msgnetwork_config_queue_capacity(self, C.size_t(capacity))
}

func NewMsgNetwork(ec EventContext, config MsgNetworkConfig) MsgNetwork {
    return C.msgnetwork_new(ec, config)
}

func (self MsgNetwork) Free() { C.msgnetwork_free(self) }
func (self MsgNetwork) Listen(addr NetAddr, err *Error) { C.msgnetwork_listen(self, addr, err) }
func (self MsgNetwork) Start() { C.msgnetwork_start(self) }

func (self MsgNetwork) SendMsgByMove(msg Msg, conn MsgNetworkConn) { C.msgnetwork_send_msg_by_move(self, msg, conn) }
func (self MsgNetwork) Connect(addr NetAddr, err *Error) MsgNetworkConn { return C.msgnetwork_connect(self, addr, err) }
func (self MsgNetwork) Terminate(conn MsgNetworkConn) { C.msgnetwork_terminate(self, conn) }

func (self MsgNetworkConn) Copy() MsgNetworkConn { return C.msgnetwork_conn_copy(self) }
func (self MsgNetworkConn) Free() { C.msgnetwork_conn_free(self) }

type MsgNetworkMsgCallback = C.msgnetwork_msg_callback_t
type MsgNetworkConnCallback = C.msgnetwork_conn_callback_t
type MsgNetworkErrorCallback = C.msgnetwork_error_callback_t

func (self MsgNetwork) RegHandler(opcode Opcode, callback MsgNetworkMsgCallback, userdata rawptr_t) {
    C.msgnetwork_reg_handler(self, C._opcode_t(opcode), callback, userdata)
}

func (self MsgNetwork) RegConnHandler(callback MsgNetworkConnCallback, userdata rawptr_t) {
    C.msgnetwork_reg_conn_handler(self, callback, userdata)
}

func (self MsgNetwork) RegErrorHandler(callback MsgNetworkErrorCallback, userdata rawptr_t) {
    C.msgnetwork_reg_error_handler(self, callback, userdata)
}

type PeerNetwork = *C.struct_peernetwork_t

type PeerNetworkConn = *C.struct_peernetwork_conn_t

type PeerNetworkIdMode = C.peernetwork_id_mode_t

var (
    ID_MODE_IP_BASED = PeerNetworkIdMode(C.ID_MODE_IP_BASED)
    ID_MODE_IP_PORT_BASED = PeerNetworkIdMode(C.ID_MODE_IP_PORT_BASED)
)

type PeerNetworkConfig = *C.struct_peernetwork_config_t

func NewPeerNetworkConfig() PeerNetworkConfig { return C.peernetwork_config_new() }

func (self PeerNetworkConfig) Free() { C.peernetwork_config_free(self) }

func (self PeerNetworkConfig) RetryConnDelay(t_sec float64) {
    C.peernetwork_config_retry_conn_delay(self, C.double(t_sec))
}

func (self PeerNetworkConfig) PingPeriod(t_sec float64) {
    C.peernetwork_config_ping_period(self, C.double(t_sec))
}

func (self PeerNetworkConfig) ConnTimeout(t_sec float64) {
    C.peernetwork_config_conn_timeout(self, C.double(t_sec))
}

func (self PeerNetworkConfig) IdMode(mode PeerNetworkIdMode) {
    C.peernetwork_config_id_mode(self, mode)
}

func (self PeerNetworkConfig) AsMsgNetworkConfig() MsgNetworkConfig {
    return C.peernetwork_config_as_msgnetwork_config(self)
}

func NewPeerNetwork(ec EventContext, config PeerNetworkConfig) PeerNetwork {
    return C.peernetwork_new(ec, config)
}

func (self PeerNetwork) Free() { C.peernetwork_free(self) }

func (self PeerNetwork) AddPeer(paddr NetAddr) { C.peernetwork_add_peer(self, paddr) }

func (self PeerNetwork) HasPeer(paddr NetAddr) bool { return bool(C.peernetwork_has_peer(self, paddr)) }

func (self PeerNetwork) GetPeerConn(paddr NetAddr, err *Error) PeerNetworkConn { return C.peernetwork_get_peer_conn(self, paddr, err) }

func (self PeerNetwork) AsMsgNetwork() MsgNetwork { return C.peernetwork_as_msgnetwork(self) }

func NewMsgNetworkConnFromPeerNetWorkConn(conn PeerNetworkConn) MsgNetworkConn { return C.msgnetwork_conn_new_from_peernetwork_conn(conn) }

func (self PeerNetworkConn) Copy() PeerNetworkConn { return C.peernetwork_conn_copy(self) }

func (self PeerNetworkConn) Free() { C.peernetwork_conn_free(self) }

func (self PeerNetwork) SendMsgByMove(_moved_msg Msg, paddr NetAddr) { C.peernetwork_send_msg_by_move(self, _moved_msg, paddr) }

func (self PeerNetwork) MulticastMsgByMove(_moved_msg Msg, paddrs []NetAddr) {
    na := NewAddrArrayFromAddrs(paddrs)
    C.peernetwork_multicast_msg_by_move(self, _moved_msg, na)
}

func (self PeerNetwork) Listen(listenAddr NetAddr, err *Error) { C.peernetwork_listen(self, listenAddr, err) }
