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
    return salticidae.NewMsgMovedFromByteArray(
        MSG_OPCODE_HELLO, salticidae.NewByteArrayMovedFromDataStream(serialized))
}

func msgHelloUnserialize(msg salticidae.Msg) (name string, text string) {
    p := msg.GetPayloadByMove()
    t := p.GetDataInPlace(4); length := binary.LittleEndian.Uint32(t.Get()); t.Release()
    t = p.GetDataInPlace(int(length)); name = string(t.Get()); t.Release()
    t = p.GetDataInPlace(p.Size()); text = string(t.Get()); t.Release()
    return
}

func msgAckSerialize() salticidae.Msg {
    return salticidae.NewMsgMovedFromByteArray(MSG_OPCODE_ACK, salticidae.NewByteArray())
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
    peerCert salticidae.UInt256
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
    res := true
    if connected {
        certHash := conn.GetPeerCert().GetDer().GetHash()
        res = certHash.IsEq(n.peerCert)
        if conn.GetMode() == salticidae.CONN_MODE_ACTIVE {
            fmt.Printf("[%s] Connected, sending hello.\n", myName)
            hello := msgHelloSerialize(myName, "Hello there!")
            n.net.SendMsg(hello, conn)
        } else {
            status := "fail"
            if res { status = "ok" }
            fmt.Printf("[%s] Accepted, waiting for greetings.\n" +
                        "The peer certificate footprint is %s (%s)\n",
                        myName, certHash.GetHex(), status)
        }
    } else {
        fmt.Printf("[%s] Disconnected, retrying.\n", myName)
        err := salticidae.NewError()
        net.Connect(conn.GetAddr(), false, &err)
    }
    return C.bool(res)
}

//export errorHandler
func errorHandler(_err *C.struct_SalticidaeCError, fatal C.bool, _ unsafe.Pointer) {
    err := (*salticidae.Error)(unsafe.Pointer(_err))
    s := "recoverable"
    if fatal { s = "fatal" }
    fmt.Printf("Captured %s error during an async call: %s\n", s, salticidae.StrError(err.GetCode()))
}

func genMyNet(ec salticidae.EventContext,
            name string, peerCert string,
            myAddr salticidae.NetAddr, otherAddr salticidae.NetAddr) MyNet {
    err := salticidae.NewError()
    netconfig := salticidae.NewMsgNetworkConfig()
    netconfig.EnableTLS(true)
    netconfig.TLSKeyFile(name + ".pem")
    netconfig.TLSCertFile(name + ".pem")
    net := salticidae.NewMsgNetwork(ec, netconfig, &err); checkError(&err)
    _peerCert := salticidae.NewUInt256FromByteArray(salticidae.NewByteArrayFromHex(peerCert))
    cname := C.CString(name)
    n := MyNet { net: net, name: name, peerCert: _peerCert, cname: cname }
    n.net.RegHandler(MSG_OPCODE_HELLO, salticidae.MsgNetworkMsgCallback(C.onReceiveHello), unsafe.Pointer(cname))
    n.net.RegHandler(MSG_OPCODE_ACK, salticidae.MsgNetworkMsgCallback(C.onReceiveAck), unsafe.Pointer(cname))
    n.net.RegConnHandler(salticidae.MsgNetworkConnCallback(C.connHandler), unsafe.Pointer(cname))
    n.net.RegErrorHandler(salticidae.MsgNetworkErrorCallback(C.errorHandler), unsafe.Pointer(cname))

    n.net.Start()
    n.net.Listen(myAddr, &err); checkError(&err)
    n.net.Connect(otherAddr, false, &err); checkError(&err)
    return n
}

func main() {
    ec = salticidae.NewEventContext()

    aliceAddr := salticidae.NewAddrFromIPPortString("127.0.0.1:12345")
    bobAddr := salticidae.NewAddrFromIPPortString("127.0.0.1:12346")

    alice = genMyNet(ec, "alice", "ed5a9a8c7429dcb235a88244bc69d43d16b35008ce49736b27aaa3042a674043", aliceAddr, bobAddr)
    bob = genMyNet(ec, "bob", "ef3bea4e72f4d0e85da7643545312e2ff6dded5e176560bdffb1e53b1cef4896", bobAddr, aliceAddr)

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
