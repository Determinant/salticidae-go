package salticidae

// #include <stdlib.h>
// #include "salticidae/network.h"
import "C"
import "runtime"

//// begin MsgNetwork def

// CMsgNetwork is the C pointer type for a MsgNetwork handle.
type CMsgNetwork = *C.msgnetwork_t
type msgNetwork struct {
	inner    CMsgNetwork
	autoFree bool
}

// MsgNetwork is the handle for a message network.
type MsgNetwork = *msgNetwork

// MsgNetworkFromC converts an existing C pointer into a go pointer. Notice that
// when the go object does *not* own the resource of the C pointer, so it is
// only valid to the extent in which the given C pointer is valid. The C memory
// will not be deallocated when the go object is finalized by GC. This applies
// to all other "FromC" functions.
func MsgNetworkFromC(ptr CMsgNetwork) MsgNetwork {
	return &msgNetwork{inner: ptr}
}

func msgNetworkSetFinalizer(res MsgNetwork, autoFree bool) {
	res.autoFree = autoFree
	if res != nil && autoFree {
		runtime.SetFinalizer(res, func(self MsgNetwork) { self.Free() })
	}
}

// Free the underlying C pointer manually.
func (net MsgNetwork) Free() {
	C.msgnetwork_free(net.inner)
	if net.autoFree {
		runtime.SetFinalizer(net, nil)
	}
}

//// end MsgNetwork def

//// begin MsgNetworkConn def

// CMsgNetworkConn is the C pointer type for a MsgNetworkConn handle.
type CMsgNetworkConn = *C.msgnetwork_conn_t
type msgNetworkConn struct {
	inner    CMsgNetworkConn
	autoFree bool
}

// MsgNetworkConn is the handle for a message network connection.
type MsgNetworkConn = *msgNetworkConn

// MsgNetworkConnFromC converts an existing C pointer into a go pointer, without
// managing the underlying C poiner.
func MsgNetworkConnFromC(ptr CMsgNetworkConn) MsgNetworkConn {
	return &msgNetworkConn{inner: ptr}
}

func msgNetworkConnSetFinalizer(res MsgNetworkConn, autoFree bool) {
	res.autoFree = autoFree
	if res != nil && autoFree {
		runtime.SetFinalizer(res, func(self MsgNetworkConn) { self.Free() })
	}
}

// Free manually frees the underlying C pointer.
func (conn MsgNetworkConn) Free() {
	C.msgnetwork_conn_free(conn.inner)
	if conn.autoFree {
		runtime.SetFinalizer(conn, nil)
	}
}

//// end MsgNetworkConn def

//// begin MsgNetworkConfig def

// CMsgNetworkConfig is the C pointer type for a MsgNetworkConfig object.
type CMsgNetworkConfig = *C.msgnetwork_config_t
type msgNetworkConfig struct {
	inner    CMsgNetworkConfig
	autoFree bool
}

// MsgNetworkConfig is the configuration for MsgNetwork.
type MsgNetworkConfig = *msgNetworkConfig

// MsgNetworkConfigFromC converts a C pointer into a go pointer.
func MsgNetworkConfigFromC(ptr CMsgNetworkConfig) MsgNetworkConfig {
	return &msgNetworkConfig{inner: ptr}
}

func msgNetworkConfigSetFinalizer(res MsgNetworkConfig, autoFree bool) {
	res.autoFree = autoFree
	if res != nil && autoFree {
		runtime.SetFinalizer(res, func(self MsgNetworkConfig) { self.Free() })
	}
}

// Free manually frees the underlying C pointer
func (config MsgNetworkConfig) Free() {
	C.msgnetwork_config_free(config.inner)
	if config.autoFree {
		runtime.SetFinalizer(config, nil)
	}
}

//// end MsgNetworkConfig def

//// begin handler def

// MsgNetworkMsgCallback is the C function pointer type which takes msg_t*,
// msgnetwork_conn_t* and void* (passing in the custom user data allocated by
// C.malloc) as parameters.
type MsgNetworkMsgCallback = C.msgnetwork_msg_callback_t

// MsgNetworkConnCallback is the C function pointer type which takes
// msgnetwork_conn_t*, bool (true for the connection is established, false for
// the connection is terminated) and void* as parameters.
type MsgNetworkConnCallback = C.msgnetwork_conn_callback_t

// MsgNetworkErrorCallback is the C function Pointer type which takes
// SalticidaeCError* and void* as parameters.
type MsgNetworkErrorCallback = C.msgnetwork_error_callback_t

//// end handler def

//// begin PeerNetwork def

// CPeerNetwork is the C pointer type for a PeerNetwork handle.
type CPeerNetwork = *C.peernetwork_t
type peerNetwork struct {
	inner    CPeerNetwork
	autoFree bool
}

