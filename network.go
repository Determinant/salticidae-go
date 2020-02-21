package salticidae

// #include <stdlib.h>
// #include "salticidae/network.h"
import "C"
import "runtime"

// The C pointer type for a MsgNetwork handle.
type CMsgNetwork = *C.msgnetwork_t
type msgNetwork struct{ inner CMsgNetwork }

// The handle for a message network.
type MsgNetwork = *msgNetwork

// Convert an existing C pointer into a go object. Notice that when the go
// object does *not* own the resource of the C pointer, so it is only valid to
// the extent in which the given C pointer is valid. The C memory will not be
// deallocated when the go object is finalized by GC. This applies to all other
// "FromC" functions.
func MsgNetworkFromC(ptr CMsgNetwork) MsgNetwork {
	return &msgNetwork{inner: ptr}
}

// The C pointer type for a MsgNetworkConn handle.
type CMsgNetworkConn = *C.msgnetwork_conn_t
type msgNetworkConn struct {
	inner    CMsgNetworkConn
	autoFree bool
}

// The handle for a message network connection.
type MsgNetworkConn = *msgNetworkConn

func MsgNetworkConnFromC(ptr CMsgNetworkConn) MsgNetworkConn {
	return &msgNetworkConn{inner: ptr}
}

var (
	CONN_MODE_ACTIVE  = MsgNetworkConnMode(C.CONN_MODE_ACTIVE)
	CONN_MODE_PASSIVE = MsgNetworkConnMode(C.CONN_MODE_PASSIVE)
)

// The connection mode. CONN_MODE_ACTIVE: a connection established from the
// local side. CONN_MODE_PASSIVE: a connection established from the remote
// side. CONN_MODE_DEAD: a connection that is already closed.
type MsgNetworkConnMode = C.msgnetwork_conn_mode_t

func msgNetworkConnSetFinalizer(res MsgNetworkConn, autoFree bool) {
	res.autoFree = autoFree
	if res != nil && autoFree {
		runtime.SetFinalizer(res, func(self MsgNetworkConn) { self.Free() })
	}
}

func (self MsgNetworkConn) Free() {
	C.msgnetwork_conn_free(self.inner)
	if self.autoFree {
		runtime.SetFinalizer(self, nil)
	}
}

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

// Check if the connection has been terminated.
func (self MsgNetworkConn) IsTerminated() bool {
	return bool(C.msgnetwork_conn_is_terminated(self.inner))
}

// Get the certificate of the remote end of this connection. Use Copy() to make a
// copy of the certificate if you want to use the certificate object beyond the
// lifetime of the connection.
func (self MsgNetworkConn) GetPeerCert() X509 {
	return &x509{inner: C.msgnetwork_conn_get_peer_cert(self.inner)}
}

// Make a copy of the object. This is required if you want to keep the
// connection passed as a callback parameter by other salticidae methods (such
// like MsgNetwork/PeerNetwork).
func (self MsgNetworkConn) Copy(autoFree bool) MsgNetworkConn {
	res := MsgNetworkConnFromC(C.msgnetwork_conn_copy(self.inner))
	msgNetworkConnSetFinalizer(res, autoFree)
	return res
}

// The C pointer type for a MsgNetworkConfig object.
type CMsgNetworkConfig = *C.msgnetwork_config_t
type msgNetworkConfig struct{ inner CMsgNetworkConfig }

// Configuration for MsgNetwork.
type MsgNetworkConfig = *msgNetworkConfig

func MsgNetworkConfigFromC(ptr CMsgNetworkConfig) MsgNetworkConfig {
	return &msgNetworkConfig{inner: ptr}
}

// Create the configuration object with default settings.
func NewMsgNetworkConfig() MsgNetworkConfig {
	res := MsgNetworkConfigFromC(C.msgnetwork_config_new())
	runtime.SetFinalizer(res, func(self MsgNetworkConfig) { self.free() })
	return res
}

func (self MsgNetworkConfig) free() { C.msgnetwork_config_free(self.inner) }

// Set the maximum message length (NOTE: by default it is 1 KBytes)
func (self MsgNetworkConfig) MaxMsgSize(size int) {
	C.msgnetwork_config_max_msg_size(self.inner, C.size_t(size))
}

// Set the maximum message queue size (the queue for buffering received
// messages to be processed by message handlers).
func (self MsgNetworkConfig) MaxMsgQueueSize(size int) {
	C.msgnetwork_config_max_msg_queue_size(self.inner, C.size_t(size))
}

