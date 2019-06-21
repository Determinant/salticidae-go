package salticidae

// #include "salticidae/network.h"
import "C"
import "runtime"

// The C pointer type for a MsgNetwork handle.
type CMsgNetwork = *C.msgnetwork_t
type msgNetwork struct { inner CMsgNetwork }
// The handle for a message network.
type MsgNetwork = *msgNetwork

// Convert an existing C pointer into a go object. Notice that when the go
// object does *not* own the resource of the C pointer, so it is only valid to
// the extent in which the given C pointer is valid. The C memory will not be
// deallocated when the go object is finalized by GC. This applies to all other
// "FromC" functions.
func MsgNetworkFromC(ptr CMsgNetwork) MsgNetwork {
    return &msgNetwork{ inner: ptr }
}

// The C pointer type for a MsgNetworkConn handle.
type CMsgNetworkConn = *C.msgnetwork_conn_t
type msgNetworkConn struct { inner CMsgNetworkConn }
// The handle for a message network connection.
type MsgNetworkConn = *msgNetworkConn

func MsgNetworkConnFromC(ptr CMsgNetworkConn) MsgNetworkConn {
    return &msgNetworkConn{ inner: ptr }
}

var (
    CONN_MODE_ACTIVE = MsgNetworkConnMode(C.CONN_MODE_ACTIVE)
    CONN_MODE_PASSIVE = MsgNetworkConnMode(C.CONN_MODE_PASSIVE)
    CONN_MODE_DEAD = MsgNetworkConnMode(C.CONN_MODE_DEAD)
)

// The connection mode. CONN_MODE_ACTIVE: a connection established from the
// local side. CONN_MODE_PASSIVE: a connection established from the remote
// side. CONN_MODE_DEAD: a connection that is already closed.
type MsgNetworkConnMode = C.msgnetwork_conn_mode_t

func (self MsgNetworkConn) free() { C.msgnetwork_conn_free(self.inner) }

// Get the corresponding MsgNetwork handle that manages this connection. The
// returned handle is only valid during the lifetime of this connection.
func (self MsgNetworkConn) GetNet() MsgNetwork {
    return MsgNetworkFromC(C.msgnetwork_conn_get_net(self.inner))
}

func (self MsgNetworkConn) GetMode() MsgNetworkConnMode {
    return C.msgnetwork_conn_get_mode(self.inner)
}

// Get the address of the remote end of this connection. Use Copy() to make a
// copy of the address if you want to use the address object beyond the
// lifetime of the connection.
func (self MsgNetworkConn) GetAddr() NetAddr {
    return NetAddrFromC(C.msgnetwork_conn_get_addr(self.inner))
}

// Make a copy of the object. This is required if you want to keep the
// connection passed as a callback parameter by other salticidae methods (such
// like MsgNetwork/PeerNetwork).
func (self MsgNetworkConn) Copy() MsgNetworkConn {
    res := MsgNetworkConnFromC(C.msgnetwork_conn_copy(self.inner))
    runtime.SetFinalizer(res, func(self MsgNetworkConn) { self.free() })
    return res
}


// The C pointer type for a MsgNetworkConfig object.
type CMsgNetworkConfig = *C.msgnetwork_config_t
type msgNetworkConfig struct { inner CMsgNetworkConfig }
// Configuration for MsgNetwork.
type MsgNetworkConfig = *msgNetworkConfig

func MsgNetworkConfigFromC(ptr CMsgNetworkConfig) MsgNetworkConfig {
    return &msgNetworkConfig{ inner: ptr }
}

// Create the configuration object with default settings.
func NewMsgNetworkConfig() MsgNetworkConfig {
    res := MsgNetworkConfigFromC(C.msgnetwork_config_new())
    runtime.SetFinalizer(res, func(self MsgNetworkConfig) { self.free() })
    return res
}

func (self MsgNetworkConfig) free() { C.msgnetwork_config_free(self.inner) }

// Set the number of consecutive read attempts in the message delivery queue.
// Usually the default value is good enough. This is used to make the tradeoff
// between the event loop fairness and the amortization of syscall cost.
func (self MsgNetworkConfig) BurstSize(size int) {
    C.msgnetwork_config_burst_size(self.inner, C.size_t(size))
}

// Maximum backlogs (see POSIX TCP backlog).
func (self MsgNetworkConfig) MaxListenBacklog(backlog int) {
    C.msgnetwork_config_max_listen_backlog(self.inner, C.int(backlog))
}