// PeerNetwork is the handle for a peer-to-peer network.
type PeerNetwork = *peerNetwork

// PeerNetworkFromC converts an existing C pointer into a go pointer.
func PeerNetworkFromC(ptr CPeerNetwork) PeerNetwork {
	return &peerNetwork{inner: ptr}
}

func peerNetworkSetFinalizer(res PeerNetwork, autoFree bool) {
	res.autoFree = autoFree
	if res != nil && autoFree {
		runtime.SetFinalizer(res, func(self PeerNetwork) { self.Free() })
	}
}

// Free the underlying C pointer manually.
func (net PeerNetwork) Free() {
	C.peernetwork_free(net.inner)
	if net.autoFree {
		runtime.SetFinalizer(net, nil)
	}
}

//// end PeerNetwork def

//// begin PeerNetworkConn def

// CPeerNetworkConn is the C pointer type for a PeerNetworkConn handle.
type CPeerNetworkConn = *C.peernetwork_conn_t
type peerNetworkConn struct {
	inner    CPeerNetworkConn
	autoFree bool
}

// PeerNetworkConn is the handle for a PeerNetwork connection.
type PeerNetworkConn = *peerNetworkConn

// PeerNetworkConnFromC converts an existing C pointer into a go pointer.
func PeerNetworkConnFromC(ptr CPeerNetworkConn) PeerNetworkConn {
	return &peerNetworkConn{inner: ptr}
}

func peerNetworkConnSetFinalizer(res PeerNetworkConn, autoFree bool) {
	res.autoFree = autoFree
	if res != nil && autoFree {
		runtime.SetFinalizer(res, func(self PeerNetworkConn) { self.Free() })
	}
}

// Free the underlying C pointer manually.
func (conn PeerNetworkConn) Free() {
	C.peernetwork_conn_free(conn.inner)
	if conn.autoFree {
		runtime.SetFinalizer(conn, nil)
	}
}

//// end PeerNetworkConn def

//// begin PeerNetworkConfig def

// CPeerNetworkConfig is the C pointer type for a PeerNetworkConfig handle.
type CPeerNetworkConfig = *C.peernetwork_config_t
type peerNetworkConfig struct {
	inner    CPeerNetworkConfig
	autoFree bool
}

// PeerNetworkConfig is the configuration for PeerNetwork.
type PeerNetworkConfig = *peerNetworkConfig

// PeerNetworkConfigFromC converts a C pointer into a go pointer.
func PeerNetworkConfigFromC(ptr CPeerNetworkConfig) PeerNetworkConfig {
	return &peerNetworkConfig{inner: ptr}
}

func peerNetworkConfigSetFinalizer(res PeerNetworkConfig, autoFree bool) {
	res.autoFree = autoFree
	if res != nil && autoFree {
		runtime.SetFinalizer(res, func(self PeerNetworkConfig) { self.Free() })
	}
}

// Free manually frees the underlying C pointer
func (config PeerNetworkConfig) Free() {
	C.peernetwork_config_free(config.inner)
	if config.autoFree {
		runtime.SetFinalizer(config, nil)
	}
}

//// end PeerNetworkConfig def

//// begin PeerID def

// CPeerID is the C pointer type for a PeerID object.
type CPeerID = *C.peerid_t
type peerID struct {
	inner    CPeerID
	autoFree bool
}

// PeerID is the peer identity object.
type PeerID = *peerID

// PeerIDFromC converts an existing C pointer into a go pointer.
func PeerIDFromC(ptr CPeerID) PeerID {
	return &peerID{inner: ptr}
}

func peerIDSetFinalizer(res PeerID, autoFree bool) {
	res.autoFree = autoFree
	if res != nil && autoFree {
		runtime.SetFinalizer(res, func(self PeerID) { self.Free() })
	}
}

// Free the underlying C pointer manually.
func (pid PeerID) Free() {
	C.peerid_free(pid.inner)
	if pid.autoFree {
		runtime.SetFinalizer(pid, nil)
	}
}

//// end PeerID def

//// begin PeerIDArray def

// CPeerIDArray is the C pointer type for a PeerIDArray object.
type CPeerIDArray = *C.peerid_array_t
type peerIDArray struct {
	inner    CPeerIDArray
	autoFree bool
}

// PeerIDArray is an array of peer ids.
type PeerIDArray = *peerIDArray

// PeerIDArrayFromC converts a C pointer into a go pointer.
func PeerIDArrayFromC(ptr CPeerIDArray) PeerIDArray {
	return &peerIDArray{inner: ptr}
}

func peerIDArraySetFinalizer(res PeerIDArray, autoFree bool) {
	res.autoFree = autoFree
	if res != nil && autoFree {
		runtime.SetFinalizer(res, func(self PeerIDArray) { self.Free() })
	}
}

