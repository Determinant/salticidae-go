package main

// #cgo CFLAGS: -I${SRCDIR}/../salticidae/include/
// #include <stdlib.h>
// #include "salticidae/network.h"
// void onTerm(int sig, void *);
// void onBobStop(threadcall_handle_t *, void *);
// void onTrigger(threadcall_handle_t *, void *);
// bool connHandler(msgnetwork_conn_t *, bool, void *);
// void onReceiveBytes(msg_t *, msgnetwork_conn_t *, void *);
// void onPeriodStat(timerev_t *, void *);
import "C"

import (
    "os"
    "fmt"
    "unsafe"
    "github.com/Determinant/salticidae-go"
)

const (
    MSG_OPCODE_BYTES salticidae.Opcode = iota
)

func msgBytesSerialize(size int) (res salticidae.Msg) {
    serialized := salticidae.NewDataStream(false)
    serialized.PutU32(salticidae.ToLittleEndianU32(uint32(size)))
    serialized.PutData(make([]byte, size))
    ba := salticidae.NewByteArrayMovedFromDataStream(serialized, false)
    serialized.Free()
    res = salticidae.NewMsgMovedFromByteArray(MSG_OPCODE_BYTES, ba, false)
    ba.Free()
    return
}

func checkError(err *salticidae.Error) {
    if err.GetCode() != 0 {
        fmt.Printf("error during a sync call: %s\n", salticidae.StrError(err.GetCode()))
        os.Exit(1)
    }
}

type MyNet struct {
    id *C.int
    net salticidae.MsgNetwork
    conn salticidae.MsgNetworkConn
    name string
    evPeriodStat salticidae.TimerEvent
    tcall salticidae.ThreadCall
    nrecv uint32
    statTimeout float64
}

var (
    mynets []MyNet
    ec salticidae.EventContext
    tec salticidae.EventContext
    bobThread chan struct{}
)

//export onBobStop
func onBobStop(_ *C.threadcall_handle_t, userdata unsafe.Pointer) {
    tec.Stop()
}

//export onTerm
func onTerm(_ C.int, _ unsafe.Pointer) {
    bob := &mynets[1]
    bob.tcall.AsyncCall(salticidae.ThreadCallCallback(C.onBobStop), unsafe.Pointer(bob.id))
    ec.Stop()
    <-bobThread
}

//export onTrigger
func onTrigger(_ *C.threadcall_handle_t, userdata unsafe.Pointer) {
    id := *(*int)(userdata)
    mynet := &mynets[id]
    payload := msgBytesSerialize(256)
    mynet.net.SendMsg(payload, mynet.conn)
    payload.Free()
    if !mynet.conn.IsTerminated() {
        mynet.tcall.AsyncCall(salticidae.ThreadCallCallback(C.onTrigger), userdata)
    }
}

//export onReceiveBytes
func onReceiveBytes(_ *C.struct_msg_t, _conn *C.struct_msgnetwork_conn_t, userdata unsafe.Pointer) {
    id := *(*int)(userdata)
    mynet := &mynets[id]
    mynet.nrecv++
}

//export connHandler
func connHandler(_conn *C.struct_msgnetwork_conn_t, connected C.bool, userdata unsafe.Pointer) C.bool {
    conn := salticidae.MsgNetworkConnFromC(salticidae.CMsgNetworkConn(_conn))
    id := *(*int)(userdata)
    mynet := &mynets[id]
    if connected {
        if conn.GetMode() == salticidae.CONN_MODE_ACTIVE {
            fmt.Printf("[%s] connected, sending hello.\n", mynet.name)
            mynet.conn = conn.Copy()
            mynet.tcall.AsyncCall(salticidae.ThreadCallCallback(C.onTrigger), userdata)
        } else {
            fmt.Printf("[%s] passively connected, waiting for greetings.\n", mynet.name)
        }
    } else {
        fmt.Printf("[%s] disconnected, retrying.\n", mynet.name)
        err := salticidae.NewError()
        mynet.net.Connect(conn.GetAddr(), false, &err)
    }
    return true
}

//export onPeriodStat
func onPeriodStat(_ *C.timerev_t, userdata unsafe.Pointer) {
    id := *(*int)(userdata)
    mynet := &mynets[id]
    fmt.Printf("%.2f mps\n", float64(mynet.nrecv) / mynet.statTimeout)
    mynet.nrecv = 0
    mynet.evPeriodStat.Add(mynet.statTimeout)
}

func genMyNet(ec salticidae.EventContext, name string, statTimeout float64, _id int) MyNet {
    err := salticidae.NewError()
    nc := salticidae.NewMsgNetworkConfig()
    nc.QueueCapacity(65536)
    nc.BurstSize(1000)
    net := salticidae.NewMsgNetwork(ec, nc, &err); checkError(&err)
    id := (*C.int)(C.malloc(C.sizeof_int))
    *id = C.int(_id)
    n := MyNet {
        id: id,
        net: net,
        conn: nil,
        name: name,
        evPeriodStat: salticidae.NewTimerEvent(ec, salticidae.TimerEventCallback(C.onPeriodStat), unsafe.Pointer(id)),
        tcall: salticidae.NewThreadCall(ec),
        nrecv: 0,
        statTimeout: statTimeout,
    }
    n.net.RegHandler(MSG_OPCODE_BYTES, salticidae.MsgNetworkMsgCallback(C.onReceiveBytes), unsafe.Pointer(id))
    n.net.RegConnHandler(salticidae.MsgNetworkConnCallback(C.connHandler), unsafe.Pointer(id))
    if statTimeout > 0 {
        n.evPeriodStat.Add(0)
    }
    return n
}

func main() {
    ec = salticidae.NewEventContext()
    err := salticidae.NewError()

    aliceAddr := salticidae.NewAddrFromIPPortString("127.0.0.1:12345", &err)
    //bobAddr := salticidae.NewAddrFromIPPortString("127.0.0.1:12346", &err)

    mynets = append(mynets, genMyNet(ec, "alice", 10, 0))
    alice := &mynets[0]
    alice.net.Start()
    alice.net.Listen(aliceAddr, &err); checkError(&err)
    bobThread = make(chan struct{})

    tec = salticidae.NewEventContext()
    mynets = append(mynets, genMyNet(tec, "bob", -1, 1))

    go func() {
        bob := &mynets[1]
        bob.net.Start()
        bob.net.Connect(aliceAddr, false, &err); checkError(&err)
        tec.Dispatch()
        bobThread <-struct{}{}
    }()

    ev_int := salticidae.NewSigEvent(ec, salticidae.SigEventCallback(C.onTerm), nil)
    ev_int.Add(salticidae.SIGINT)
    ev_term := salticidae.NewSigEvent(ec, salticidae.SigEventCallback(C.onTerm), nil)
    ev_term.Add(salticidae.SIGTERM)

    ec.Dispatch()

    for i, _ := range mynets {
        mynets[i].net.Stop()
        C.free(unsafe.Pointer(mynets[i].id))
    }
}