// The timeout for connecting to the remote (in seconds).
func (self MsgNetworkConfig) ConnServerTimeout(timeout float64) {
    C.msgnetwork_config_conn_server_timeout(self.inner, C.double(timeout))
}

// The size for an inbound data chunk (per read() syscall).
func (self MsgNetworkConfig) SegBuffSize(size int) {
    C.msgnetwork_config_seg_buff_size(self.inner, C.size_t(size))
}

// The number of worker threads.
func (self MsgNetworkConfig) NWorker(nworker int) {
    C.msgnetwork_config_nworker(self.inner, C.size_t(nworker))
}

// The capacity of the send buffer.
func (self MsgNetworkConfig) QueueCapacity(capacity int) {
    C.msgnetwork_config_queue_capacity(self.inner, C.size_t(capacity))
}

// Create a message network handle which is attached to given event loop.
func NewMsgNetwork(ec EventContext, config MsgNetworkConfig) MsgNetwork {
    res := MsgNetworkFromC(C.msgnetwork_new(ec.inner, config.inner))
    ec.attach(rawptr_t(res.inner), res)
    runtime.SetFinalizer(res, func(self MsgNetwork) { self.free() })
    return res
}

func (self MsgNetwork) free() { C.msgnetwork_free(self.inner) }

// Start the message network (by spawning worker threads). This should be
// called before using any other methods.
func (self MsgNetwork) Start() { C.msgnetwork_start(self.inner) }

// Listen to the specified network address.
func (self MsgNetwork) Listen(addr NetAddr, err *Error) { C.msgnetwork_listen(self.inner, addr.inner, err) }

// Stop the message network. No other methods should be called after this.
func (self MsgNetwork) Stop() { C.msgnetwork_stop(self.inner) }

// Send a message through the given connection.
func (self MsgNetwork) SendMsg(msg Msg, conn MsgNetworkConn) {
    C.msgnetwork_send_msg(self.inner, msg.inner, conn.inner)
}

// Send a message through the given connection, using a worker thread to
// seralize and put data to the send buffer. The payload contained in the given
// msg will be moved and sent. Thus, no methods of msg involving the payload
// should be called afterwards.
func (self MsgNetwork) SendMsgDeferredByMove(msg Msg, conn MsgNetworkConn) {
    C.msgnetwork_send_msg_deferred_by_move(self.inner, msg.inner, conn.inner)
}

// Try to connect to a remote address. The connection handle is returned. The
// returned connection handle could be kept in your program.
func (self MsgNetwork) Connect(addr NetAddr, err *Error) MsgNetworkConn {
    res := MsgNetworkConnFromC(C.msgnetwork_connect(self.inner, addr.inner, err))
    runtime.SetFinalizer(res, func(self MsgNetworkConn) { self.free() })
    return res
}

// Terminate the given connection.
func (self MsgNetwork) Terminate(conn MsgNetworkConn) { C.msgnetwork_terminate(self.inner, conn.inner) }

// The C function pointer type which takes msg_t*, msgnetwork_conn_t* and void*
// (passing in the custom user data allocated by C.malloc) as parameters.
type MsgNetworkMsgCallback = C.msgnetwork_msg_callback_t
// The C function pointer type which takes msgnetwork_conn_t*, bool (true for
// the connection is established, false for the connection is terminated) and
// void* as parameters.
type MsgNetworkConnCallback = C.msgnetwork_conn_callback_t
// The C function Pointer type which takes SalticidaeCError* and void* as parameters.
type MsgNetworkErrorCallback = C.msgnetwork_error_callback_t

// Register a message handler for the type of message identified by opcode. The
// callback function will be invoked upon the delivery of each message with the
// given opcode, by the thread of the event loop the MsgNetwork is attached to.
func (self MsgNetwork) RegHandler(opcode Opcode, callback MsgNetworkMsgCallback, userdata rawptr_t) {
    C.msgnetwork_reg_handler(self.inner, C._opcode_t(opcode), callback, userdata)
}

// Register a connection handler invoked when the connection state is changed.
func (self MsgNetwork) RegConnHandler(callback MsgNetworkConnCallback, userdata rawptr_t) {
    C.msgnetwork_reg_conn_handler(self.inner, callback, userdata)
}

