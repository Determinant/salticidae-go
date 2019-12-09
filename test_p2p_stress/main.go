package main

// #include <stdlib.h>
// #include <stdint.h>
// #include <arpa/inet.h>
// #include "salticidae/network.h"
// void onTerm(int sig, void *);
// void onReceiveRand(msg_t *, msgnetwork_conn_t *, void *);
// void onReceiveAck(msg_t *, msgnetwork_conn_t *, void *);
// void onStopLoop(threadcall_handle_t *, void *);
// void peerHandler(peernetwork_conn_t *, bool, void *);
// void errorHandler(SalticidaeCError *, bool, int32_t, void *);
// void onTimeout(timerev_t *, void *);
// typedef struct timeout_callback_context_t {
//     int app_id;
//     uint64_t addr_id;
//     msgnetwork_conn_t *conn;
// } timeout_callback_context_t;
// static timeout_callback_context_t *timeout_callback_context_new() {
//     timeout_callback_context_t *ctx = malloc(sizeof(timeout_callback_context_t));
//     ctx->conn = NULL;
//     return ctx;
// }
//
import "C"

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"unsafe"

	"github.com/ava-labs/salticidae-go"
)

const (
	MSG_OPCODE_RAND salticidae.Opcode = iota
	MSG_OPCODE_ACK
)

func msgRandSerialize(view uint32, size int) (salticidae.Msg, salticidae.UInt256) {
	serialized := salticidae.NewDataStream(false)
	defer serialized.Free()
	serialized.PutU32(salticidae.ToLittleEndianU32(view))
	buffer := make([]byte, size)
	_, err := rand.Read(buffer)
	if err != nil {
		panic("rand source failed")
	}
	serialized.PutData(buffer)
	ba := salticidae.NewByteArrayFromBytes(buffer, false)
	defer ba.Free()
	payload := salticidae.NewByteArrayMovedFromDataStream(serialized, false)
	defer payload.Free()
	return salticidae.NewMsgMovedFromByteArray(MSG_OPCODE_RAND, payload, false), ba.GetHash(true)
}

func msgRandUnserialize(msg salticidae.Msg) (view uint32, hash salticidae.UInt256) {
	payload := msg.GetPayloadByMove()
	succ := true
	view = salticidae.FromLittleEndianU32(payload.GetU32(&succ))
	ba := salticidae.NewByteArrayCopiedFromDataStream(payload, false)
	defer ba.Free()
	hash = ba.GetHash(false)
	return
}

func msgAckSerialize(view uint32, hash salticidae.UInt256) salticidae.Msg {
	serialized := salticidae.NewDataStream(false)
	defer serialized.Free()
	serialized.PutU32(salticidae.ToLittleEndianU32(view))
	hash.Serialize(serialized)
	payload := salticidae.NewByteArrayMovedFromDataStream(serialized, false)
	defer payload.Free()
	return salticidae.NewMsgMovedFromByteArray(MSG_OPCODE_ACK, payload, false)
}

func msgAckUnserialize(msg salticidae.Msg) (view uint32, hash salticidae.UInt256) {
	payload := msg.GetPayloadByMove()
	hash = salticidae.NewUInt256(false)
	succ := true
	view = salticidae.FromLittleEndianU32(payload.GetU32(&succ))
	hash.Unserialize(payload)
	return
}

func checkError(err *salticidae.Error) {
	if err.GetCode() != 0 {
		fmt.Printf("error during a sync call: %s\n", salticidae.StrError(err.GetCode()))
		os.Exit(1)
	}
}

type TestContext struct {
	timer      salticidae.TimerEvent
	timer_ctx  *C.struct_timeout_callback_context_t
	state      int
	view       uint32
	hash       salticidae.UInt256
	ncompleted int
}

type AppContext struct {
	addr  salticidae.NetAddr
	ec    salticidae.EventContext
	net   salticidae.PeerNetwork
	tcall salticidae.ThreadCall
	tc    map[uint64]*TestContext
}

func (self AppContext) Free() {
	for _, tc := range self.tc {
		if tc.timer != nil {
			C.free(unsafe.Pointer(tc.timer_ctx))
		}
	}
}

func NewTestContext() TestContext {
	return TestContext{view: 0, ncompleted: 0}
}

func addr2id(addr salticidae.NetAddr) uint64 {
	return uint64(addr.GetIP()) | (uint64(addr.GetPort()) << 32)
}

