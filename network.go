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
	dep      interface{}
	autoFree bool
	freed    bool
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
	if res.inner != nil && autoFree {
		runtime.SetFinalizer(res, func(self MsgNetwork) { self.Free() })
	}
}

// Free the underlying C pointer manually.
func (net MsgNetwork) Free() {
	if doubleFreeWarn(&net.freed) {
		return
	}
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
	freed    bool
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
	if res.inner != nil && autoFree {
		runtime.SetFinalizer(res, func(self MsgNetworkConn) { self.Free() })
	}
}

// Free manually frees the underlying C pointer.
func (conn MsgNetworkConn) Free() {
	if doubleFreeWarn(&conn.freed) {
		return
	}
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
	dep      interface{}
	autoFree bool
	freed    bool
}

// MsgNetworkConfig is the configuration for MsgNetwork.
type MsgNetworkConfig = *msgNetworkConfig

// MsgNetworkConfigFromC converts a C pointer into a go pointer.
func MsgNetworkConfigFromC(ptr CMsgNetworkConfig) MsgNetworkConfig {
	return &msgNetworkConfig{inner: ptr}
}

func msgNetworkConfigSetFinalizer(res MsgNetworkConfig, autoFree bool) {
	res.autoFree = autoFree
	if res.inner != nil && autoFree {
		runtime.SetFinalizer(res, func(self MsgNetworkConfig) { self.Free() })
	}
}

// Free manually frees the underlying C pointer
func (config MsgNetworkConfig) Free() {
	if doubleFreeWarn(&config.freed) {
		return
	}
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
	dep      interface{}
	autoFree bool
	freed    bool
}

// PeerNetwork is the handle for a peer-to-peer network.
type PeerNetwork = *peerNetwork

// PeerNetworkFromC converts an existing C pointer into a go pointer.
func PeerNetworkFromC(ptr CPeerNetwork) PeerNetwork {
	return &peerNetwork{inner: ptr}
}

func peerNetworkSetFinalizer(res PeerNetwork, autoFree bool) {
	res.autoFree = autoFree
	if res.inner != nil && autoFree {
		runtime.SetFinalizer(res, func(self PeerNetwork) { self.Free() })
	}
}

// Free the underlying C pointer manually.
func (net PeerNetwork) Free() {
	if doubleFreeWarn(&net.freed) {
		return
	}
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
	freed    bool
}

// PeerNetworkConn is the handle for a PeerNetwork connection.
type PeerNetworkConn = *peerNetworkConn

// PeerNetworkConnFromC converts an existing C pointer into a go pointer.
func PeerNetworkConnFromC(ptr CPeerNetworkConn) PeerNetworkConn {
	return &peerNetworkConn{inner: ptr}
}

func peerNetworkConnSetFinalizer(res PeerNetworkConn, autoFree bool) {
	res.autoFree = autoFree
	if res.inner != nil && autoFree {
		runtime.SetFinalizer(res, func(self PeerNetworkConn) { self.Free() })
	}
}

// Free the underlying C pointer manually.
func (conn PeerNetworkConn) Free() {
	if doubleFreeWarn(&conn.freed) {
		return
	}
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
	freed    bool
}

// PeerNetworkConfig is the configuration for PeerNetwork.
type PeerNetworkConfig = *peerNetworkConfig

// PeerNetworkConfigFromC converts a C pointer into a go pointer.
func PeerNetworkConfigFromC(ptr CPeerNetworkConfig) PeerNetworkConfig {
	return &peerNetworkConfig{inner: ptr}
}

func peerNetworkConfigSetFinalizer(res PeerNetworkConfig, autoFree bool) {
	res.autoFree = autoFree
	if res.inner != nil && autoFree {
		runtime.SetFinalizer(res, func(self PeerNetworkConfig) { self.Free() })
	}
}

// Free manually frees the underlying C pointer
func (config PeerNetworkConfig) Free() {
	if doubleFreeWarn(&config.freed) {
		return
	}
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
	freed    bool
}

// PeerID is the peer identity object.
type PeerID = *peerID

// PeerIDFromC converts an existing C pointer into a go pointer.
func PeerIDFromC(ptr CPeerID) PeerID {
	return &peerID{inner: ptr}
}

func peerIDSetFinalizer(res PeerID, autoFree bool) {
	res.autoFree = autoFree
	if res.inner != nil && autoFree {
		runtime.SetFinalizer(res, func(self PeerID) { self.Free() })
	}
}

