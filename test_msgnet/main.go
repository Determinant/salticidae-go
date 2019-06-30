package main

// #cgo CFLAGS: -I${SRCDIR}/../salticidae/include/
// #include <stdlib.h>
// #include "salticidae/network.h"
// void onTerm(int sig, void *);
// void onReceiveHello(msg_t *, msgnetwork_conn_t *, void *);
// void onReceiveAck(msg_t *, msgnetwork_conn_t *, void *);
// bool connHandler(msgnetwork_conn_t *, bool, void *);
// void errorHandler(SalticidaeCError *, bool, void *);
import "C"

import (
    "os"
    "fmt"
    "unsafe"
    "github.com/Determinant/salticidae-go"
)

const (
    MSG_OPCODE_HELLO salticidae.Opcode = iota
    MSG_OPCODE_ACK
)

func msgHelloSerialize(name string, text string) salticidae.Msg {
    serialized := salticidae.NewDataStream(true)
    serialized.PutU32(salticidae.ToLittleEndianU32(uint32(len(name))))
    serialized.PutData([]byte(name))
    serialized.PutData([]byte(text))
    return salticidae.NewMsgMovedFromByteArray(
        MSG_OPCODE_HELLO, salticidae.NewByteArrayMovedFromDataStream(serialized, true), true)
}

func msgHelloUnserialize(msg salticidae.Msg) (name string, text string) {
    p := msg.GetPayloadByMove()
    succ := true
    length := salticidae.FromLittleEndianU32(p.GetU32(&succ))
    t := p.GetDataInPlace(int(length)); name = string(t.Get()); t.Release()
    t = p.GetDataInPlace(p.Size()); text = string(t.Get()); t.Release()
    return
}

func msgAckSerialize() salticidae.Msg {
    return salticidae.NewMsgMovedFromByteArray(MSG_OPCODE_ACK, salticidae.NewByteArray(true), true)
}

func checkError(err *salticidae.Error) {
    if err.GetCode() != 0 {
        fmt.Printf("error during a sync call: %s\n", salticidae.StrError(err.GetCode()))
        os.Exit(1)
    }
}

type MyNet struct {
    net salticidae.MsgNetwork
    name string
    cname *C.char
}

var (
    alice, bob MyNet
    ec salticidae.EventContext
)

//export onTerm
func onTerm(_ C.int, _ unsafe.Pointer) {
    ec.Stop()
}

//export onReceiveHello
func onReceiveHello(_msg *C.struct_msg_t, _conn *C.struct_msgnetwork_conn_t, userdata unsafe.Pointer) {
    msg := salticidae.MsgFromC(salticidae.CMsg(_msg))
    conn := salticidae.MsgNetworkConnFromC(salticidae.CMsgNetworkConn(_conn))
    net := conn.GetNet()
    name, text := msgHelloUnserialize(msg)
    myName := C.GoString((*C.char)(userdata))
    fmt.Printf("[%s] %s says %s\n", myName, name, text)
    ack := msgAckSerialize()
    net.SendMsg(ack, conn)
}

//export onReceiveAck
func onReceiveAck(_ *C.struct_msg_t, _conn *C.struct_msgnetwork_conn_t, userdata unsafe.Pointer) {
    myName := C.GoString((*C.char)(userdata))
    fmt.Printf("[%s] the peer knows\n", myName)
}

//export connHandler
func connHandler(_conn *C.struct_msgnetwork_conn_t, connected C.bool, userdata unsafe.Pointer) C.bool {
    conn := salticidae.MsgNetworkConnFromC(salticidae.CMsgNetworkConn(_conn))
    net := conn.GetNet()
    myName := C.GoString((*C.char)(userdata))
    n := alice
    if myName == "bob" { n = bob }
    if connected {
        if conn.GetMode() == salticidae.CONN_MODE_ACTIVE {
            fmt.Printf("[%s] connected, sending hello.\n", myName)
            hello := msgHelloSerialize(myName, "Hello there!")
            n.net.SendMsg(hello, conn)
        } else {
            fmt.Printf("[%s] accepted, waiting for greetings.\n", myName)
        }
    } else {
        fmt.Printf("[%s] disconnected, retrying.\n", myName)
        err := salticidae.NewError()
        net.Connect(conn.GetAddr(), false, &err)
    }
    return true
}

//export errorHandler
func errorHandler(_err *C.struct_SalticidaeCError, fatal C.bool, _ unsafe.Pointer) {
    err := (*salticidae.Error)(unsafe.Pointer(_err))
    s := "recoverable"
    if fatal { s = "fatal" }
    fmt.Printf("Captured %s error during an async call: %s\n", s, salticidae.StrError(err.GetCode()))
}

func genMyNet(ec salticidae.EventContext,
            name string,
            myAddr salticidae.NetAddr, otherAddr salticidae.NetAddr) MyNet {
    err := salticidae.NewError()
    netconfig := salticidae.NewMsgNetworkConfig()
    net := salticidae.NewMsgNetwork(ec, netconfig, &err); checkError(&err)
    n := MyNet { net: net, name: name, cname: C.CString(name) }
    cname := unsafe.Pointer(n.cname)
    n.net.RegHandler(MSG_OPCODE_HELLO, salticidae.MsgNetworkMsgCallback(C.onReceiveHello), cname)
    n.net.RegHandler(MSG_OPCODE_ACK, salticidae.MsgNetworkMsgCallback(C.onReceiveAck), cname)
    n.net.RegConnHandler(salticidae.MsgNetworkConnCallback(C.connHandler), cname)
    n.net.RegErrorHandler(salticidae.MsgNetworkErrorCallback(C.errorHandler), cname)

    n.net.Start()
    n.net.Listen(myAddr, &err); checkError(&err)
    n.net.Connect(otherAddr, false, &err); checkError(&err)
    return n
}

func main() {
    ec = salticidae.NewEventContext()
    err := salticidae.NewError()

    aliceAddr := salticidae.NewAddrFromIPPortString("127.0.0.1:12345", &err)
    bobAddr := salticidae.NewAddrFromIPPortString("127.0.0.1:12346", &err)

    alice = genMyNet(ec, "alice", aliceAddr, bobAddr)
    bob = genMyNet(ec, "bob", bobAddr, aliceAddr)

    ev_int := salticidae.NewSigEvent(ec, salticidae.SigEventCallback(C.onTerm), nil)
    ev_int.Add(salticidae.SIGINT)
    ev_term := salticidae.NewSigEvent(ec, salticidae.SigEventCallback(C.onTerm), nil)
    ev_term.Add(salticidae.SIGTERM)

    ec.Dispatch()
    alice.net.Stop()
    bob.net.Stop()
    C.free(unsafe.Pointer(alice.cname))
    C.free(unsafe.Pointer(bob.cname))
}