// Register an error handler invoked when there is recoverable errors during any
// asynchronous call/execution inside the MsgNetwork.
func (self MsgNetwork) RegErrorHandler(callback MsgNetworkErrorCallback, userdata rawptr_t) {
    C.msgnetwork_reg_error_handler(self.inner, callback, userdata)
}

// The C pointer type for a PeerNetwork handle.
type CPeerNetwork = *C.peernetwork_t
type peerNetwork struct { inner CPeerNetwork }
// The handle for a peer-to-peer network.
type PeerNetwork = *peerNetwork

func PeerNetworkFromC(ptr CPeerNetwork) PeerNetwork {
    return &peerNetwork{ inner: ptr }
}

// The C pointer type for a PeerNetworkConn handle.
type CPeerNetworkConn = *C.peernetwork_conn_t
type peerNetworkConn struct { inner CPeerNetworkConn }
// The handle for a PeerNetwork connection.
type PeerNetworkConn = *peerNetworkConn

func PeerNetworkConnFromC(ptr CPeerNetworkConn) PeerNetworkConn {
    return &peerNetworkConn{ inner: ptr }
}

var (
    ID_MODE_IP_BASED = PeerNetworkIdMode(C.ID_MODE_IP_BASED)
    ID_MODE_IP_PORT_BASED = PeerNetworkIdMode(C.ID_MODE_IP_PORT_BASED)
)

// The identity mode.  ID_MODE_IP_BASED: a remote peer is identified by the IP
// only. ID_MODE_IP_PORT_BASED: a remote peer is identified by IP + port
// number.
type PeerNetworkIdMode = C.peernetwork_id_mode_t

// The C pointer type for a PeerNetworkConfig handle.
type CPeerNetworkConfig = *C.peernetwork_config_t
type peerNetworkConfig struct { inner CPeerNetworkConfig }
// Configuration for PeerNetwork.
type PeerNetworkConfig = *peerNetworkConfig

func PeerNetworkConfigFromC(ptr CPeerNetworkConfig) PeerNetworkConfig {
    return &peerNetworkConfig{ inner: ptr }
}

// Create the configuration object with default settings.
func NewPeerNetworkConfig() PeerNetworkConfig {
    res := PeerNetworkConfigFromC(C.peernetwork_config_new())
    runtime.SetFinalizer(res, func(self PeerNetworkConfig) { self.free() })
    return res
}

func (self PeerNetworkConfig) free() { C.peernetwork_config_free(self.inner) }

// Set the connection retry delay (in seconds).
func (self PeerNetworkConfig) RetryConnDelay(t_sec float64) {
    C.peernetwork_config_retry_conn_delay(self.inner, C.double(t_sec))
}

// Set the period for sending ping messsages (in seconds).
func (self PeerNetworkConfig) PingPeriod(t_sec float64) {
    C.peernetwork_config_ping_period(self.inner, C.double(t_sec))
}

// Set the time it takes after sending a ping message before a connection is
// considered as broken.
func (self PeerNetworkConfig) ConnTimeout(t_sec float64) {
    C.peernetwork_config_conn_timeout(self.inner, C.double(t_sec))
}

// Set the identity mode.
func (self PeerNetworkConfig) IdMode(mode PeerNetworkIdMode) {
    C.peernetwork_config_id_mode(self.inner, mode)
}

// Use the PeerNetworkConfig object as a MsgNetworkConfig object (to invoke the
// methods inherited from MsgNetworkConfig, such as NWorker).
func (self PeerNetworkConfig) AsMsgNetworkConfig() MsgNetworkConfig {
    return MsgNetworkConfigFromC(C.peernetwork_config_as_msgnetwork_config(self.inner))
}

// Create a peer-to-peer message network handle.
func NewPeerNetwork(ec EventContext, config PeerNetworkConfig) PeerNetwork {
    res := PeerNetworkFromC(C.peernetwork_new(ec.inner, config.inner))
    ec.attach(rawptr_t(res.inner), res)
    runtime.SetFinalizer(res, func(self PeerNetwork) { self.free() })
    return res
}

func (self PeerNetwork) free() { C.peernetwork_free(self.inner) }

// Add a peer to the list of known peers. The P2P network will try to keep
// bi-direction connections to all known peers in the list (through
// reconnection and connection deduplication).
func (self PeerNetwork) AddPeer(paddr NetAddr) { C.peernetwork_add_peer(self.inner, paddr.inner) }

// Test whether a peer is already in the list.
func (self PeerNetwork) HasPeer(paddr NetAddr) bool { return bool(C.peernetwork_has_peer(self.inner, paddr.inner)) }

