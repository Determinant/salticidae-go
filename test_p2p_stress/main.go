package main

// #cgo CFLAGS: -I${SRCDIR}/../salticidae/include/
// #include <stdlib.h>
// #include <arpa/inet.h>
// #include "salticidae/network.h"
// void onTerm(int sig, void *);
// void onReceiveRand(msg_t *, msgnetwork_conn_t *, void *);
// void onReceiveAck(msg_t *, msgnetwork_conn_t *, void *);
// void onStopLoop(threadcall_handle_t *, void *);
// void connHandler(msgnetwork_conn_t *, bool, void *);
// void errorHandler(SalticidaeCError *, bool, void *);
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
    "github.com/Determinant/salticidae-go"
    "math/rand"
    "os"
    "fmt"
    "sync"
    "unsafe"
    "strconv"
)

const (
    MSG_OPCODE_RAND salticidae.Opcode = iota
    MSG_OPCODE_ACK
)

func msgRandSerialize(size int) (salticidae.Msg, salticidae.UInt256) {
    buffer := make([]byte, size)
    _, err := rand.Read(buffer)
    if err != nil {
        panic("rand source failed")
    }
    serialized := salticidae.NewDataStreamFromBytes(buffer)
    hash := serialized.GetHash()
    return salticidae.NewMsgMovedFromByteArray(
        MSG_OPCODE_RAND,
        salticidae.NewByteArrayMovedFromDataStream(serialized)), hash
}

func msgRandUnserialize(msg salticidae.Msg) salticidae.DataStream {
    return msg.ConsumePayload()
}

func msgAckSerialize(hash salticidae.UInt256) salticidae.Msg {
    serialized := salticidae.NewDataStream()
    hash.Serialize(serialized)
    return salticidae.NewMsgMovedFromByteArray(
        MSG_OPCODE_ACK,
        salticidae.NewByteArrayMovedFromDataStream(serialized))
}

func msgAckUnserialize(msg salticidae.Msg) salticidae.UInt256 {
    p := msg.ConsumePayload()
    hash := salticidae.NewUInt256()
    hash.Unserialize(p)
    p.Free()
    return hash
}

func checkError(err *salticidae.Error) {
    if err.GetCode() != 0 {
        fmt.Printf("error during a sync call: %s\n", salticidae.StrError(err.GetCode()))
        os.Exit(1)
    }
}

type TestContext struct {
    timer salticidae.TimerEvent
    timer_ctx *C.struct_timeout_callback_context_t
    state int
    hash salticidae.UInt256
    ncompleted int
}

func (self TestContext) Free() {
    if self.timer != nil {
        self.timer.Free()
        C.free(unsafe.Pointer(self.timer_ctx))
    }
    if self.hash != nil {
        self.hash.Free()
    }
}

type AppContext struct {
    addr salticidae.NetAddr
    ec salticidae.EventContext
    net salticidae.PeerNetwork
    tcall salticidae.ThreadCall
    tc map[uint64] *TestContext
}

func (self AppContext) Free() {
    self.addr.Free()
    self.net.Free()
    self.tcall.Free()
    for _, tc:= range self.tc {
        tc.Free()
    }
    self.ec.Free()
}

