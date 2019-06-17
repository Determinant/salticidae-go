package salticidae

// #include "salticidae/network.h"
import "C"
import "runtime"

type CMsgNetwork = *C.msgnetwork_t
type msgNetwork struct { inner CMsgNetwork }
type MsgNetwork = *msgNetwork

func MsgNetworkFromC(ptr CMsgNetwork) MsgNetwork {
    return &msgNetwork{ inner: ptr }
}

type CMsgNetworkConn = *C.msgnetwork_conn_t
type msgNetworkConn struct { inner CMsgNetworkConn }
type MsgNetworkConn = *msgNetworkConn

func MsgNetworkConnFromC(ptr CMsgNetworkConn) MsgNetworkConn {
    return &msgNetworkConn{ inner: ptr }
}

type MsgNetworkConnMode = C.msgnetwork_conn_mode_t

func (self MsgNetworkConn) free() { C.msgnetwork_conn_free(self.inner) }

func (self MsgNetworkConn) GetNet() MsgNetwork {
    return &msgNetwork{ inner: C.msgnetwork_conn_get_net(self.inner) }
}

var (
    CONN_MODE_ACTIVE = MsgNetworkConnMode(C.CONN_MODE_ACTIVE)
    CONN_MODE_PASSIVE = MsgNetworkConnMode(C.CONN_MODE_PASSIVE)
    CONN_MODE_DEAD = MsgNetworkConnMode(C.CONN_MODE_DEAD)
)

func (self MsgNetworkConn) GetMode() MsgNetworkConnMode {
    return C.msgnetwork_conn_get_mode(self.inner)
}

func (self MsgNetworkConn) GetAddr() NetAddr {
    res := &netAddr{ inner: C.msgnetwork_conn_get_addr(self.inner) }
    runtime.SetFinalizer(res, func(self NetAddr) { self.free() })
    return res
}

type CMsgNetworkConfig = *C.msgnetwork_config_t
type msgNetworkConfig struct { inner CMsgNetworkConfig }
type MsgNetworkConfig = *msgNetworkConfig

func MsgNetworkConfigFromC(ptr CMsgNetworkConfig) MsgNetworkConfig {
    return &msgNetworkConfig{ inner: ptr }
}

func NewMsgNetworkConfig() MsgNetworkConfig {
    res := &msgNetworkConfig{ inner: C.msgnetwork_config_new() }
    runtime.SetFinalizer(res, func(self MsgNetworkConfig) { self.free() })
    return res
}

func (self MsgNetworkConfig) free() { C.msgnetwork_config_free(self.inner) }

func (self MsgNetworkConfig) BurstSize(size int) {
    C.msgnetwork_config_burst_size(self.inner, C.size_t(size))
}

func (self MsgNetworkConfig) MaxListenBacklog(backlog int) {
    C.msgnetwork_config_max_listen_backlog(self.inner, C.int(backlog))
}

func (self MsgNetworkConfig) ConnServerTimeout(timeout float64) {
    C.msgnetwork_config_conn_server_timeout(self.inner, C.double(timeout))
}

func (self MsgNetworkConfig) SegBuffSize(size int) {
    C.msgnetwork_config_seg_buff_size(self.inner, C.size_t(size))
}

func (self MsgNetworkConfig) NWorker(nworker int) {
    C.msgnetwork_config_nworker(self.inner, C.size_t(nworker))
}

func (self MsgNetworkConfig) QueueCapacity(capacity int) {
    C.msgnetwork_config_queue_capacity(self.inner, C.size_t(capacity))
}

func NewMsgNetwork(ec EventContext, config MsgNetworkConfig) MsgNetwork {
    res := &msgNetwork { inner: C.msgnetwork_new(ec.inner, config.inner) }
    ec.attach(rawptr_t(res.inner), res)
    runtime.SetFinalizer(res, func(self MsgNetwork) { self.free() })
    return res
}

func (self MsgNetwork) free() { C.msgnetwork_free(self.inner) }
func (self MsgNetwork) Listen(addr NetAddr, err *Error) { C.msgnetwork_listen(self.inner, addr.inner, err) }
func (self MsgNetwork) Start() { C.msgnetwork_start(self.inner) }

func (self MsgNetwork) SendMsgByMove(msg Msg, conn MsgNetworkConn) { C.msgnetwork_send_msg_by_move(self.inner, msg.inner, conn.inner) }
func (self MsgNetwork) Connect(addr NetAddr, err *Error) MsgNetworkConn {
    res := &msgNetworkConn { inner: C.msgnetwork_connect(self.inner, addr.inner, err) }
    runtime.SetFinalizer(res, func(self MsgNetworkConn) { self.free() })
    return res
}
func (self MsgNetwork) Terminate(conn MsgNetworkConn) { C.msgnetwork_terminate(self.inner, conn.inner) }


func (self MsgNetwork) RegHandler(opcode Opcode, callback MsgNetworkMsgCallback, userdata rawptr_t) {
    C.msgnetwork_reg_handler(self.inner, C._opcode_t(opcode), callback, userdata)
}

func (self MsgNetwork) RegConnHandler(callback MsgNetworkConnCallback, userdata rawptr_t) {
    C.msgnetwork_reg_conn_handler(self.inner, callback, userdata)
}

func (self MsgNetwork) RegErrorHandler(callback MsgNetworkErrorCallback, userdata rawptr_t) {
    C.msgnetwork_reg_error_handler(self.inner, callback, userdata)
}


