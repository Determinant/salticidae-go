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

type MsgNetworkMsgCallback = C.msgnetwork_msg_callback_t
type MsgNetworkConnCallback = C.msgnetwork_conn_callback_t

func (self MsgNetwork) RegHandler(opcode Opcode, callback MsgNetworkMsgCallback) {
    C.msgnetwork_reg_handler(self, C._opcode_t(opcode), callback)
}

func (self MsgNetwork) RegConnHandler(callback MsgNetworkConnCallback) {
    C.msgnetwork_reg_conn_handler(self, callback)
}