func (self AppContext) getTC(addr_id uint64) (_tc *TestContext) {
	if tc, ok := self.tc[addr_id]; ok {
		_tc = tc
	} else {
		_tc = new(TestContext)
		self.tc[addr_id] = _tc
	}
	return
}

func sendRand(size int, app *AppContext, conn salticidae.MsgNetworkConn, tc *TestContext) {
	msg, hash := msgRandSerialize(tc.view, size)
	defer msg.Free()
	tc.hash = hash
	app.net.AsMsgNetwork().SendMsg(msg, conn)
}

var apps []AppContext
var threads sync.WaitGroup
var segBuffSize = 4096
var ec salticidae.EventContext
var ids []*C.int

//export onTimeout
func onTimeout(_ *C.timerev_t, userdata unsafe.Pointer) {
	ctx := (*C.struct_timeout_callback_context_t)(userdata)
	app := &apps[int(ctx.app_id)]
	tc := app.getTC(uint64(ctx.addr_id))
	tc.ncompleted++
	app.net.AsMsgNetwork().Terminate(
		salticidae.MsgNetworkConnFromC(
			salticidae.CMsgNetworkConn(ctx.conn)))
	var s string
	for addr_id, v := range app.tc {
		s += fmt.Sprintf(" %d(%d)", C.ntohs(C.ushort(addr_id>>32)), v.ncompleted)
	}
	fmt.Printf("INFO: %d completed:%s\n", C.ntohs(C.ushort(app.addr.GetPort())), s)
}

//export onReceiveRand
func onReceiveRand(_msg *C.struct_msg_t, _conn *C.struct_msgnetwork_conn_t, _ unsafe.Pointer) {
	msg := salticidae.MsgFromC(salticidae.CMsg(_msg))
	conn := salticidae.MsgNetworkConnFromC(salticidae.CMsgNetworkConn(_conn))
	net := conn.GetNet()
	view, hash := msgRandUnserialize(msg)
	defer hash.Free()
	ack := msgAckSerialize(view, hash)
	defer ack.Free()
	net.SendMsg(ack, conn)
}

//export onReceiveAck
func onReceiveAck(_msg *C.struct_msg_t, _conn *C.struct_msgnetwork_conn_t, userdata unsafe.Pointer) {
	view, hash := msgAckUnserialize(salticidae.MsgFromC(salticidae.CMsg(_msg)))
	defer hash.Free()
	id := int(*(*C.int)(userdata))
	app := &apps[id]
	conn := salticidae.MsgNetworkConnFromC(salticidae.CMsgNetworkConn(_conn))
	pconn := salticidae.NewPeerNetworkConnFromMsgNetworkConnUnsafe(conn, false)
	defer pconn.Free()
	addr := pconn.GetPeerAddr(false)
	defer addr.Free()
	if addr.IsNull() {
		return
	}
	addrID := addr2id(addr)
	tc := app.getTC(addrID)

	if view != tc.view {
		fmt.Printf("dropping stale MsgAck\n")
		return
	}

	if !hash.IsEq(tc.hash) {
		panic("corrupted I/O!")
	}

	if tc.state == segBuffSize*2 {
		sendRand(tc.state, app, conn, tc)
		tc.state = -1
		ctx := C.timeout_callback_context_new()
		ctx.app_id = C.int(id)
		ctx.addr_id = C.uint64_t(addrID)
		ctx.conn = C.msgnetwork_conn_copy(_conn)
		if tc.timer != nil {
			C.msgnetwork_conn_free(tc.timer_ctx.conn)
			C.free(unsafe.Pointer(tc.timer_ctx))
			tc.timer.Del()
		}
		tc.timer = salticidae.NewTimerEvent(app.ec, salticidae.TimerEventCallback(C.onTimeout), unsafe.Pointer(ctx))
		tc.timer_ctx = ctx
		t := rand.Float64() * 10
		tc.timer.Add(t)
		fmt.Printf("rand-bomboard phase, ending in %.2f secs\n", t)
	} else if tc.state == -1 {
		sendRand(rand.Int()%(segBuffSize*10), app, conn, tc)
	} else {
		tc.state++
		sendRand(tc.state, app, conn, tc)
	}
}