// Free the underlying C pointer manually.
func (pid PeerID) Free() {
	if doubleFreeWarn(&pid.freed) {
		return
	}
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
	freed    bool
}

// PeerIDArray is an array of peer ids.
type PeerIDArray = *peerIDArray

// PeerIDArrayFromC converts a C pointer into a go pointer.
func PeerIDArrayFromC(ptr CPeerIDArray) PeerIDArray {
	return &peerIDArray{inner: ptr}
}

func peerIDArraySetFinalizer(res PeerIDArray, autoFree bool) {
	res.autoFree = autoFree
	if res.inner != nil && autoFree {
		runtime.SetFinalizer(res, func(self PeerIDArray) { self.Free() })
	}
}

// Free the underlying C pointer manually.
func (pidarr PeerIDArray) Free() {
	if doubleFreeWarn(&pidarr.freed) {
		return
	}
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
	dep      interface{}
	autoFree bool
	freed    bool
}

// ClientNetwork is the handle for a client-server network.
type ClientNetwork = *clientNetwork

// ClientNetworkFromC converts a C pointer into a go pointer.
func ClientNetworkFromC(ptr CClientNetwork) ClientNetwork {
	return &clientNetwork{inner: ptr}
}

func clientNetworkSetFinalizer(res ClientNetwork, autoFree bool) {
	res.autoFree = autoFree
	if res.inner != nil && autoFree {
		runtime.SetFinalizer(res, func(self ClientNetwork) { self.Free() })
	}
}

// Free manually frees the underlying C pointer.
func (net ClientNetwork) Free() {
	if doubleFreeWarn(&net.freed) {
		return
	}
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
	freed    bool
}

// ClientNetworkConn is the handle for a ClientNetwork connection.
type ClientNetworkConn = *clientNetworkConn

// ClientNetworkConnFromC converts a C pointer into a go pointer.
func ClientNetworkConnFromC(ptr CClientNetworkConn) ClientNetworkConn {
	return &clientNetworkConn{inner: ptr}
}

func clientNetworkConnSetFinalizer(res ClientNetworkConn, autoFree bool) {
	res.autoFree = autoFree
	if res.inner != nil && autoFree {
		runtime.SetFinalizer(res, func(self ClientNetworkConn) { self.Free() })
	}
}

