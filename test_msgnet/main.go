package main

// #cgo CFLAGS: -I${SRCDIR}/../salticidae/include/
// #include "salticidae/network.h"
// void onTerm(int sig, void *);
// void onReceiveHello(msg_t *, msgnetwork_conn_t *, void *);
// void onReceiveAck(msg_t *, msgnetwork_conn_t *, void *);
// void connHandler(msgnetwork_conn_t *, bool, void *);
// void errorHandler(SalticidaeCError *, bool, void *);
import "C"

import (
    "encoding/binary"
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
    serialized := salticidae.NewDataStream()
    t := make([]byte, 4)
    binary.LittleEndian.PutUint32(t, uint32(len(name)))
    serialized.PutData(t)
    serialized.PutData([]byte(name))
    serialized.PutData([]byte(text))
    s := salticidae.NewByteArrayMovedFromDataStream(serialized)
    return salticidae.NewMsgMovedFromByteArray(
        MSG_OPCODE_HELLO, s)
}

func msgHelloUnserialize(msg salticidae.Msg) (name string, text string) {
    p := msg.ConsumePayload()
    length := binary.LittleEndian.Uint32(p.GetDataInPlace(4))
    name = string(p.GetDataInPlace(int(length)))
    text = string(p.GetDataInPlace(p.Size()))
    return
}

func msgAckSerialize() salticidae.Msg {
    s := salticidae.NewByteArray()
    return salticidae.NewMsgMovedFromByteArray(MSG_OPCODE_ACK, s)
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
func onReceiveHello(_msg *C.struct_msg_t, _conn *C.struct_msgnetwork_conn_t, _ unsafe.Pointer) {
    msg := salticidae.MsgFromC(salticidae.CMsg(_msg))
    conn := salticidae.MsgNetworkConnFromC(salticidae.CMsgNetworkConn(_conn))
    net := conn.GetNet()
    myName := bob.name
    if net == alice.net {
        myName = alice.name
    }
    name, text := msgHelloUnserialize(msg)
    fmt.Printf("[%s] %s says %s\n", myName, name, text)
    ack := msgAckSerialize()
    net.SendMsgByMove(ack, conn)
}

//export onReceiveAck
func onReceiveAck(_ *C.struct_msg_t, _conn *C.struct_msgnetwork_conn_t, _ unsafe.Pointer) {
    conn := salticidae.MsgNetworkConnFromC(salticidae.CMsgNetworkConn(_conn))
    net := conn.GetNet()
    name := bob.name
    if net == alice.net {
        name = alice.name
    }
    fmt.Printf("[%s] the peer knows\n", name)
}

//export connHandler
func connHandler(_conn *C.struct_msgnetwork_conn_t, connected C.bool, _ unsafe.Pointer) {
    conn := salticidae.MsgNetworkConnFromC(salticidae.CMsgNetworkConn(_conn))
    net := conn.GetNet()
    n := bob
    if net == alice.net {
        n = alice
    }
    name := n.name
    if connected {
        if conn.GetMode() == salticidae.CONN_MODE_ACTIVE {
            fmt.Printf("[%s] Connected, sending hello.\n", name)
            hello := msgHelloSerialize(name, "Hello there!")
            (*n.net).SendMsgByMove(hello, conn)
        } else {
            fmt.Printf("[%s] Accepted, waiting for greetings.\n", name)
        }
    } else {
        fmt.Printf("[%s] Disconnected, retrying.\n", name)
        err := salticidae.NewError()
        net.Connect(conn.GetAddr(), &err)
    }
}

//export errorHandler
func errorHandler(_err *C.struct_SalticidaeCError, fatal C.bool, _ unsafe.Pointer) {
    err := (*salticidae.Error)(unsafe.Pointer(_err))
    s := "recoverable"
    if fatal { s = "fatal" }
    fmt.Printf("Captured %s error during an async call: %s\n", s, salticidae.StrError(err.GetCode()))
}

func genMyNet(ec salticidae.EventContext, name string) MyNet {
    netconfig := salticidae.NewMsgNetworkConfig()
    n := MyNet {
        net: salticidae.NewMsgNetwork(ec, netconfig),
        name: name,
    }
    n.net.RegHandler(MSG_OPCODE_HELLO, salticidae.MsgNetworkMsgCallback(C.onReceiveHello), nil)
    n.net.RegHandler(MSG_OPCODE_ACK, salticidae.MsgNetworkMsgCallback(C.onReceiveAck), nil)
    n.net.RegConnHandler(salticidae.MsgNetworkConnCallback(C.connHandler), nil)
    n.net.RegErrorHandler(salticidae.MsgNetworkErrorCallback(C.errorHandler), nil)
    return n
}

func main() {
    ec = salticidae.NewEventContext()
    err := salticidae.NewError()

    alice_addr := salticidae.NewAddrFromIPPortString("127.0.0.1:12345")
    bob_addr := salticidae.NewAddrFromIPPortString("127.0.0.1:12346")

    alice = genMyNet(ec, "Alice")
    bob = genMyNet(ec, "Bob")

    alice.net.Start()
    bob.net.Start()

    alice.net.Listen(alice_addr, &err)
    checkError(&err)
    bob.net.Listen(bob_addr, &err)
    checkError(&err)

    alice.net.Connect(bob_addr, &err)
    checkError(&err)
    bob.net.Connect(alice_addr, &err)
    checkError(&err)

    ev_int := salticidae.NewSigEvent(ec, salticidae.SigEventCallback(C.onTerm), nil)
    ev_int.Add(salticidae.SIGINT)
    ev_term := salticidae.NewSigEvent(ec, salticidae.SigEventCallback(C.onTerm), nil)
    ev_term.Add(salticidae.SIGTERM)

    ec.Dispatch()
    alice.net.Stop()
    bob.net.Stop()
}