//export peerHandler
func peerHandler(_conn *C.struct_peernetwork_conn_t, connected C.bool, userdata unsafe.Pointer) {
	if connected {
		pconn := salticidae.PeerNetworkConnFromC(salticidae.CPeerNetworkConn(_conn))
		conn := salticidae.NewMsgNetworkConnFromPeerNetworkConn(pconn, false)
		defer conn.Free()
		id := int(*(*C.int)(userdata))
		app := &apps[id]
		addr := pconn.GetPeerAddr(false)
		defer addr.Free()
		tc := app.getTC(addr2id(addr))
		tc.state = 1
		tc.view++
		sendRand(tc.state, app, conn, tc)
	}
}

//export errorHandler
func errorHandler(_err *C.struct_SalticidaeCError, fatal C.bool, asyncID C.int32_t, _ unsafe.Pointer) {
	err := (*salticidae.Error)(unsafe.Pointer(_err))
	s := "recoverable"
	if fatal {
		s = "fatal"
	}
	fmt.Printf("Captured %s error during an async call %d: %s\n", s, asyncID, salticidae.StrError(err.GetCode()))
}

//export onStopLoop
func onStopLoop(_ *C.threadcall_handle_t, userdata unsafe.Pointer) {
	id := int(*(*C.int)(userdata))
	ec := apps[id].ec
	ec.Stop()
}

//export onTerm
func onTerm(_ C.int, _ unsafe.Pointer) {
	for i, _ := range apps {
		a := &apps[i]
		a.tcall.AsyncCall(
			salticidae.ThreadCallCallback(C.onStopLoop),
			unsafe.Pointer(ids[i]))
	}
	threads.Wait()
	ec.Stop()
}

func main() {
	ec = salticidae.NewEventContext()
	err := salticidae.NewError()

	var addrs []salticidae.NetAddr
	for i := 0; i < 5; i++ {
		addrs = append(addrs,
			salticidae.NewNetAddrFromIPPortString("127.0.0.1:"+strconv.Itoa(12345+i), true, &err))
	}
	netconfig := salticidae.NewPeerNetworkConfig()
	nc := netconfig.AsMsgNetworkConfig()
	nc.SegBuffSize(segBuffSize)
	nc.NWorker(2)
	netconfig.ConnTimeout(5)
	netconfig.PingPeriod(2)
	apps = make([]AppContext, len(addrs))
	ids = make([](*C.int), len(addrs))
	for i, addr := range addrs {
		ec := salticidae.NewEventContext()
		net := salticidae.NewPeerNetwork(ec, netconfig, &err)
		checkError(&err)
		apps[i] = AppContext{
			addr:  addr,
			ec:    ec,
			net:   net,
			tcall: salticidae.NewThreadCall(ec),
			tc:    make(map[uint64]*TestContext),
		}
		ids[i] = (*C.int)(C.malloc(C.sizeof_int))
		*ids[i] = C.int(i)
		_i := unsafe.Pointer(ids[i])
		mnet := net.AsMsgNetwork()
		mnet.RegHandler(MSG_OPCODE_RAND, salticidae.MsgNetworkMsgCallback(C.onReceiveRand), _i)
		mnet.RegHandler(MSG_OPCODE_ACK, salticidae.MsgNetworkMsgCallback(C.onReceiveAck), _i)
		net.RegPeerHandler(salticidae.PeerNetworkPeerCallback(C.peerHandler), _i)
		mnet.RegErrorHandler(salticidae.MsgNetworkErrorCallback(C.errorHandler), _i)
		mnet.Start()
	}

	threads.Add(len(apps))
	for i, _ := range apps {
		app_id := i
		go func() {
			err := salticidae.NewError()
			a := &apps[app_id]
			a.net.Listen(a.addr, &err)
			checkError(&err)
			for _, addr := range addrs {
				if !addr.IsEq(a.addr) {
					a.net.AddPeer(addr)
				}
			}
			a.ec.Dispatch()
			a.net.AsMsgNetwork().Stop()
			a.Free()
			C.free(unsafe.Pointer(ids[app_id]))
			threads.Done()
		}()
	}

	ev_int := salticidae.NewSigEvent(ec, salticidae.SigEventCallback(C.onTerm), nil)
	ev_int.Add(salticidae.SIGINT)
	ev_term := salticidae.NewSigEvent(ec, salticidae.SigEventCallback(C.onTerm), nil)
	ev_term.Add(salticidae.SIGTERM)

	ec.Dispatch()
}