// Free the underlying C pointer manually.
func (pidarr PeerIDArray) Free() {
	C.peerid_array_free(pidarr.inner)
	if pidarr.autoFree {
		runtime.SetFinalizer(pidarr, nil)
	}
}

//// end PeerIDArray def

//// begin ClientNetwork def

// CClientNetwork is the C pointer type for a ClientNetwork handle.
type CClientNetwork = *C.clientnetwork_t
type clientNetwork struct {
	inner    CClientNetwork
	autoFree bool
}

// ClientNetwork is the handle for a client-server network.
type ClientNetwork = *clientNetwork

// ClientNetworkFromC converts a C pointer into a go pointer.
func ClientNetworkFromC(ptr CClientNetwork) ClientNetwork {
	return &clientNetwork{inner: ptr}
}

func clientNetworkSetFinalizer(res ClientNetwork, autoFree bool) {
	res.autoFree = autoFree
	if res != nil && autoFree {
		runtime.SetFinalizer(res, func(self ClientNetwork) { self.Free() })
	}
}

// Free manually frees the underlying C pointer.
func (net ClientNetwork) Free() {
	C.clientnetwork_free(net.inner)
	if net.autoFree {
		runtime.SetFinalizer(net, nil)
	}
}

//// end ClientNetwork def

//// begin ClientNetworkConn def

// CClientNetworkConn is the C pointer type for a ClientNetworkConn handle.
type CClientNetworkConn = *C.clientnetwork_conn_t
type clientNetworkConn struct {
	inner    CClientNetworkConn
	autoFree bool
}

// ClientNetworkConn is the handle for a ClientNetwork connection.
type ClientNetworkConn = *clientNetworkConn

// ClientNetworkConnFromC converts a C pointer into a go pointer.
func ClientNetworkConnFromC(ptr CClientNetworkConn) ClientNetworkConn {
	return &clientNetworkConn{inner: ptr}
}

func clientNetworkConnSetFinalizer(res ClientNetworkConn, autoFree bool) {
	res.autoFree = autoFree
	if res != nil && autoFree {
		runtime.SetFinalizer(res, func(self ClientNetworkConn) { self.Free() })
	}
}

// Free manually frees the underlying C pointer.
func (conn ClientNetworkConn) Free() {
	C.clientnetwork_conn_free(conn.inner)
	if conn.autoFree {
		runtime.SetFinalizer(conn, nil)
	}
}

//// end ClientNetwork def

//////// methods ////////

//// begin MsgNetworkConn methods

var (
	CONN_MODE_ACTIVE  = MsgNetworkConnMode(C.CONN_MODE_ACTIVE)
	CONN_MODE_PASSIVE = MsgNetworkConnMode(C.CONN_MODE_PASSIVE)
)

// MsgNetworkConnMode specifies the connection mode. CONN_MODE_ACTIVE: a
// connection established from the local side. CONN_MODE_PASSIVE: a connection
// established from the remote side. CONN_MODE_DEAD: a connection that is
// already closed.
type MsgNetworkConnMode = C.msgnetwork_conn_mode_t

// GetNet gets the corresponding MsgNetwork handle that manages this connection. The
// returned handle is only valid during the lifetime of this connection.
func (conn MsgNetworkConn) GetNet() MsgNetwork {
	return MsgNetworkFromC(C.msgnetwork_conn_get_net(conn.inner))
}

// GetMode gets the connection mode.
func (conn MsgNetworkConn) GetMode() MsgNetworkConnMode {
	return C.msgnetwork_conn_get_mode(conn.inner)
}

// GetAddr gets the address of the remote end of this connection. Use Copy() to
// make a copy of the address if you want to use the address object beyond the
// lifetime of the connection.
func (conn MsgNetworkConn) GetAddr() NetAddr {
	return NetAddrFromC(C.msgnetwork_conn_get_addr(conn.inner))
}

// IsTerminated checks if the connection has been terminated.
func (conn MsgNetworkConn) IsTerminated() bool {
	return bool(C.msgnetwork_conn_is_terminated(conn.inner))
}

// GetPeerCert gets the certificate of the remote end of this connection. Use
// Copy() to make a copy of the certificate if you want to use the certificate
// object beyond the lifetime of the connection.
func (conn MsgNetworkConn) GetPeerCert() X509 {
	return &x509{inner: C.msgnetwork_conn_get_peer_cert(conn.inner)}
}

// Make a copy of the object. This is required if you want to keep the
// connection passed as a callback parameter by other salticidae methods (such
// like MsgNetwork/PeerNetwork).
func (conn MsgNetworkConn) Copy(autoFree bool) MsgNetworkConn {
	res := MsgNetworkConnFromC(C.msgnetwork_conn_copy(conn.inner))
	msgNetworkConnSetFinalizer(res, autoFree)
	return res
}