// Set the number of consecutive read attempts in the message delivery queue.
// Usually the default value is good enough. This is used to make the tradeoff
// between the event loop fairness and the amortization of syscall cost.
func (self MsgNetworkConfig) BurstSize(size int) {
	C.msgnetwork_config_burst_size(self.inner, C.size_t(size))
}

// Set the maximum backlogs (see POSIX TCP backlog).
func (self MsgNetworkConfig) MaxListenBacklog(backlog int) {
	C.msgnetwork_config_max_listen_backlog(self.inner, C.int(backlog))
}

// Set the timeout for connecting to the remote (in seconds).
func (self MsgNetworkConfig) ConnServerTimeout(timeout float64) {
	C.msgnetwork_config_conn_server_timeout(self.inner, C.double(timeout))
}

// Set the size for an inbound data chunk (per read() syscall).
func (self MsgNetworkConfig) RecvChunkSize(size int) {
	C.msgnetwork_config_recv_chunk_size(self.inner, C.size_t(size))
}

// Set the number of worker threads.
func (self MsgNetworkConfig) NWorker(nworker int) {
	C.msgnetwork_config_nworker(self.inner, C.size_t(nworker))
}

// Set the maximum send buffer size.
func (self MsgNetworkConfig) MaxSendBuffSize(size int) {
	C.msgnetwork_config_max_send_buff_size(self.inner, C.size_t(size))
}

// Set the maximum recv buffer size.
func (self MsgNetworkConfig) MaxRecvBuffSize(size int) {
	C.msgnetwork_config_max_recv_buff_size(self.inner, C.size_t(size))
}

// Specify whether to use SSL/TLS. When this flag is enabled, MsgNetwork uses
// TLSKey (or TLSKeyFile) or TLSCert (or TLSCertFile) to setup the underlying
// OpenSSL.
func (self MsgNetworkConfig) EnableTLS(enabled bool) {
	C.msgnetwork_config_enable_tls(self.inner, C.bool(enabled))
}

// Load the TLS key from a file. The file should be an unencrypted PEM file.
// Use TLSKey() instead for more complex usage.
func (self MsgNetworkConfig) TLSKeyFile(fname string) {
	c_str := C.CString(fname)
	C.msgnetwork_config_tls_key_file(self.inner, c_str)
	C.free(rawptr_t(c_str))
}

// Load the TLS certificate from a file. The file should be an unencrypted
// (X509) PEM file.  Use TLSCert() instead for more complex usage.
func (self MsgNetworkConfig) TLSCertFile(fname string) {
	c_str := C.CString(fname)
	C.msgnetwork_config_tls_cert_file(self.inner, c_str)
	C.free(rawptr_t(c_str))
}

// Use the given TLS key. This overrides TLSKeyFile(). pkey will be moved and
// kept by the library. Thus, no methods of msg involving the payload should be
// called afterwards.
func (self MsgNetworkConfig) TLSKeyByMove(pkey PKey) {
	C.msgnetwork_config_tls_key_by_move(self.inner, pkey.inner)
}

//// Load the TLS certificate from a file. The file should be an unencrypted
//// (X509) PEM file.  Use TLSCert() instead for more complex usage.
//func (self MsgNetworkConfig) TLSCert(fname string) {
//    c_str := C.CString(fname)
//    C.msgnetwork_config_tls_cert(self.inner, c_str)
//    C.free(rawptr_t(c_str))
//}

// Create a message network handle which is attached to given event loop.
func NewMsgNetwork(ec EventContext, config MsgNetworkConfig, err *Error) MsgNetwork {
	res := MsgNetworkFromC(C.msgnetwork_new(ec.inner, config.inner, err))
	if res != nil {
		ec.attach(rawptr_t(res.inner), res)
		runtime.SetFinalizer(res, func(self MsgNetwork) { self.free() })
	}
	return res
}

func (self MsgNetwork) free() { C.msgnetwork_free(self.inner) }

// Start the message network (by spawning worker threads). This should be
// called before using any other methods.
func (self MsgNetwork) Start() { C.msgnetwork_start(self.inner) }

// Listen to the specified network address.
func (self MsgNetwork) Listen(addr NetAddr, err *Error) {
	C.msgnetwork_listen(self.inner, addr.inner, err)
}

// Stop the message network. No other methods should be called after this.
func (self MsgNetwork) Stop() { C.msgnetwork_stop(self.inner) }