// Free manually frees the underlying C pointer.
func (conn ClientNetworkConn) Free() {
	if doubleFreeWarn(&conn.freed) {
		return
	}
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

// GetNet gets the corresponding MsgNetwork handle that manages this connection.
// The conn can only be GC'ed when res is no longer used.
func (conn MsgNetworkConn) GetNet() (res MsgNetwork) {
	res = MsgNetworkFromC(C.msgnetwork_conn_get_net(conn.inner))
	res.dep = conn
	return
}

// GetMode gets the connection mode.
func (conn MsgNetworkConn) GetMode() (res MsgNetworkConnMode) {
	res = C.msgnetwork_conn_get_mode(conn.inner)
	runtime.KeepAlive(conn)
	return
}

// GetAddr gets the address of the remote end of this connection. Use Copy() to
// make a copy of the address if you want to use the address object beyond the
// lifetime of the connection.
// The conn can only be GC'ed when res is no longer used.
func (conn MsgNetworkConn) GetAddr() (res NetAddr) {
	res = NetAddrFromC(C.msgnetwork_conn_get_addr(conn.inner))
	res.dep = conn
	return
}

// IsTerminated checks if the connection has been terminated.
func (conn MsgNetworkConn) IsTerminated() (res bool) {
	res = bool(C.msgnetwork_conn_is_terminated(conn.inner))
	runtime.KeepAlive(conn)
	return
}

// GetPeerCert gets the certificate of the remote end of this connection. Use
// Copy() to make a copy of the certificate if you want to use the certificate
// object beyond the lifetime of the connection.
// The conn can only be GC'ed when res is no longer used.
func (conn MsgNetworkConn) GetPeerCert() (res X509) {
	res = X509FromC(C.msgnetwork_conn_get_peer_cert(conn.inner))
	res.dep = conn
	return
}

// Make a copy of the object. This is required if you want to keep the
// connection passed as a callback parameter by other salticidae methods (such
// like MsgNetwork/PeerNetwork).
func (conn MsgNetworkConn) Copy(autoFree bool) (res MsgNetworkConn) {
	res = MsgNetworkConnFromC(C.msgnetwork_conn_copy(conn.inner))
	msgNetworkConnSetFinalizer(res, autoFree)
	runtime.KeepAlive(conn)
	return
}

//// end MsgNetworkConn methods

//// begin MsgNetworkConfig methods

// NewMsgNetworkConfig creates the configuration object with default settings.
func NewMsgNetworkConfig() (res MsgNetworkConfig) {
	res = MsgNetworkConfigFromC(C.msgnetwork_config_new())
	msgNetworkConfigSetFinalizer(res, true)
	return
}

// MaxMsgSize sets the maximum message length (NOTE: by default it is 1
// KBytes).
func (config MsgNetworkConfig) MaxMsgSize(size int) {
	C.msgnetwork_config_max_msg_size(config.inner, C.size_t(size))
	runtime.KeepAlive(config)
}

// MaxMsgQueueSize sets the maximum message queue size (the queue for buffering
// received messages to be processed by message handlers).
func (config MsgNetworkConfig) MaxMsgQueueSize(size int) {
	C.msgnetwork_config_max_msg_queue_size(config.inner, C.size_t(size))
	runtime.KeepAlive(config)
}

// BurstSize sets the number of consecutive read attempts in the message
// delivery queue.  Usually the default value is good enough. This is used to
// make the tradeoff between the event loop fairness and the amortization of
// syscall cost.
func (config MsgNetworkConfig) BurstSize(size int) {
	C.msgnetwork_config_burst_size(config.inner, C.size_t(size))
	runtime.KeepAlive(config)
}

// MaxListenBacklog set the maximum backlogs (see POSIX TCP backlog).
func (config MsgNetworkConfig) MaxListenBacklog(backlog int) {
	C.msgnetwork_config_max_listen_backlog(config.inner, C.int(backlog))
	runtime.KeepAlive(config)
}

// ConnServerTimeout sets the timeout for connecting to the remote (in
// seconds).
func (config MsgNetworkConfig) ConnServerTimeout(timeout float64) {
	C.msgnetwork_config_conn_server_timeout(config.inner, C.double(timeout))
	runtime.KeepAlive(config)
}

// RecvChunkSize sets the size for an inbound data chunk (per read() syscall).
func (config MsgNetworkConfig) RecvChunkSize(size int) {
	C.msgnetwork_config_recv_chunk_size(config.inner, C.size_t(size))
	runtime.KeepAlive(config)
}

// NWorker sets the number of worker threads.
func (config MsgNetworkConfig) NWorker(nworker int) {
	C.msgnetwork_config_nworker(config.inner, C.size_t(nworker))
	runtime.KeepAlive(config)
}

// MaxSendBuffSize sets the maximum send buffer size.
func (config MsgNetworkConfig) MaxSendBuffSize(size int) {
	C.msgnetwork_config_max_send_buff_size(config.inner, C.size_t(size))
	runtime.KeepAlive(config)
}

// MaxRecvBuffSize sets the maximum recv buffer size.
func (config MsgNetworkConfig) MaxRecvBuffSize(size int) {
	C.msgnetwork_config_max_recv_buff_size(config.inner, C.size_t(size))
	runtime.KeepAlive(config)
}

// EnableTLS specifies whether to use SSL/TLS. When this flag is enabled,
// MsgNetwork uses TLSKey (or TLSKeyFile) or TLSCert (or TLSCertFile) to setup
// the underlying OpenSSL.
func (config MsgNetworkConfig) EnableTLS(enabled bool) {
	C.msgnetwork_config_enable_tls(config.inner, C.bool(enabled))
	runtime.KeepAlive(config)
}

// TLSKeyFile loads the TLS key from a file. The file should be an unencrypted
// PEM file.  Use TLSKey() instead for more complex usage.
func (config MsgNetworkConfig) TLSKeyFile(fname string) {
	cStr := C.CString(fname)
	C.msgnetwork_config_tls_key_file(config.inner, cStr)
	C.free(rawPtr(cStr))
	runtime.KeepAlive(config)
}

// TLSCertFile loads the TLS certificate from a file. The file should be an
// unencrypted (X509) PEM file.  Use TLSCert() instead for more complex usage.
func (config MsgNetworkConfig) TLSCertFile(fname string) {
	cStr := C.CString(fname)
	C.msgnetwork_config_tls_cert_file(config.inner, cStr)
	C.free(rawPtr(cStr))
	runtime.KeepAlive(config)
}

// TLSKeyByMove loads the given TLS key. This overrides TLSKeyFile(). pkey will
// be moved and kept by the library. Thus, no methods of msg involving the
// payload should be called afterwards.
func (config MsgNetworkConfig) TLSKeyByMove(pkey PKey) {
	C.msgnetwork_config_tls_key_by_move(config.inner, pkey.inner)
	runtime.KeepAlive(config)
	runtime.KeepAlive(pkey)
}

/// end MsgNetworkConfig methods

//// begin MsgNetwork methods

// NewMsgNetwork creates a message network handle which is attached to given
// event loop.
func NewMsgNetwork(ec EventContext, config MsgNetworkConfig, err *Error) (res MsgNetwork) {
	res = MsgNetworkFromC(C.msgnetwork_new(ec.inner, config.inner, err))
	if res.inner != nil {
		ec.attach(rawPtr(res.inner), res)
	}
	msgNetworkSetFinalizer(res, true)
	runtime.KeepAlive(ec)
	runtime.KeepAlive(config)
	return
}

// Start the message network (by spawning worker threads). This should be
// called before using any other methods.
func (net MsgNetwork) Start() {
	C.msgnetwork_start(net.inner)
	runtime.KeepAlive(net)
}

// Listen to the specified network address.
func (net MsgNetwork) Listen(addr NetAddr, err *Error) {
	C.msgnetwork_listen(net.inner, addr.inner, err)
	runtime.KeepAlive(net)
	runtime.KeepAlive(addr)
}

// Stop the message network. No other methods should be called after this.
func (net MsgNetwork) Stop() {
	C.msgnetwork_stop(net.inner)
	runtime.KeepAlive(net)
}

// SendMsg sends a message through the given connection.
func (net MsgNetwork) SendMsg(msg Msg, conn MsgNetworkConn) (res bool) {
	res = bool(C.msgnetwork_send_msg(net.inner, msg.inner, conn.inner))
	runtime.KeepAlive(net)
	runtime.KeepAlive(msg)
	runtime.KeepAlive(conn)
	return
}

// SendMsgDeferredByMove sends a message through the given connection, using a
// worker thread to seralize and put data to the send buffer. The payload
// contained in the given msg will be moved and sent. Thus, no methods of msg
// involving the payload should be called afterwards.
func (net MsgNetwork) SendMsgDeferredByMove(msg Msg, conn MsgNetworkConn) (res int32) {
	res = int32(C.msgnetwork_send_msg_deferred_by_move(net.inner, msg.inner, conn.inner))
	runtime.KeepAlive(net)
	runtime.KeepAlive(msg)
	runtime.KeepAlive(conn)
	return
}

// ConnectSync tries to connect to a remote address. The connection handle is
// returned. The returned connection handle could be kept in your program.
func (net MsgNetwork) ConnectSync(addr NetAddr, autoFree bool, err *Error) (res MsgNetworkConn) {
	res = MsgNetworkConnFromC(C.msgnetwork_connect_sync(net.inner, addr.inner, err))
	msgNetworkConnSetFinalizer(res, autoFree)
	runtime.KeepAlive(net)
	runtime.KeepAlive(addr)
	return
}

// Connect tries to connect to a remote address (async). It returns an id which
// could be used to identify the corresponding error in the error callback.
func (net MsgNetwork) Connect(addr NetAddr) (res int32) {
	res = int32(C.msgnetwork_connect(net.inner, addr.inner))
	runtime.KeepAlive(net)
	runtime.KeepAlive(addr)
	return
}

// Terminate the given connection.
func (net MsgNetwork) Terminate(conn MsgNetworkConn) {
	C.msgnetwork_terminate(net.inner, conn.inner)
	runtime.KeepAlive(net)
	runtime.KeepAlive(conn)
}

// RegHandler registers a message handler for the type of message identified by
// opcode. The callback function will be invoked upon the delivery of each
// message with the given opcode, by the thread of the event loop the
// MsgNetwork is attached to.
func (net MsgNetwork) RegHandler(opcode Opcode, callback MsgNetworkMsgCallback, userdata rawPtr) {
	C.msgnetwork_reg_handler(net.inner, C._opcode_t(opcode), callback, userdata)
	runtime.KeepAlive(net)
}

// RegConnHandler registers a connection handler invoked when the connection
// state is changed.
func (net MsgNetwork) RegConnHandler(callback MsgNetworkConnCallback, userdata rawPtr) {
	C.msgnetwork_reg_conn_handler(net.inner, callback, userdata)
	runtime.KeepAlive(net)
}

// RegErrorHandler registers an error handler invoked when there is recoverable
// errors during any asynchronous call/execution inside the MsgNetwork.
func (net MsgNetwork) RegErrorHandler(callback MsgNetworkErrorCallback, userdata rawPtr) {
	C.msgnetwork_reg_error_handler(net.inner, callback, userdata)
	runtime.KeepAlive(net)
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
func NewPeerNetworkConfig() (res PeerNetworkConfig) {
	res = PeerNetworkConfigFromC(C.peernetwork_config_new())
	peerNetworkConfigSetFinalizer(res, true)
	return
}

// PingPeriod sets the period for sending ping messsages (in seconds).
func (config PeerNetworkConfig) PingPeriod(sec float64) {
	C.peernetwork_config_ping_period(config.inner, C.double(sec))
	runtime.KeepAlive(config)
}

// ConnTimeout sets the time it takes after sending a ping message before a
// connection is considered as broken.
func (config PeerNetworkConfig) ConnTimeout(sec float64) {
	C.peernetwork_config_conn_timeout(config.inner, C.double(sec))
	runtime.KeepAlive(config)
}

// IDMode sets the identity mode.
func (config PeerNetworkConfig) IDMode(mode PeerNetworkIDMode) {
	C.peernetwork_config_id_mode(config.inner, mode)
	runtime.KeepAlive(config)
}

// AsMsgNetworkConfig uses the PeerNetworkConfig object as a MsgNetworkConfig
// object (to invoke the methods inherited from MsgNetworkConfig, such as
// NWorker).
// The config can only be GC'ed when res is no longer used.
func (config PeerNetworkConfig) AsMsgNetworkConfig() (res MsgNetworkConfig) {
	res = MsgNetworkConfigFromC(C.peernetwork_config_as_msgnetwork_config(config.inner))
	res.dep = config
	return
}

//// end PeerNetworkConfig methods

//// begin PeerID methods

// NewPeerIDFromNetAddr creates a PeerID from the NetAddr.
func NewPeerIDFromNetAddr(addr NetAddr, autoFree bool) (res PeerID) {
	res = PeerIDFromC(C.peerid_new_from_netaddr(addr.inner))
	peerIDSetFinalizer(res, autoFree)
	runtime.KeepAlive(addr)
	return
}

// NewPeerIDFromX509 creates a PeerID from the X509 certificate.
func NewPeerIDFromX509(cert X509, autoFree bool) (res PeerID) {
	res = PeerIDFromC(C.peerid_new_from_x509(cert.inner))
	peerIDSetFinalizer(res, autoFree)
	runtime.KeepAlive(cert)
	return
}

// NewPeerIDMovedFromUInt256 creates a PeerID from raw id by moving.
func NewPeerIDMovedFromUInt256(movedRawID UInt256, autoFree bool) (res PeerID) {
	res = PeerIDFromC(C.peerid_new_moved_from_uint256(movedRawID.inner))
	peerIDSetFinalizer(res, autoFree)
	runtime.KeepAlive(movedRawID)
	return
}

// NewPeerIDCopiedFromUInt256 creates a PeerID from raw id by copying.
func NewPeerIDCopiedFromUInt256(rawID UInt256, autoFree bool) (res PeerID) {
	res = PeerIDFromC(C.peerid_new_copied_from_uint256(rawID.inner))
	peerIDSetFinalizer(res, autoFree)
	runtime.KeepAlive(rawID)
	return
}

// AsUInt256 treats PeerID as its underlying UInt256.
// The pid can only be GC'ed when res is no longer used.
func (pid PeerID) AsUInt256() (res UInt256) {
	res = UInt256FromC(C.peerid_as_uint256(pid.inner))
	res.dep = pid
	return
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
		base := (**C.peerid_t)(rawPtr(&_arr[0]))
		res = PeerIDArrayFromC(C.peerid_array_new_from_peers(base, C.size_t(size)))
	} else {
		res = PeerIDArrayFromC(C.peerid_array_new())
	}
	runtime.KeepAlive(arr)
	runtime.KeepAlive(_arr)
	peerIDArraySetFinalizer(res, autoFree)
	return
}

//// end PeerIDArray methods

//// begin PeerNetwork methods

// NewPeerNetwork creates a peer-to-peer message network handle.
func NewPeerNetwork(ec EventContext, config PeerNetworkConfig, err *Error) (res PeerNetwork) {
	res = PeerNetworkFromC(C.peernetwork_new(ec.inner, config.inner, err))
	if res.inner != nil {
		ec.attach(rawPtr(res.inner), res)
	}
	peerNetworkSetFinalizer(res, true)
	runtime.KeepAlive(ec)
	runtime.KeepAlive(config)
	return
}

// Listen to the specified network address. Notice that this method overrides
// Listen() in MsgNetwork, so you should always call this one instead of
// AsMsgNetwork().Listen().
func (net PeerNetwork) Listen(listenAddr NetAddr, err *Error) {
	C.peernetwork_listen(net.inner, listenAddr.inner, err)
	runtime.KeepAlive(net)
	runtime.KeepAlive(listenAddr)
}

// AddPeer registers a peer to the list of known peers. The P2P network will
// try to keep bi-direction connections to all known peers in the list (through
// reconnection and connection deduplication). This is an async call and the
// call id is returned as the reference for error handling.
func (net PeerNetwork) AddPeer(peer PeerID) (res int32) {
	res = int32(C.peernetwork_add_peer(net.inner, peer.inner))
	runtime.KeepAlive(net)
	runtime.KeepAlive(peer)
	return
}

// DelPeer removes a peer from the list of known peers. The P2P network will
// teardown the existing bi-direction connection and the PeerHandler callback
// will not be called. This is an async call.
func (net PeerNetwork) DelPeer(peer PeerID) (res int32) {
	res = int32(C.peernetwork_del_peer(net.inner, peer.inner))
	runtime.KeepAlive(net)
	runtime.KeepAlive(peer)
	return
}

// HasPeer tests whether a peer is already in the list.
func (net PeerNetwork) HasPeer(peer PeerID) (res bool) {
	res = bool(C.peernetwork_has_peer(net.inner, peer.inner))
	runtime.KeepAlive(net)
	runtime.KeepAlive(peer)
	return
}

// GetPeerConn gets the connection of the known peer. The connection handle is
// returned. The returned connection handle could be kept in your program.
func (net PeerNetwork) GetPeerConn(peer PeerID, autoFree bool, err *Error) (res PeerNetworkConn) {
	res = PeerNetworkConnFromC(C.peernetwork_get_peer_conn(net.inner, peer.inner, err))
	peerNetworkConnSetFinalizer(res, autoFree)
	runtime.KeepAlive(net)
	runtime.KeepAlive(peer)
	return
}

// SetPeerAddr sets the IP address of the registered peer, used to connect to
// the peer. The address for a peer is by default empty and a p2p connection
// can only be established from the other side in this case (which is common
// for the peers behind some firewall/router). This is an async call.
func (net PeerNetwork) SetPeerAddr(peer PeerID, addr NetAddr) (res int32) {
	res = int32(C.peernetwork_set_peer_addr(net.inner, peer.inner, addr.inner))
	runtime.KeepAlive(net)
	runtime.KeepAlive(peer)
	runtime.KeepAlive(addr)
	return
}

// ConnPeer tries to connect to the peer. If ntry > 0, it specifies the maximum
// number of attempts before giving up. If ntry = 0, it stops any
// ongoing/established connection and future attempts.  If ntry = -1,
// reconnection is attempted indefinitely. retryDelay specifies the minimum
// delay (in seconds) between two attempts. When ntry != 0, once peer is
// connected, the retry state is reset with ntry.
func (net PeerNetwork) ConnPeer(peer PeerID, ntry int32, retryDelay float64) (res int32) {
	res = int32(C.peernetwork_conn_peer(net.inner, peer.inner, C.int(ntry), C.double(retryDelay)))
	runtime.KeepAlive(net)
	runtime.KeepAlive(peer)
	return
}

// SendMsg sends a message to the given peer.
func (net PeerNetwork) SendMsg(msg Msg, peer PeerID) (res bool) {
	res = bool(C.peernetwork_send_msg(net.inner, msg.inner, peer.inner))
	runtime.KeepAlive(net)
	runtime.KeepAlive(msg)
	runtime.KeepAlive(peer)
	return
}

// SendMsgDeferredByMove sends a message to the given peer, using a worker
// thread to seralize and put data to the send buffer. The payload contained in
// the given msg will be moved and sent. Thus, no methods of msg involving the
// payload should be called afterwards.
func (net PeerNetwork) SendMsgDeferredByMove(msg Msg, peer PeerID) (res int32) {
	res = int32(C.peernetwork_send_msg_deferred_by_move(net.inner, msg.inner, peer.inner))
	runtime.KeepAlive(net)
	runtime.KeepAlive(msg)
	runtime.KeepAlive(peer)
	return
}

// MulticastMsgByMove sends a message to the given list of peers. The payload
// contained in the given msg will be moved and sent. Thus, no methods of msg
// involving the payload should be called afterwards.
func (net PeerNetwork) MulticastMsgByMove(msg Msg, peers []PeerID) (res int32) {
	na := NewPeerIDArrayFromPeers(peers, false)
	res = int32(C.peernetwork_multicast_msg_by_move(net.inner, msg.inner, na.inner))
	na.Free()
	runtime.KeepAlive(net)
	runtime.KeepAlive(msg)
	runtime.KeepAlive(na)
	return
}

// PeerNetworkPeerCallback is the C function pointer type which takes
// peernetwork_conn_t*, bool and void* (passing in the custom user data
// allocated by C.malloc) as parameters.
type PeerNetworkPeerCallback = C.peernetwork_peer_callback_t

// RegPeerHandler registers a connection handler invoked when p2p connection is
// established/broken.
func (net PeerNetwork) RegPeerHandler(callback PeerNetworkPeerCallback, userdata rawPtr) {
	C.peernetwork_reg_peer_handler(net.inner, callback, userdata)
	runtime.KeepAlive(net)
}

// PeerNetworkUnknownPeerCallback is the C function pointer type which takes
// netaddr_t*, x509_t* and void* (passing in the custom user data allocated by
// C.malloc) as parameters.
type PeerNetworkUnknownPeerCallback = C.peernetwork_unknown_peer_callback_t

// RegUnknownPeerHandler registers a connection handler invoked when a remote
// peer that is not in the list of known peers attempted to connect. By default
// configuration, the connection was rejected, and you can call AddPeer() to
// enroll this peer if you hope to establish the connection soon.
func (net PeerNetwork) RegUnknownPeerHandler(callback PeerNetworkUnknownPeerCallback, userdata rawPtr) {
	C.peernetwork_reg_unknown_peer_handler(net.inner, callback, userdata)
	runtime.KeepAlive(net)
}

// AsMsgNetwork uses the PeerNetwork handle as a MsgNetwork handle (to invoke
// the methods inherited from MsgNetwork, such as RegHandler).
// The net can only be GC'ed when res is no longer used.
func (net PeerNetwork) AsMsgNetwork() (res MsgNetwork) {
	res = MsgNetworkFromC(C.peernetwork_as_msgnetwork(net.inner))
	res.dep = net
	return
}

// AsPeerNetworkUnsafe use the MsgNetwork handle as a PeerNetwork handle
// (forcing the conversion).
// The net can only be GC'ed when res is no longer used.
func (net MsgNetwork) AsPeerNetworkUnsafe() (res PeerNetwork) {
	res = PeerNetworkFromC(C.msgnetwork_as_peernetwork_unsafe(net.inner))
	res.dep = net
	return
}

//// end PeerNetwork methods

//// begin PeerNetworkConn methods

// NewMsgNetworkConnFromPeerNetworkConn creates a MsgNetworkConn handle from a
// PeerNetworkConn (representing the same connection).
func NewMsgNetworkConnFromPeerNetworkConn(conn PeerNetworkConn, autoFree bool) (res MsgNetworkConn) {
	res = MsgNetworkConnFromC(C.msgnetwork_conn_new_from_peernetwork_conn(conn.inner))
	msgNetworkConnSetFinalizer(res, autoFree)
	runtime.KeepAlive(conn)
	return
}

// NewPeerNetworkConnFromMsgNetworkConnUnsafe creates a PeerNetworkConn handle
// from a MsgNetworkConn (representing the same connection and forcing the
// conversion).
func NewPeerNetworkConnFromMsgNetworkConnUnsafe(conn MsgNetworkConn, autoFree bool) (res PeerNetworkConn) {
	res = PeerNetworkConnFromC(C.peernetwork_conn_new_from_msgnetwork_conn_unsafe(conn.inner))
	peerNetworkConnSetFinalizer(res, autoFree)
	runtime.KeepAlive(conn)
	return
}

// Copy the connection handle.
func (conn PeerNetworkConn) Copy(autoFree bool) (res PeerNetworkConn) {
	res = PeerNetworkConnFromC(C.peernetwork_conn_copy(conn.inner))
	peerNetworkConnSetFinalizer(res, autoFree)
	runtime.KeepAlive(conn)
	return
}

// GetPeerAddr gets the listening address of the remote peer (no Copy() is needed).
func (conn PeerNetworkConn) GetPeerAddr(autoFree bool) (res NetAddr) {
	res = NetAddrFromC(C.peernetwork_conn_get_peer_addr(conn.inner))
	netAddrSetFinalizer(res, autoFree)
	runtime.KeepAlive(conn)
	return
}

// GetPeerID gets the id of the remote peer (no Copy() is needed).
func (conn PeerNetworkConn) GetPeerID(autoFree bool) (res PeerID) {
	res = PeerIDFromC(C.peernetwork_conn_get_peer_id(conn.inner))
	peerIDSetFinalizer(res, autoFree)
	runtime.KeepAlive(conn)
	return
}

//// end PeerNetworkConn methods

//// begin ClientNetwork methods

// NewClientNetwork creates a client-server message network handle.
func NewClientNetwork(ec EventContext, config MsgNetworkConfig, err *Error) (res ClientNetwork) {
	res = ClientNetworkFromC(C.clientnetwork_new(ec.inner, config.inner, err))
	if res.inner != nil {
		ec.attach(rawPtr(res.inner), res)
	}
	clientNetworkSetFinalizer(res, true)
	runtime.KeepAlive(ec)
	runtime.KeepAlive(config)
	return
}

// SendMsg sends a message to the given client.
func (net ClientNetwork) SendMsg(msg Msg, addr NetAddr) (res bool) {
	res = bool(C.clientnetwork_send_msg(net.inner, msg.inner, addr.inner))
	runtime.KeepAlive(net)
	runtime.KeepAlive(msg)
	runtime.KeepAlive(addr)
	return
}

// SendMsgDeferredByMove sends a message to the given client, using a worker
// thread to seralize and put data to the send buffer. The payload contained in
// the given msg will be moved and sent. Thus, no methods of msg involving the
// payload should be called afterwards.
func (net ClientNetwork) SendMsgDeferredByMove(msg Msg, addr NetAddr) (res int32) {
	res = int32(C.clientnetwork_send_msg_deferred_by_move(net.inner, msg.inner, addr.inner))
	runtime.KeepAlive(net)
	runtime.KeepAlive(msg)
	runtime.KeepAlive(addr)
	return
}

// AsMsgNetwork uses the ClientNetwork handle as a MsgNetwork handle (to invoke
// the methods inherited from MsgNetwork, such as RegHandler).
// The net can only be GC'ed when res is no longer used.
func (net ClientNetwork) AsMsgNetwork() (res MsgNetwork) {
	res = MsgNetworkFromC(C.clientnetwork_as_msgnetwork(net.inner))
	res.dep = net
	return
}

// AsClientNetworkUnsafe uses the MsgNetwork handle as a ClientNetwork handle
// (forcing the conversion).
// The net can only be GC'ed when res is no longer used.
func (net MsgNetwork) AsClientNetworkUnsafe() (res ClientNetwork) {
	res = ClientNetworkFromC(C.msgnetwork_as_clientnetwork_unsafe(net.inner))
	res.dep = net
	return
}

//// end ClientNetwork methods

//// begin ClientNetworkConn methods

// NewMsgNetworkConnFromClientNetworkConn creates a MsgNetworkConn handle from
// a ClientNetworkConn (representing the same connection).
func NewMsgNetworkConnFromClientNetworkConn(conn ClientNetworkConn, autoFree bool) (res MsgNetworkConn) {
	res = MsgNetworkConnFromC(C.msgnetwork_conn_new_from_clientnetwork_conn(conn.inner))
	msgNetworkConnSetFinalizer(res, autoFree)
	runtime.KeepAlive(conn)
	return
}

// NewClientNetworkConnFromMsgNetworkConnUnsafe creates a ClientNetworkConn
// handle from a MsgNetworkConn (representing the same connection and forcing
// the conversion).
func NewClientNetworkConnFromMsgNetworkConnUnsafe(conn MsgNetworkConn, autoFree bool) (res ClientNetworkConn) {
	res = ClientNetworkConnFromC(C.clientnetwork_conn_new_from_msgnetwork_conn_unsafe(conn.inner))
	clientNetworkConnSetFinalizer(res, autoFree)
	runtime.KeepAlive(conn)
	return
}

// Copy the connection handle.
func (conn ClientNetworkConn) Copy(autoFree bool) (res ClientNetworkConn) {
	res = ClientNetworkConnFromC(C.clientnetwork_conn_copy(conn.inner))
	clientNetworkConnSetFinalizer(res, autoFree)
	runtime.KeepAlive(conn)
	return
}

//// end ClientNetworkConn methods