// Get the connection of the known peer. The connection handle is returned. The
// returned connection handle could be kept in your program.
func (self PeerNetwork) GetPeerConn(paddr NetAddr, err *Error) PeerNetworkConn {
    res := PeerNetworkConnFromC(C.peernetwork_get_peer_conn(self.inner, paddr.inner, err))
    runtime.SetFinalizer(res, func(self PeerNetworkConn) { self.free() })
    return res
}

// Use the PeerNetwork handle as a MsgNetwork handle (to invoke the methods
// inherited from MsgNetwork, such as RegHandler).
func (self PeerNetwork) AsMsgNetwork() MsgNetwork {
    return MsgNetworkFromC(C.peernetwork_as_msgnetwork(self.inner))
}

// Use the MsgNetwork handle as a PeerNetwork handle (forcing the conversion).
func (self MsgNetwork) AsPeerNetworkUnsafe() PeerNetwork {
    return PeerNetworkFromC(C.msgnetwork_as_peernetwork_unsafe(self.inner))
}

// Create a MsgNetworkConn handle from a PeerNetworkConn (representing the same
// connection).
func NewMsgNetworkConnFromPeerNetWorkConn(conn PeerNetworkConn) MsgNetworkConn {
    res := MsgNetworkConnFromC(C.msgnetwork_conn_new_from_peernetwork_conn(conn.inner))
    runtime.SetFinalizer(res, func(self MsgNetworkConn) { self.free() })
    return res
}

// Create a PeerNetworkConn handle from a MsgNetworkConn (representing the same
// connection and forcing the conversion).
func NewPeerNetworkConnFromMsgNetWorkConnUnsafe(conn MsgNetworkConn) PeerNetworkConn {
    res := PeerNetworkConnFromC(C.peernetwork_conn_new_from_msgnetwork_conn_unsafe(conn.inner))
    runtime.SetFinalizer(res, func(self PeerNetworkConn) { self.free() })
    return res
}

// Make a copy of the connection handle.
func (self PeerNetworkConn) Copy() PeerNetworkConn {
    res := PeerNetworkConnFromC(C.peernetwork_conn_copy(self.inner))
    runtime.SetFinalizer(res, func(self PeerNetworkConn) { self.free() })
    return res
}

func (self PeerNetworkConn) free() { C.peernetwork_conn_free(self.inner) }

// Listen to the specified network address. Notice that this method overrides
// Listen() in MsgNetwork, so you should always call this one instead of
// AsMsgNetwork().Listen().
func (self PeerNetwork) Listen(listenAddr NetAddr, err *Error) {
    C.peernetwork_listen(self.inner, listenAddr.inner, err)
}

// Send a message to the given peer.
func (self PeerNetwork) SendMsg(msg Msg, addr NetAddr) {
    C.peernetwork_send_msg(self.inner, msg.inner, addr.inner)
}

// Send a message to the given peer, using a worker thread to seralize and put
// data to the send buffer. The payload contained in the given msg will be
// moved and sent. Thus, no methods of msg involving the payload should be
// called afterwards.
func (self PeerNetwork) SendMsgDeferredByMove(msg Msg, addr NetAddr) {
    C.peernetwork_send_msg_deferred_by_move(self.inner, msg.inner, addr.inner)
}

// Send a message to the given list of peers. The payload contained in the
// given msg will be moved and sent. Thus, no methods of msg involving the
// payload should be called afterwards.
func (self PeerNetwork) MulticastMsgByMove(msg Msg, paddrs []NetAddr) {
    na := NewAddrArrayFromAddrs(paddrs)
    C.peernetwork_multicast_msg_by_move(self.inner, msg.inner, na.inner)
}

// The C function pointer type which takes netaddr_t* and void* (passing in the
// custom user data allocated by C.malloc) as parameters.
type MsgNetworkUnknownPeerCallback = C.msgnetwork_unknown_peer_callback_t

// Register a connection handler invoked when a remote peer that is not in the
// list of known peers attempted to connect. By default configuration, the
// connection was rejected, and you can call AddPeer() to enroll this peer if
// you hope to establish the connection soon.
func (self PeerNetwork) RegUnknownPeerHandler(callback MsgNetworkUnknownPeerCallback, userdata rawptr_t) {
    C.peernetwork_reg_unknown_peer_handler(self.inner, callback, userdata)
}