func (self MsgNetworkConn) Copy() MsgNetworkConn {
    res := &msgNetworkConn { inner: C.msgnetwork_conn_copy(self.inner) }
    runtime.SetFinalizer(res, func(self MsgNetworkConn) { self.free() })
    return res
}
func (self MsgNetworkConn) Free() { C.msgnetwork_conn_free(self.inner) }

type MsgNetworkMsgCallback = C.msgnetwork_msg_callback_t
type MsgNetworkConnCallback = C.msgnetwork_conn_callback_t
type MsgNetworkErrorCallback = C.msgnetwork_error_callback_t

type CPeerNetwork = *C.peernetwork_t
type peerNetwork struct { inner CPeerNetwork }
type PeerNetwork = *peerNetwork

func PeerNetworkFromC(ptr CPeerNetwork) PeerNetwork {
    return &peerNetwork{ inner: ptr }
}


type CPeerNetworkConn = *C.peernetwork_conn_t
type peerNetworkConn struct { inner CPeerNetworkConn }
type PeerNetworkConn = *peerNetworkConn

func PeerNetworkConnFromC(ptr CPeerNetworkConn) PeerNetworkConn {
    return &peerNetworkConn{ inner: ptr }
}

type PeerNetworkIdMode = C.peernetwork_id_mode_t

var (
    ID_MODE_IP_BASED = PeerNetworkIdMode(C.ID_MODE_IP_BASED)
    ID_MODE_IP_PORT_BASED = PeerNetworkIdMode(C.ID_MODE_IP_PORT_BASED)
)

type CPeerNetworkConfig = *C.peernetwork_config_t
type peerNetworkConfig struct { inner CPeerNetworkConfig }
type PeerNetworkConfig = *peerNetworkConfig

func PeerNetworkConfigFromC(ptr CPeerNetworkConfig) PeerNetworkConfig {
    return &peerNetworkConfig{ inner: ptr }
}

func NewPeerNetworkConfig() PeerNetworkConfig {
    res := &peerNetworkConfig { inner: C.peernetwork_config_new() }
    runtime.SetFinalizer(res, func(self PeerNetworkConfig) { self.free() })
    return res
}

func (self PeerNetworkConfig) free() { C.peernetwork_config_free(self.inner) }

func (self PeerNetworkConfig) RetryConnDelay(t_sec float64) {
    C.peernetwork_config_retry_conn_delay(self.inner, C.double(t_sec))
}

func (self PeerNetworkConfig) PingPeriod(t_sec float64) {
    C.peernetwork_config_ping_period(self.inner, C.double(t_sec))
}

func (self PeerNetworkConfig) ConnTimeout(t_sec float64) {
    C.peernetwork_config_conn_timeout(self.inner, C.double(t_sec))
}

func (self PeerNetworkConfig) IdMode(mode PeerNetworkIdMode) {
    C.peernetwork_config_id_mode(self.inner, mode)
}

func (self PeerNetworkConfig) AsMsgNetworkConfig() MsgNetworkConfig {
    return &msgNetworkConfig { inner: C.peernetwork_config_as_msgnetwork_config(self.inner) }
}

func NewPeerNetwork(ec EventContext, config PeerNetworkConfig) PeerNetwork {
    res := &peerNetwork { inner: C.peernetwork_new(ec.inner, config.inner) }
    ec.attach(rawptr_t(res.inner), res)
    runtime.SetFinalizer(res, func(self PeerNetwork) { self.free() })
    return res
}

func (self PeerNetwork) free() { C.peernetwork_free(self.inner) }

func (self PeerNetwork) AddPeer(paddr NetAddr) { C.peernetwork_add_peer(self.inner, paddr.inner) }

func (self PeerNetwork) HasPeer(paddr NetAddr) bool { return bool(C.peernetwork_has_peer(self.inner, paddr.inner)) }

func (self PeerNetwork) GetPeerConn(paddr NetAddr, err *Error) PeerNetworkConn {
    res := &peerNetworkConn{ inner: C.peernetwork_get_peer_conn(self.inner, paddr.inner, err) }
    runtime.SetFinalizer(res, func(self PeerNetworkConn) { self.free() })
    return res
}

func (self PeerNetwork) AsMsgNetwork() MsgNetwork { return &msgNetwork{ inner: C.peernetwork_as_msgnetwork(self.inner) } }

func NewMsgNetworkConnFromPeerNetWorkConn(conn PeerNetworkConn) MsgNetworkConn {
    res := &msgNetworkConn{ inner: C.msgnetwork_conn_new_from_peernetwork_conn(conn.inner) }
    runtime.SetFinalizer(res, func(self MsgNetworkConn) { self.free() })
    return res
}

func (self PeerNetworkConn) Copy() PeerNetworkConn {
    res := &peerNetworkConn { inner: C.peernetwork_conn_copy(self.inner) }
    runtime.SetFinalizer(res, func(self PeerNetworkConn) { self.free() })
    return res
}

func (self PeerNetworkConn) free() { C.peernetwork_conn_free(self.inner) }

func (self PeerNetwork) SendMsgByMove(_moved_msg Msg, paddr NetAddr) {
    C.peernetwork_send_msg_by_move(self.inner, _moved_msg.inner, paddr.inner)
}

func (self PeerNetwork) MulticastMsgByMove(_moved_msg Msg, paddrs []NetAddr) {
    na := NewAddrArrayFromAddrs(paddrs)
    C.peernetwork_multicast_msg_by_move(self.inner, _moved_msg.inner, na.inner)
}

func (self PeerNetwork) Listen(listenAddr NetAddr, err *Error) {
    C.peernetwork_listen(self.inner, listenAddr.inner, err)
}