// Send a message through the given connection.
func (self MsgNetwork) SendMsg(msg Msg, conn MsgNetworkConn) bool {
	return bool(C.msgnetwork_send_msg(self.inner, msg.inner, conn.inner))
}

// Send a message through the given connection, using a worker thread to
// seralize and put data to the send buffer. The payload contained in the given
// msg will be moved and sent. Thus, no methods of msg involving the payload
// should be called afterwards.
func (self MsgNetwork) SendMsgDeferredByMove(msg Msg, conn MsgNetworkConn) int32 {
	return int32(C.msgnetwork_send_msg_deferred_by_move(self.inner, msg.inner, conn.inner))
}

// Try to connect to a remote address. The connection handle is returned. The
// returned connection handle could be kept in your program.
func (self MsgNetwork) ConnectSync(addr NetAddr, autoFree bool, err *Error) MsgNetworkConn {
	res := MsgNetworkConnFromC(C.msgnetwork_connect_sync(self.inner, addr.inner, err))
	msgNetworkConnSetFinalizer(res, autoFree)
	return res
}

// Try to connect to a remote address (async). It returns an id which could be
// used to identify the corresponding error in the error callback.
func (self MsgNetwork) Connect(addr NetAddr) int32 {
	return int32(C.msgnetwork_connect(self.inner, addr.inner))
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
type peerNetwork struct{ inner CPeerNetwork }

// The handle for a peer-to-peer network.
type PeerNetwork = *peerNetwork

func PeerNetworkFromC(ptr CPeerNetwork) PeerNetwork {
	return &peerNetwork{inner: ptr}
}

// The C pointer type for a PeerNetworkConn handle.
type CPeerNetworkConn = *C.peernetwork_conn_t
type peerNetworkConn struct {
	inner    CPeerNetworkConn
	autoFree bool
}

// The handle for a PeerNetwork connection.
type PeerNetworkConn = *peerNetworkConn

func PeerNetworkConnFromC(ptr CPeerNetworkConn) PeerNetworkConn {
	return &peerNetworkConn{inner: ptr}
}

func peerNetworkConnSetFinalizer(res PeerNetworkConn, autoFree bool) {
	res.autoFree = autoFree
	if res != nil && autoFree {
		runtime.SetFinalizer(res, func(self PeerNetworkConn) { self.Free() })
	}
}

var (
	ADDR_BASED = PeerNetworkIdMode(C.ID_MODE_ADDR_BASED)
	CERT_BASED = PeerNetworkIdMode(C.ID_MODE_CERT_BASED)
)

// The identity mode.  ID_MODE_IP_BASED: a remote peer is identified by the IP
// only. ID_MODE_IP_PORT_BASED: a remote peer is identified by IP + port
// number.
type PeerNetworkIdMode = C.peernetwork_id_mode_t

// The C pointer type for a PeerNetworkConfig handle.
type CPeerNetworkConfig = *C.peernetwork_config_t
type peerNetworkConfig struct{ inner CPeerNetworkConfig }

// Configuration for PeerNetwork.
type PeerNetworkConfig = *peerNetworkConfig

func PeerNetworkConfigFromC(ptr CPeerNetworkConfig) PeerNetworkConfig {
	return &peerNetworkConfig{inner: ptr}
}

// Create the configuration object with default settings.
func NewPeerNetworkConfig() PeerNetworkConfig {
	res := PeerNetworkConfigFromC(C.peernetwork_config_new())
	runtime.SetFinalizer(res, func(self PeerNetworkConfig) { self.free() })
	return res
}

func (self PeerNetworkConfig) free() { C.peernetwork_config_free(self.inner) }

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

// The C pointer type for a PeerId object
type CPeerId = *C.peerid_t
type peerId struct {
	inner    CPeerId
	autoFree bool
}

// Peer identity object
type PeerId = *peerId

// Convert an existing C pointer into a go object.
func PeerIdFromC(ptr CPeerId) PeerId {
	return &peerId{inner: ptr}
}

func peerIdSetFinalizer(res PeerId, autoFree bool) {
	res.autoFree = autoFree
	if res != nil && autoFree {
		runtime.SetFinalizer(res, func(self PeerId) { self.Free() })
	}
}

// Create PeerId from a NetAddr
func NewPeerIdFromNetAddr(addr NetAddr, autoFree bool) (res PeerId) {
	res = &peerId{inner: C.peerid_new_from_netaddr(addr.inner)}
	peerIdSetFinalizer(res, autoFree)
	return
}

// Create PeerId from a X509 certificate
func NewPeerIdFromX509(cert X509, autoFree bool) (res PeerId) {
	res = &peerId{inner: C.peerid_new_from_x509(cert.inner)}
	peerIdSetFinalizer(res, autoFree)
	return
}

func (self PeerId) Free() {
	C.peerid_free(self.inner)
	if self.autoFree {
		runtime.SetFinalizer(self, nil)
	}
}

// The C pointer type for a PeerIdArray object
type CPeerIdArray = *C.peerid_array_t
type peerIdArray struct {
	inner    CPeerIdArray
	autoFree bool
}

// An array of peer ids.
type PeerIdArray = *peerIdArray

func PeerIdArrayFromC(ptr CPeerIdArray) PeerIdArray {
	return &peerIdArray{inner: ptr}
}

func peerIdArraySetFinalizer(res PeerIdArray, autoFree bool) {
	res.autoFree = autoFree
	if res != nil && autoFree {
		runtime.SetFinalizer(res, func(self PeerIdArray) { self.Free() })
	}
}

// Convert a Go slice of peer ids to PeerIdArray.
func NewPeerIdArrayFromPeers(arr []PeerId, autoFree bool) (res PeerIdArray) {
	size := len(arr)
	_arr := make([]CPeerId, size)
	for i, v := range arr {
		_arr[i] = v.inner
	}
	if size > 0 {
		// FIXME: here we assume struct of a single pointer has the same memory
		// footprint the pointer
		base := (**C.peerid_t)(rawptr_t(&_arr[0]))
		res = PeerIdArrayFromC(C.peerid_array_new_from_peers(base, C.size_t(size)))
	} else {
		res = PeerIdArrayFromC(C.peerid_array_new())
	}
	runtime.KeepAlive(_arr)
	peerIdArraySetFinalizer(res, autoFree)
	return
}

func (self PeerIdArray) Free() {
	C.peerid_array_free(self.inner)
	if self.autoFree {
		runtime.SetFinalizer(self, nil)
	}
}

// Create a peer-to-peer message network handle.
func NewPeerNetwork(ec EventContext, config PeerNetworkConfig, err *Error) PeerNetwork {
	res := PeerNetworkFromC(C.peernetwork_new(ec.inner, config.inner, err))
	if res != nil {
		ec.attach(rawptr_t(res.inner), res)
		runtime.SetFinalizer(res, func(self PeerNetwork) { self.free() })
	}
	return res
}

func (self PeerNetwork) free() { C.peernetwork_free(self.inner) }

// Register a peer to the list of known peers. The P2P network will try to keep
// bi-direction connections to all known peers in the list (through
// reconnection and connection deduplication). This is an async call and the
// call id is returned as the reference for error handling.
func (self PeerNetwork) AddPeer(peer PeerId) int32 {
	return int32(C.peernetwork_add_peer(self.inner, peer.inner))
}

// Remove a peer from the list of known peers. The P2P network will teardown
// the existing bi-direction connection and the PeerHandler callback will not
// be called. This is an async call.
func (self PeerNetwork) DelPeer(peer PeerId) int32 {
	return int32(C.peernetwork_del_peer(self.inner, peer.inner))
}

// Test whether a peer is already in the list.
func (self PeerNetwork) HasPeer(peer PeerId) bool {
	return bool(C.peernetwork_has_peer(self.inner, peer.inner))
}

// Get the connection of the known peer. The connection handle is returned. The
// returned connection handle could be kept in your program.
func (self PeerNetwork) GetPeerConn(peer PeerId, autoFree bool, err *Error) PeerNetworkConn {
	res := PeerNetworkConnFromC(C.peernetwork_get_peer_conn(self.inner, peer.inner, err))
	peerNetworkConnSetFinalizer(res, autoFree)
	return res
}

// Set the IP address of the registered peer, used to connect to the peer. The
// address for a peer is by default empty and a p2p connection can only be
// established from the other side in this case (which is common for the peers
// behind some firewall/router). This is an async call.
func (self PeerNetwork) SetPeerAddr(peer PeerId, addr NetAddr) int32 {
	return int32(C.peernetwork_set_peer_addr(self.inner, peer.inner, addr.inner))
}

// Try to connect to the peer. If ntry > 0, it specifies the maximum number of
// attempts before giving up. If ntry = 0, it stops any ongoing/established
// connection and future attempts.  If ntry = -1, reconnection is attempted
// indefinitely. retryDelay specifies the minimum delay (in seconds) between
// two attempts. When ntry != 0, once peer is connected, the retry state is
// reset with ntry.
func (self PeerNetwork) ConnPeer(peer PeerId, ntry int32, retryDelay float64) int32 {
	return int32(C.peernetwork_conn_peer(self.inner, peer.inner, C.int(ntry), C.double(retryDelay)))
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
func NewMsgNetworkConnFromPeerNetworkConn(conn PeerNetworkConn, autoFree bool) MsgNetworkConn {
	res := MsgNetworkConnFromC(C.msgnetwork_conn_new_from_peernetwork_conn(conn.inner))
	msgNetworkConnSetFinalizer(res, autoFree)
	return res
}

// Create a PeerNetworkConn handle from a MsgNetworkConn (representing the same
// connection and forcing the conversion).
func NewPeerNetworkConnFromMsgNetworkConnUnsafe(conn MsgNetworkConn, autoFree bool) PeerNetworkConn {
	res := PeerNetworkConnFromC(C.peernetwork_conn_new_from_msgnetwork_conn_unsafe(conn.inner))
	peerNetworkConnSetFinalizer(res, autoFree)
	return res
}

// Make a copy of the connection handle.
func (self PeerNetworkConn) Copy(autoFree bool) PeerNetworkConn {
	res := PeerNetworkConnFromC(C.peernetwork_conn_copy(self.inner))
	peerNetworkConnSetFinalizer(res, autoFree)
	return res
}

// Get the listening address of the remote peer (no Copy() is needed).
func (self PeerNetworkConn) GetPeerAddr(autoFree bool) NetAddr {
	res := NetAddrFromC(C.peernetwork_conn_get_peer_addr(self.inner))
	netAddrSetFinalizer(res, autoFree)
	return res
}

func (self PeerNetworkConn) Free() {
	C.peernetwork_conn_free(self.inner)
	if self.autoFree {
		runtime.SetFinalizer(self, nil)
	}
}

// Listen to the specified network address. Notice that this method overrides
// Listen() in MsgNetwork, so you should always call this one instead of
// AsMsgNetwork().Listen().
func (self PeerNetwork) Listen(listenAddr NetAddr, err *Error) {
	C.peernetwork_listen(self.inner, listenAddr.inner, err)
}

// Send a message to the given peer.
func (self PeerNetwork) SendMsg(msg Msg, peer PeerId) bool {
	return bool(C.peernetwork_send_msg(self.inner, msg.inner, peer.inner))
}

// Send a message to the given peer, using a worker thread to seralize and put
// data to the send buffer. The payload contained in the given msg will be
// moved and sent. Thus, no methods of msg involving the payload should be
// called afterwards.
func (self PeerNetwork) SendMsgDeferredByMove(msg Msg, peer PeerId) int32 {
	return int32(C.peernetwork_send_msg_deferred_by_move(self.inner, msg.inner, peer.inner))
}

// Send a message to the given list of peers. The payload contained in the
// given msg will be moved and sent. Thus, no methods of msg involving the
// payload should be called afterwards.
func (self PeerNetwork) MulticastMsgByMove(msg Msg, peers []PeerId) (res int32) {
	na := NewPeerIdArrayFromPeers(peers, false)
	res = int32(C.peernetwork_multicast_msg_by_move(self.inner, msg.inner, na.inner))
	na.Free()
	return res
}

// The C function pointer type which takes peernetwork_conn_t*, bool and void*
// (passing in the custom user data allocated by C.malloc) as parameters.
type PeerNetworkPeerCallback = C.peernetwork_peer_callback_t

// Register a connection handler invoked when p2p connection is established/broken.
func (self PeerNetwork) RegPeerHandler(callback PeerNetworkPeerCallback, userdata rawptr_t) {
	C.peernetwork_reg_peer_handler(self.inner, callback, userdata)
}

// The C function pointer type which takes netaddr_t*, x509_t* and void* (passing in the
// custom user data allocated by C.malloc) as parameters.
type PeerNetworkUnknownPeerCallback = C.peernetwork_unknown_peer_callback_t

// Register a connection handler invoked when a remote peer that is not in the
// list of known peers attempted to connect. By default configuration, the
// connection was rejected, and you can call AddPeer() to enroll this peer if
// you hope to establish the connection soon.
func (self PeerNetwork) RegUnknownPeerHandler(callback PeerNetworkUnknownPeerCallback, userdata rawptr_t) {
	C.peernetwork_reg_unknown_peer_handler(self.inner, callback, userdata)
}

// The C pointer type for a ClientNetwork handle.
type CClientNetwork = *C.clientnetwork_t
type clientNetwork struct{ inner CClientNetwork }

// The handle for a client-server network.
type ClientNetwork = *clientNetwork

func ClientNetworkFromC(ptr CClientNetwork) ClientNetwork {
	return &clientNetwork{inner: ptr}
}

// The C pointer type for a ClientNetworkConn handle.
type CClientNetworkConn = *C.clientnetwork_conn_t
type clientNetworkConn struct {
	inner    CClientNetworkConn
	autoFree bool
}

// The handle for a ClientNetwork connection.
type ClientNetworkConn = *clientNetworkConn

func ClientNetworkConnFromC(ptr CClientNetworkConn) ClientNetworkConn {
	return &clientNetworkConn{inner: ptr}
}

func clientNetworkConnSetFinalizer(res ClientNetworkConn, autoFree bool) {
	res.autoFree = autoFree
	if res != nil && autoFree {
		runtime.SetFinalizer(res, func(self ClientNetworkConn) { self.Free() })
	}
}

// Create a client-server message network handle.
func NewClientNetwork(ec EventContext, config MsgNetworkConfig, err *Error) ClientNetwork {
	res := ClientNetworkFromC(C.clientnetwork_new(ec.inner, config.inner, err))
	if res != nil {
		ec.attach(rawptr_t(res.inner), res)
		runtime.SetFinalizer(res, func(self ClientNetwork) { self.free() })
	}
	return res
}

func (self ClientNetwork) free() { C.clientnetwork_free(self.inner) }

// Use the ClientNetwork handle as a MsgNetwork handle (to invoke the methods
// inherited from MsgNetwork, such as RegHandler).
func (self ClientNetwork) AsMsgNetwork() MsgNetwork {
	return MsgNetworkFromC(C.clientnetwork_as_msgnetwork(self.inner))
}

// Use the MsgNetwork handle as a ClientNetwork handle (forcing the conversion).
func (self MsgNetwork) AsClientNetworkUnsafe() ClientNetwork {
	return ClientNetworkFromC(C.msgnetwork_as_clientnetwork_unsafe(self.inner))
}

// Create a MsgNetworkConn handle from a ClientNetworkConn (representing the same
// connection).
func NewMsgNetworkConnFromClientNetworkConn(conn ClientNetworkConn, autoFree bool) MsgNetworkConn {
	res := MsgNetworkConnFromC(C.msgnetwork_conn_new_from_clientnetwork_conn(conn.inner))
	msgNetworkConnSetFinalizer(res, autoFree)
	return res
}

// Create a ClientNetworkConn handle from a MsgNetworkConn (representing the same
// connection and forcing the conversion).
func NewClientNetworkConnFromMsgNetworkConnUnsafe(conn MsgNetworkConn, autoFree bool) ClientNetworkConn {
	res := ClientNetworkConnFromC(C.clientnetwork_conn_new_from_msgnetwork_conn_unsafe(conn.inner))
	clientNetworkConnSetFinalizer(res, autoFree)
	return res
}

// Make a copy of the connection handle.
func (self ClientNetworkConn) Copy(autoFree bool) ClientNetworkConn {
	res := ClientNetworkConnFromC(C.clientnetwork_conn_copy(self.inner))
	clientNetworkConnSetFinalizer(res, autoFree)
	return res
}

func (self ClientNetworkConn) Free() {
	C.clientnetwork_conn_free(self.inner)
	if self.autoFree {
		runtime.SetFinalizer(self, nil)
	}
}

// Send a message to the given client.
func (self ClientNetwork) SendMsg(msg Msg, addr NetAddr) bool {
	return bool(C.clientnetwork_send_msg(self.inner, msg.inner, addr.inner))
}

// Send a message to the given client, using a worker thread to seralize and put
// data to the send buffer. The payload contained in the given msg will be
// moved and sent. Thus, no methods of msg involving the payload should be
// called afterwards.
func (self ClientNetwork) SendMsgDeferredByMove(msg Msg, addr NetAddr) int32 {
	return int32(C.clientnetwork_send_msg_deferred_by_move(self.inner, msg.inner, addr.inner))
}