func NewTestContext() TestContext {
    return TestContext { ncompleted: 0 }
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

func sendRand(size int, app *AppContext, conn salticidae.MsgNetworkConn) {
    msg, hash := msgRandSerialize(size)
    addr := conn.GetAddr()
    tc := app.getTC(addr2id(addr))
    addr.Free()
    if tc.hash != nil {
        salticidae.UInt256(tc.hash).Free()
    }
    tc.hash = hash
    app.net.AsMsgNetwork().SendMsgByMove(msg, conn)
}

var apps []AppContext
var threads sync.WaitGroup
var segBuffSize = 4096
var ec salticidae.EventContext

//export onTimeout
func onTimeout(_ *C.timerev_t, userdata unsafe.Pointer) {
    ctx := (*C.struct_timeout_callback_context_t)(userdata)
    app := &apps[ctx.app_id]
    tc := app.getTC(uint64(ctx.addr_id))
    tc.ncompleted++
    app.net.AsMsgNetwork().Terminate(salticidae.MsgNetworkConn(ctx.conn))
    var s string
    for addr_id, v := range app.tc {
        s += fmt.Sprintf(" %d(%d)", C.ntohs(C.ushort(addr_id >> 32)), v.ncompleted)
    }
    fmt.Printf("INFO: %d completed:%s\n", C.ntohs(C.ushort(app.addr.GetPort())), s)
}

//export onReceiveRand
func onReceiveRand(_msg *C.struct_msg_t, _conn *C.struct_msgnetwork_conn_t, userdata unsafe.Pointer) {
    msg := salticidae.Msg(_msg)
    bytes := msgRandUnserialize(msg)
    hash := bytes.GetHash()
    bytes.Free()
    conn := salticidae.MsgNetworkConn(_conn)
    net := conn.GetNet()
    ack := msgAckSerialize(hash)
    hash.Free()
    net.SendMsgByMove(ack, conn)
}

//export onReceiveAck
func onReceiveAck(_msg *C.struct_msg_t, _conn *C.struct_msgnetwork_conn_t, userdata unsafe.Pointer) {
    hash := msgAckUnserialize(salticidae.Msg(_msg))
    id := *(* int)(userdata)
    app := &apps[id]
    conn := salticidae.MsgNetworkConn(_conn)
    _addr := conn.GetAddr()
    addr := addr2id(_addr)
    _addr.Free()
    tc := app.getTC(addr)

    if !hash.IsEq(tc.hash) {
        //fmt.Printf("%s %s\n", hash.GetHex(), tc.hash.GetHex())
        panic("corrupted I/O!")
    }
    hash.Free()

    if tc.state == segBuffSize * 2 {
        sendRand(tc.state, app, conn)
        tc.state = -1
        ctx := C.timeout_callback_context_new()
        ctx.app_id = C.int(id)
        ctx.addr_id = C.uint64_t(addr)
        ctx.conn = (*C.struct_msgnetwork_conn_t)(conn.Copy())
        if tc.timer != nil {
            tc.timer.Free()
            salticidae.MsgNetworkConn(tc.timer_ctx.conn).Free()
            C.free(unsafe.Pointer(tc.timer_ctx))
        }
        tc.timer = salticidae.NewTimerEvent(app.ec, salticidae.TimerEventCallback(C.onTimeout), unsafe.Pointer(ctx))
        tc.timer_ctx = ctx
        t := rand.Float64() * 10
        tc.timer.Add(t)
        fmt.Printf("rand-bomboard phase, ending in %.2f secs\n", t)
    } else if tc.state == -1 {
        sendRand(rand.Int() % (segBuffSize * 10), app, conn)
    } else {
        tc.state++
        sendRand(tc.state, app, conn)
    }
}

//export connHandler
func connHandler(_conn *C.struct_msgnetwork_conn_t, connected C.bool, userdata unsafe.Pointer) {
    conn := salticidae.MsgNetworkConn(_conn)
    id := *(*int)(userdata)
    app := &apps[id]
    if connected {
        if conn.GetMode() == salticidae.CONN_MODE_ACTIVE {
            addr := conn.GetAddr()
            tc := app.getTC(addr2id(addr))
            addr.Free()
            tc.state = 1
            fmt.Printf("INFO: increasing phase\n")
            sendRand(tc.state, app, conn)
        }
    }
}

//export errorHandler
func errorHandler(_err *C.struct_SalticidaeCError, fatal C.bool, _ unsafe.Pointer) {
    err := (*salticidae.Error)(unsafe.Pointer(_err))
    s := "recoverable"
    if fatal { s = "fatal" }
    fmt.Printf("Captured %s error during an async call: %s\n", s, salticidae.StrError(err.GetCode()))
}

//export onStopLoop
func onStopLoop(_ *C.threadcall_handle_t, userdata unsafe.Pointer) {
    ec := salticidae.EventContext(userdata)
    ec.Stop()
}

//export onTerm
func onTerm(_ C.int, _ unsafe.Pointer) {
    for i, _ := range apps {
        a := &apps[i]
        a.tcall.AsyncCall(
            salticidae.ThreadCallCallback(C.onStopLoop),
            unsafe.Pointer(a.ec))
    }
    threads.Wait()
    ec.Stop()
}

func main() {
    ec = salticidae.NewEventContext()

    var addrs []salticidae.NetAddr
    for i := 0; i < 5; i++ {
        addrs = append(addrs,
            salticidae.NewAddrFromIPPortString("127.0.0.1:" + strconv.Itoa(12345 + i)))
    }
    netconfig := salticidae.NewPeerNetworkConfig()
    nc := netconfig.AsMsgNetworkConfig()
    nc.SegBuffSize(segBuffSize)
    nc.NWorker(2)
    netconfig.ConnTimeout(5)
    netconfig.PingPeriod(2)
    apps = make([]AppContext, len(addrs))
    ids := make([](*C.int), len(addrs))
    for i, addr := range addrs {
        ec := salticidae.NewEventContext()
        apps[i] = AppContext {
            addr: addr,
            ec: ec,
            net: salticidae.NewPeerNetwork(ec, netconfig),
            tcall: salticidae.NewThreadCall(ec),
            tc: make(map[uint64] *TestContext),
        }
        ids[i] = (*C.int)(C.malloc(C.sizeof_int))
        *ids[i] = C.int(i)
        _i := unsafe.Pointer(ids[i])
        net := apps[i].net.AsMsgNetwork()
        net.RegHandler(MSG_OPCODE_RAND, salticidae.MsgNetworkMsgCallback(C.onReceiveRand), _i)
        net.RegHandler(MSG_OPCODE_ACK, salticidae.MsgNetworkMsgCallback(C.onReceiveAck), _i)
        net.RegConnHandler(salticidae.MsgNetworkConnCallback(C.connHandler), _i)
        net.RegErrorHandler(salticidae.MsgNetworkErrorCallback(C.errorHandler), _i)
        net.Start()
    }
    netconfig.Free()

    threads.Add(len(apps))
    for i, _ := range apps {
        app_id := i
        go func() {
            err := salticidae.NewError()
            a := &apps[app_id]
            a.net.Listen(a.addr, &err)
            for _, addr := range addrs {
                if !addr.IsEq(a.addr) {
                    a.net.AddPeer(addr)
                }
            }
            a.ec.Dispatch()
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

    ev_int.Free()
    ev_term.Free()
    ec.Free()
}
