package main

import "salticidae-go"

func run(my_addr string, other_addr string) {
    netconfig := salticidae.NewMsgNetworkConfig()
    ec := salticidae.NewEventContext()
    net := salticidae.NewMsgNetwork(ec, netconfig)
    listen_addr := salticidae.NewAddrFromIPPortString(my_addr)
    connect_addr := salticidae.NewAddrFromIPPortString(other_addr)

    net.Start()
    net.Listen(listen_addr)
    net.Connect(connect_addr)
    ec.Dispatch()
}

func main() {
    alice := "127.0.0.1:10000"
    bob := "127.0.0.1:10001"
    go run(alice, bob)
    run(bob, alice)
}