//// end MsgNetworkConn methods

//// begin MsgNetworkConfig methods

// NewMsgNetworkConfig creates the configuration object with default settings.
func NewMsgNetworkConfig() MsgNetworkConfig {
	res := MsgNetworkConfigFromC(C.msgnetwork_config_new())
	msgNetworkConfigSetFinalizer(res, true)
	return res
}

// MaxMsgSize sets the maximum message length (NOTE: by default it is 1
// KBytes).
func (config MsgNetworkConfig) MaxMsgSize(size int) {
	C.msgnetwork_config_max_msg_size(config.inner, C.size_t(size))
}

// MaxMsgQueueSize sets the maximum message queue size (the queue for buffering
// received messages to be processed by message handlers).
func (config MsgNetworkConfig) MaxMsgQueueSize(size int) {
	C.msgnetwork_config_max_msg_queue_size(config.inner, C.size_t(size))
}

// BurstSize sets the number of consecutive read attempts in the message
// delivery queue.  Usually the default value is good enough. This is used to
// make the tradeoff between the event loop fairness and the amortization of
// syscall cost.
func (config MsgNetworkConfig) BurstSize(size int) {
	C.msgnetwork_config_burst_size(config.inner, C.size_t(size))
}

// MaxListenBacklog set the maximum backlogs (see POSIX TCP backlog).
func (config MsgNetworkConfig) MaxListenBacklog(backlog int) {
	C.msgnetwork_config_max_listen_backlog(config.inner, C.int(backlog))
}

// ConnServerTimeout sets the timeout for connecting to the remote (in
// seconds).
func (config MsgNetworkConfig) ConnServerTimeout(timeout float64) {
	C.msgnetwork_config_conn_server_timeout(config.inner, C.double(timeout))
}

// RecvChunkSize sets the size for an inbound data chunk (per read() syscall).
func (config MsgNetworkConfig) RecvChunkSize(size int) {
	C.msgnetwork_config_recv_chunk_size(config.inner, C.size_t(size))
}

// NWorker sets the number of worker threads.
func (config MsgNetworkConfig) NWorker(nworker int) {
	C.msgnetwork_config_nworker(config.inner, C.size_t(nworker))
}

// MaxSendBuffSize sets the maximum send buffer size.
func (config MsgNetworkConfig) MaxSendBuffSize(size int) {
	C.msgnetwork_config_max_send_buff_size(config.inner, C.size_t(size))
}

// MaxRecvBuffSize sets the maximum recv buffer size.
func (config MsgNetworkConfig) MaxRecvBuffSize(size int) {
	C.msgnetwork_config_max_recv_buff_size(config.inner, C.size_t(size))
}

// EnableTLS specifies whether to use SSL/TLS. When this flag is enabled,
// MsgNetwork uses TLSKey (or TLSKeyFile) or TLSCert (or TLSCertFile) to setup
// the underlying OpenSSL.
func (config MsgNetworkConfig) EnableTLS(enabled bool) {
	C.msgnetwork_config_enable_tls(config.inner, C.bool(enabled))
}

// TLSKeyFile loads the TLS key from a file. The file should be an unencrypted
// PEM file.  Use TLSKey() instead for more complex usage.
func (config MsgNetworkConfig) TLSKeyFile(fname string) {
	cStr := C.CString(fname)
	C.msgnetwork_config_tls_key_file(config.inner, cStr)
	C.free(RawPtr(cStr))
}

// TLSCertFile loads the TLS certificate from a file. The file should be an
// unencrypted (X509) PEM file.  Use TLSCert() instead for more complex usage.
func (config MsgNetworkConfig) TLSCertFile(fname string) {
	cStr := C.CString(fname)
	C.msgnetwork_config_tls_cert_file(config.inner, cStr)
	C.free(RawPtr(cStr))
}

// TLSKeyByMove loads the given TLS key. This overrides TLSKeyFile(). pkey will
// be moved and kept by the library. Thus, no methods of msg involving the
// payload should be called afterwards.
func (config MsgNetworkConfig) TLSKeyByMove(pkey PKey) {
	C.msgnetwork_config_tls_key_by_move(config.inner, pkey.inner)
}

/// end MsgNetworkConfig methods

//// begin MsgNetwork methods

// NewMsgNetwork creates a message network handle which is attached to given
// event loop.
func NewMsgNetwork(ec EventContext, config MsgNetworkConfig, err *Error) MsgNetwork {
	res := MsgNetworkFromC(C.msgnetwork_new(ec.inner, config.inner, err))
	if res != nil {
		ec.attach(RawPtr(res.inner), res)
	}
	msgNetworkSetFinalizer(res, true)
	return res
}

