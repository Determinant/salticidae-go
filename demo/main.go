package main

// void onTerm_cgo(int sig);
import "C"

import "salticidae-go"
import "unsafe"

var ec salticidae.EventContext

//export onTerm
func onTerm(_ int) {
    ec.Stop()
}

func run(ec salticidae.EventContext, my_addr string, other_addr string) salticidae.MsgNetwork {
    netconfig := salticidae.NewMsgNetworkConfig()
    net := salticidae.NewMsgNetwork(ec, netconfig)
    listen_addr := salticidae.NewAddrFromIPPortString(my_addr)
    connect_addr := salticidae.NewAddrFromIPPortString(other_addr)

    net.Start()
    net.Listen(listen_addr)
    net.Connect(connect_addr)
    return net
}

func main() {
    ec = salticidae.NewEventContext()
    ev_int := salticidae.NewSigEvent(ec, salticidae.SigEventCallback(unsafe.Pointer(C.onTerm_cgo)))
    ev_int.Add(salticidae.SIGINT)
    ev_term := salticidae.NewSigEvent(ec, salticidae.SigEventCallback(unsafe.Pointer(C.onTerm_cgo)))
    ev_term.Add(salticidae.SIGTERM)

    alice := "127.0.0.1:10000"
    bob := "127.0.0.1:10001"
    alice_net := run(ec, alice, bob)
    bob_net := run(ec, bob, alice)
    ec.Dispatch()
    ev_int.Free()
    ev_term.Free()
    alice_net.Free()
    bob_net.Free()
}
