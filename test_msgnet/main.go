package main

// #cgo CFLAGS: -I${SRCDIR}/../salticidae/include/
// #include "salticidae/network.h"
// void onTerm(int sig);
// void onReceiveHello(msg_t *, msgnetwork_conn_t *);
// void onReceiveAck(msg_t *, msgnetwork_conn_t *);
// void connHandler(msgnetwork_conn_t *, bool);
import "C"

import (
    "encoding/binary"
    "fmt"
    "salticidae-go"
)

var ec salticidae.EventContext
const (
    MSG_OPCODE_HELLO salticidae.Opcode = iota
    MSG_OPCODE_ACK
)

//export onTerm
func onTerm(_ C.int) {
    ec.Stop()
}

type MsgHello struct {
    name string
    text string
}

func msgHelloSerialize(name string, text string) salticidae.Msg {
    serialized := salticidae.NewDataStream()
    t := make([]byte, 4)
    binary.LittleEndian.PutUint32(t, uint32(len(name)))
    serialized.PutData(t)
    serialized.PutData([]byte(name))
    serialized.PutData([]byte(text))
    return salticidae.NewMsg(MSG_OPCODE_HELLO, serialized.ToByteArray())
}

func msgHelloUnserialize(msg salticidae.Msg) MsgHello {
    p := msg.GetPayload()
    length := binary.LittleEndian.Uint32(p.GetData(4))
    name := string(p.GetData(int(length)))
    text := string(p.GetData(p.Size()))
    p.Free()
    return MsgHello { name: name, text: text }
}

func msgAckSerialize() salticidae.Msg {
    return salticidae.NewMsg(MSG_OPCODE_ACK, salticidae.NewByteArray())
}

type MyNet struct {
    net salticidae.MsgNetwork
    name string
}

var alice, bob MyNet

//export onReceiveHello
func onReceiveHello(__msg *C.struct_msg_t, _conn *C.struct_msgnetwork_conn_t) {
    _msg := salticidae.Msg(__msg)
    conn := salticidae.MsgNetworkConn(_conn)
    net := conn.GetNet()
    name := bob.name
    if net == alice.net.GetInner() {
        name = alice.name
    }
    msg := msgHelloUnserialize(_msg)
    fmt.Printf("[%s] %s says %s\n", name, msg.name, msg.text)
    ack := msgAckSerialize()
    net.SendMsg(ack, conn)
    ack.Free()
}

//export onReceiveAck
func onReceiveAck(_ *C.struct_msg_t, _conn *C.struct_msgnetwork_conn_t) {
    conn := salticidae.MsgNetworkConn(_conn)
    net := conn.GetNet()
    name := bob.name
    if net == alice.net.GetInner() {
        name = alice.name
    }
    fmt.Printf("[%s] the peer knows\n", name)
}

//export connHandler
func connHandler(_conn *C.struct_msgnetwork_conn_t, connected C.bool) {
    conn := salticidae.MsgNetworkConn(_conn)
    net := conn.GetNet()
    n := &bob
    if net == alice.net.GetInner() {
        n = &alice
    }
    name := n.name
    if connected {
        if conn.GetMode() == salticidae.CONN_MODE_ACTIVE {
            fmt.Printf("[%s] Connected, sending hello.", name)
            hello := msgHelloSerialize(name, "Hello there!")
            n.net.SendMsg(hello, conn)
            hello.Free()
        } else {
            fmt.Printf("[%s] Accepted, waiting for greetings.\n", name)
        }
    } else {
        fmt.Printf("[%s] Disconnected, retrying.\n", name)
        net.Connect(conn.GetAddr())
    }
}

func genMyNet(ec salticidae.EventContext, name string) MyNet {
    netconfig := salticidae.NewMsgNetworkConfig()
    n := MyNet {
        net: salticidae.NewMsgNetwork(ec, netconfig),
        name: name,
    }
    netconfig.Free()
    return n
}

func main() {
    ec = salticidae.NewEventContext()
    alice_addr := salticidae.NewAddrFromIPPortString("127.0.0.1:12345")
    bob_addr := salticidae.NewAddrFromIPPortString("127.0.0.1:12346")

    alice = genMyNet(ec, "Alice")
    bob = genMyNet(ec, "Bob")

    alice.net.RegHandler(MSG_OPCODE_HELLO, salticidae.MsgNetworkMsgCallback(C.onReceiveHello))
    alice.net.RegHandler(MSG_OPCODE_ACK, salticidae.MsgNetworkMsgCallback(C.onReceiveAck))
    alice.net.RegConnHandler(salticidae.MsgNetworkConnCallback(C.connHandler))

    bob.net.RegHandler(MSG_OPCODE_HELLO, salticidae.MsgNetworkMsgCallback(C.onReceiveHello))
    bob.net.RegHandler(MSG_OPCODE_ACK, salticidae.MsgNetworkMsgCallback(C.onReceiveAck))
    bob.net.RegConnHandler(salticidae.MsgNetworkConnCallback(C.connHandler))

    alice.net.Start()
    bob.net.Start()

    alice.net.Listen(alice_addr)
    bob.net.Listen(bob_addr)

    alice.net.Connect(bob_addr)
    bob.net.Connect(alice_addr)

    alice_addr.Free()
    bob_addr.Free()

    ev_int := salticidae.NewSigEvent(ec, salticidae.SigEventCallback(C.onTerm))
    ev_int.Add(salticidae.SIGINT)
    ev_term := salticidae.NewSigEvent(ec, salticidae.SigEventCallback(C.onTerm))
    ev_term.Add(salticidae.SIGTERM)

    ec.Dispatch()

    ev_int.Free()
    ev_term.Free()
    alice.net.Free()
    bob.net.Free()
    ec.Free()
}