// Start the message network (by spawning worker threads). This should be
// called before using any other methods.
func (net MsgNetwork) Start() { C.msgnetwork_start(net.inner) }

// Listen to the specified network address.
func (net MsgNetwork) Listen(addr NetAddr, err *Error) {
	C.msgnetwork_listen(net.inner, addr.inner, err)
}

// Stop the message network. No other methods should be called after this.
func (net MsgNetwork) Stop() { C.msgnetwork_stop(net.inner) }

// SendMsg sends a message through the given connection.
func (net MsgNetwork) SendMsg(msg Msg, conn MsgNetworkConn) bool {
	return bool(C.msgnetwork_send_msg(net.inner, msg.inner, conn.inner))
}

// SendMsgDeferredByMove sends a message through the given connection, using a
// worker thread to seralize and put data to the send buffer. The payload
// contained in the given msg will be moved and sent. Thus, no methods of msg
// involving the payload should be called afterwards.
func (net MsgNetwork) SendMsgDeferredByMove(msg Msg, conn MsgNetworkConn) int32 {
	return int32(C.msgnetwork_send_msg_deferred_by_move(net.inner, msg.inner, conn.inner))
}

// ConnectSync tries to connect to a remote address. The connection handle is
// returned. The returned connection handle could be kept in your program.
func (net MsgNetwork) ConnectSync(addr NetAddr, autoFree bool, err *Error) MsgNetworkConn {
	res := MsgNetworkConnFromC(C.msgnetwork_connect_sync(net.inner, addr.inner, err))
	msgNetworkConnSetFinalizer(res, autoFree)
	return res
}

// Connect tries to connect to a remote address (async). It returns an id which
// could be used to identify the corresponding error in the error callback.
func (net MsgNetwork) Connect(addr NetAddr) int32 {
	return int32(C.msgnetwork_connect(net.inner, addr.inner))
}

// Terminate the given connection.
func (net MsgNetwork) Terminate(conn MsgNetworkConn) {
	C.msgnetwork_terminate(net.inner, conn.inner)
}

// RegHandler registers a message handler for the type of message identified by
// opcode. The callback function will be invoked upon the delivery of each
// message with the given opcode, by the thread of the event loop the
// MsgNetwork is attached to.
func (net MsgNetwork) RegHandler(opcode Opcode, callback MsgNetworkMsgCallback, userdata RawPtr) {
	C.msgnetwork_reg_handler(net.inner, C._opcode_t(opcode), callback, userdata)
}

// RegConnHandler registers a connection handler invoked when the connection
// state is changed.
func (net MsgNetwork) RegConnHandler(callback MsgNetworkConnCallback, userdata RawPtr) {
	C.msgnetwork_reg_conn_handler(net.inner, callback, userdata)
}

// RegErrorHandler registers an error handler invoked when there is recoverable
// errors during any asynchronous call/execution inside the MsgNetwork.
func (net MsgNetwork) RegErrorHandler(callback MsgNetworkErrorCallback, userdata RawPtr) {
	C.msgnetwork_reg_error_handler(net.inner, callback, userdata)
}

//// begin PeerNetworkConfig methods

var (
	ADDR_BASED = PeerNetworkIDMode(C.ID_MODE_ADDR_BASED)
	CERT_BASED = PeerNetworkIDMode(C.ID_MODE_CERT_BASED)
)

// PeerNetworkIDMode specifies the identity mode.  ID_MODE_IP_BASED: a remote
// peer is identified by the IP only. ID_MODE_IP_PORT_BASED: a remote peer is
// identified by IP + port number.
type PeerNetworkIDMode = C.peernetwork_id_mode_t

// NewPeerNetworkConfig creates the configuration object with default settings.
func NewPeerNetworkConfig() PeerNetworkConfig {
	res := PeerNetworkConfigFromC(C.peernetwork_config_new())
	peerNetworkConfigSetFinalizer(res, true)
	return res
}

// PingPeriod sets the period for sending ping messsages (in seconds).
func (config PeerNetworkConfig) PingPeriod(sec float64) {
	C.peernetwork_config_ping_period(config.inner, C.double(sec))
}

// ConnTimeout sets the time it takes after sending a ping message before a
// connection is considered as broken.
func (config PeerNetworkConfig) ConnTimeout(sec float64) {
	C.peernetwork_config_conn_timeout(config.inner, C.double(sec))
}

// IDMode sets the identity mode.
func (config PeerNetworkConfig) IDMode(mode PeerNetworkIDMode) {
	C.peernetwork_config_id_mode(config.inner, mode)
}

// AsMsgNetworkConfig uses the PeerNetworkConfig object as a MsgNetworkConfig
// object (to invoke the methods inherited from MsgNetworkConfig, such as
// NWorker).
func (config PeerNetworkConfig) AsMsgNetworkConfig() MsgNetworkConfig {
	return MsgNetworkConfigFromC(C.peernetwork_config_as_msgnetwork_config(config.inner))
}

//// end PeerNetworkConfig methods

//// begin PeerID methods

// NewPeerIDFromNetAddr creates a PeerID from the NetAddr.
func NewPeerIDFromNetAddr(addr NetAddr, autoFree bool) (res PeerID) {
	res = PeerIDFromC(C.peerid_new_from_netaddr(addr.inner))
	peerIDSetFinalizer(res, autoFree)
	return
}

// NewPeerIDFromX509 creates a PeerID from the X509 certificate.
func NewPeerIDFromX509(cert X509, autoFree bool) (res PeerID) {
	res = PeerIDFromC(C.peerid_new_from_x509(cert.inner))
	peerIDSetFinalizer(res, autoFree)
	return
}

// NewPeerIDMovedFromUInt256 creates a PeerID from raw id by moving.
func NewPeerIDMovedFromUInt256(movedRawID UInt256, autoFree bool) (res PeerID) {
	res = PeerIDFromC(C.peerid_new_moved_from_uint256(movedRawID.inner))
	peerIDSetFinalizer(res, autoFree)
	return
}

// NewPeerIDCopiedFromUInt256 creates a PeerID from raw id by copying.
func NewPeerIDCopiedFromUInt256(rawID UInt256, autoFree bool) (res PeerID) {
	res = PeerIDFromC(C.peerid_new_copied_from_uint256(rawID.inner))
	peerIDSetFinalizer(res, autoFree)
	return
}

// AsUInt256 treats PeerID as its underlying UInt256.
func (pid PeerID) AsUInt256() (res UInt256) {
	return UInt256FromC(C.peerid_as_uint256(pid.inner))
}

//// end PeerID methods

//// begin PeerIDArray methods

// NewPeerIDArrayFromPeers convert a Go slice of peer ids to PeerIDArray.
func NewPeerIDArrayFromPeers(arr []PeerID, autoFree bool) (res PeerIDArray) {
	size := len(arr)
	_arr := make([]CPeerID, size)
	for i, v := range arr {
		_arr[i] = v.inner
	}
	if size > 0 {
		// FIXME: here we assume struct of a single pointer has the same memory
		// footprint the pointer
		base := (**C.peerid_t)(RawPtr(&_arr[0]))
		res = PeerIDArrayFromC(C.peerid_array_new_from_peers(base, C.size_t(size)))
	} else {
		res = PeerIDArrayFromC(C.peerid_array_new())
	}
	runtime.KeepAlive(_arr)
	peerIDArraySetFinalizer(res, autoFree)
	return
}

//// end PeerIDArray methods

//// begin PeerNetwork methods

// NewPeerNetwork creates a peer-to-peer message network handle.
func NewPeerNetwork(ec EventContext, config PeerNetworkConfig, err *Error) PeerNetwork {
	res := PeerNetworkFromC(C.peernetwork_new(ec.inner, config.inner, err))
	if res != nil {
		ec.attach(RawPtr(res.inner), res)
		peerNetworkSetFinalizer(res, true)
	}
	return res
}

// Listen to the specified network address. Notice that this method overrides
// Listen() in MsgNetwork, so you should always call this one instead of
// AsMsgNetwork().Listen().
func (net PeerNetwork) Listen(listenAddr NetAddr, err *Error) {
	C.peernetwork_listen(net.inner, listenAddr.inner, err)
}

// AddPeer registers a peer to the list of known peers. The P2P network will
// try to keep bi-direction connections to all known peers in the list (through
// reconnection and connection deduplication). This is an async call and the
// call id is returned as the reference for error handling.
func (net PeerNetwork) AddPeer(peer PeerID) int32 {
	return int32(C.peernetwork_add_peer(net.inner, peer.inner))
}

// DelPeer removes a peer from the list of known peers. The P2P network will
// teardown the existing bi-direction connection and the PeerHandler callback
// will not be called. This is an async call.
func (net PeerNetwork) DelPeer(peer PeerID) int32 {
	return int32(C.peernetwork_del_peer(net.inner, peer.inner))
}

// HasPeer tests whether a peer is already in the list.
func (net PeerNetwork) HasPeer(peer PeerID) bool {
	return bool(C.peernetwork_has_peer(net.inner, peer.inner))
}

// GetPeerConn gets the connection of the known peer. The connection handle is
// returned. The returned connection handle could be kept in your program.
func (net PeerNetwork) GetPeerConn(peer PeerID, autoFree bool, err *Error) PeerNetworkConn {
	res := PeerNetworkConnFromC(C.peernetwork_get_peer_conn(net.inner, peer.inner, err))
	peerNetworkConnSetFinalizer(res, autoFree)
	return res
}

// SetPeerAddr sets the IP address of the registered peer, used to connect to
// the peer. The address for a peer is by default empty and a p2p connection
// can only be established from the other side in this case (which is common
// for the peers behind some firewall/router). This is an async call.
func (net PeerNetwork) SetPeerAddr(peer PeerID, addr NetAddr) int32 {
	return int32(C.peernetwork_set_peer_addr(net.inner, peer.inner, addr.inner))
}

// ConnPeer tries to connect to the peer. If ntry > 0, it specifies the maximum
// number of attempts before giving up. If ntry = 0, it stops any
// ongoing/established connection and future attempts.  If ntry = -1,
// reconnection is attempted indefinitely. retryDelay specifies the minimum
// delay (in seconds) between two attempts. When ntry != 0, once peer is
// connected, the retry state is reset with ntry.
func (net PeerNetwork) ConnPeer(peer PeerID, ntry int32, retryDelay float64) int32 {
	return int32(C.peernetwork_conn_peer(net.inner, peer.inner, C.int(ntry), C.double(retryDelay)))
}

// SendMsg sends a message to the given peer.
func (net PeerNetwork) SendMsg(msg Msg, peer PeerID) bool {
	return bool(C.peernetwork_send_msg(net.inner, msg.inner, peer.inner))
}

// SendMsgDeferredByMove sends a message to the given peer, using a worker
// thread to seralize and put data to the send buffer. The payload contained in
// the given msg will be moved and sent. Thus, no methods of msg involving the
// payload should be called afterwards.
func (net PeerNetwork) SendMsgDeferredByMove(msg Msg, peer PeerID) int32 {
	return int32(C.peernetwork_send_msg_deferred_by_move(net.inner, msg.inner, peer.inner))
}

// MulticastMsgByMove sends a message to the given list of peers. The payload
// contained in the given msg will be moved and sent. Thus, no methods of msg
// involving the payload should be called afterwards.
func (net PeerNetwork) MulticastMsgByMove(msg Msg, peers []PeerID) (res int32) {
	na := NewPeerIDArrayFromPeers(peers, false)
	res = int32(C.peernetwork_multicast_msg_by_move(net.inner, msg.inner, na.inner))
	na.Free()
	return res
}

// PeerNetworkPeerCallback is the C function pointer type which takes
// peernetwork_conn_t*, bool and void* (passing in the custom user data
// allocated by C.malloc) as parameters.
type PeerNetworkPeerCallback = C.peernetwork_peer_callback_t

// RegPeerHandler registers a connection handler invoked when p2p connection is
// established/broken.
func (net PeerNetwork) RegPeerHandler(callback PeerNetworkPeerCallback, userdata RawPtr) {
	C.peernetwork_reg_peer_handler(net.inner, callback, userdata)
}

// PeerNetworkUnknownPeerCallback is the C function pointer type which takes
// netaddr_t*, x509_t* and void* (passing in the custom user data allocated by
// C.malloc) as parameters.
type PeerNetworkUnknownPeerCallback = C.peernetwork_unknown_peer_callback_t

// RegUnknownPeerHandler registers a connection handler invoked when a remote
// peer that is not in the list of known peers attempted to connect. By default
// configuration, the connection was rejected, and you can call AddPeer() to
// enroll this peer if you hope to establish the connection soon.
func (net PeerNetwork) RegUnknownPeerHandler(callback PeerNetworkUnknownPeerCallback, userdata RawPtr) {
	C.peernetwork_reg_unknown_peer_handler(net.inner, callback, userdata)
}

// AsMsgNetwork uses the PeerNetwork handle as a MsgNetwork handle (to invoke
// the methods inherited from MsgNetwork, such as RegHandler).
func (net PeerNetwork) AsMsgNetwork() MsgNetwork {
	return MsgNetworkFromC(C.peernetwork_as_msgnetwork(net.inner))
}

// AsPeerNetworkUnsafe use the MsgNetwork handle as a PeerNetwork handle
// (forcing the conversion).
func (net MsgNetwork) AsPeerNetworkUnsafe() PeerNetwork {
	return PeerNetworkFromC(C.msgnetwork_as_peernetwork_unsafe(net.inner))
}

//// end PeerNetwork methods

//// begin PeerNetworkConn methods

// NewMsgNetworkConnFromPeerNetworkConn creates a MsgNetworkConn handle from a
// PeerNetworkConn (representing the same connection).
func NewMsgNetworkConnFromPeerNetworkConn(conn PeerNetworkConn, autoFree bool) MsgNetworkConn {
	res := MsgNetworkConnFromC(C.msgnetwork_conn_new_from_peernetwork_conn(conn.inner))
	msgNetworkConnSetFinalizer(res, autoFree)
	return res
}

// NewPeerNetworkConnFromMsgNetworkConnUnsafe creates a PeerNetworkConn handle
// from a MsgNetworkConn (representing the same connection and forcing the
// conversion).
func NewPeerNetworkConnFromMsgNetworkConnUnsafe(conn MsgNetworkConn, autoFree bool) PeerNetworkConn {
	res := PeerNetworkConnFromC(C.peernetwork_conn_new_from_msgnetwork_conn_unsafe(conn.inner))
	peerNetworkConnSetFinalizer(res, autoFree)
	return res
}

// Copy the connection handle.
func (conn PeerNetworkConn) Copy(autoFree bool) PeerNetworkConn {
	res := PeerNetworkConnFromC(C.peernetwork_conn_copy(conn.inner))
	peerNetworkConnSetFinalizer(res, autoFree)
	return res
}

// GetPeerAddr gets the listening address of the remote peer (no Copy() is needed).
func (conn PeerNetworkConn) GetPeerAddr(autoFree bool) NetAddr {
	res := NetAddrFromC(C.peernetwork_conn_get_peer_addr(conn.inner))
	netAddrSetFinalizer(res, autoFree)
	return res
}

// GetPeerID gets the id of the remote peer (no Copy() is needed).
func (conn PeerNetworkConn) GetPeerID(autoFree bool) PeerID {
	res := PeerIDFromC(C.peernetwork_conn_get_peer_id(conn.inner))
	peerIDSetFinalizer(res, autoFree)
	return res
}

//// end PeerNetworkConn methods

//// begin ClientNetwork methods

// NewClientNetwork creates a client-server message network handle.
func NewClientNetwork(ec EventContext, config MsgNetworkConfig, err *Error) ClientNetwork {
	res := ClientNetworkFromC(C.clientnetwork_new(ec.inner, config.inner, err))
	if res != nil {
		ec.attach(RawPtr(res.inner), res)
		clientNetworkSetFinalizer(res, true)
	}
	return res
}

// SendMsg sends a message to the given client.
func (net ClientNetwork) SendMsg(msg Msg, addr NetAddr) bool {
	return bool(C.clientnetwork_send_msg(net.inner, msg.inner, addr.inner))
}

// SendMsgDeferredByMove sends a message to the given client, using a worker
// thread to seralize and put data to the send buffer. The payload contained in
// the given msg will be moved and sent. Thus, no methods of msg involving the
// payload should be called afterwards.
func (net ClientNetwork) SendMsgDeferredByMove(msg Msg, addr NetAddr) int32 {
	return int32(C.clientnetwork_send_msg_deferred_by_move(net.inner, msg.inner, addr.inner))
}

// AsMsgNetwork uses the ClientNetwork handle as a MsgNetwork handle (to invoke
// the methods inherited from MsgNetwork, such as RegHandler).
func (net ClientNetwork) AsMsgNetwork() MsgNetwork {
	return MsgNetworkFromC(C.clientnetwork_as_msgnetwork(net.inner))
}

// AsClientNetworkUnsafe uses the MsgNetwork handle as a ClientNetwork handle
// (forcing the conversion).
func (net MsgNetwork) AsClientNetworkUnsafe() ClientNetwork {
	return ClientNetworkFromC(C.msgnetwork_as_clientnetwork_unsafe(net.inner))
}

//// end ClientNetwork methods

//// begin ClientNetworkConn methods

// NewMsgNetworkConnFromClientNetworkConn creates a MsgNetworkConn handle from
// a ClientNetworkConn (representing the same connection).
func NewMsgNetworkConnFromClientNetworkConn(conn ClientNetworkConn, autoFree bool) MsgNetworkConn {
	res := MsgNetworkConnFromC(C.msgnetwork_conn_new_from_clientnetwork_conn(conn.inner))
	msgNetworkConnSetFinalizer(res, autoFree)
	return res
}

// NewClientNetworkConnFromMsgNetworkConnUnsafe creates a ClientNetworkConn
// handle from a MsgNetworkConn (representing the same connection and forcing
// the conversion).
func NewClientNetworkConnFromMsgNetworkConnUnsafe(conn MsgNetworkConn, autoFree bool) ClientNetworkConn {
	res := ClientNetworkConnFromC(C.clientnetwork_conn_new_from_msgnetwork_conn_unsafe(conn.inner))
	clientNetworkConnSetFinalizer(res, autoFree)
	return res
}

// Copy the connection handle.
func (conn ClientNetworkConn) Copy(autoFree bool) ClientNetworkConn {
	res := ClientNetworkConnFromC(C.clientnetwork_conn_copy(conn.inner))
	clientNetworkConnSetFinalizer(res, autoFree)
	return res
}

//// end ClientNetworkConn methods
