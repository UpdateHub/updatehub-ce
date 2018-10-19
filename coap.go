package main

import (
	"fmt"
	"net"
	"time"

	"github.com/OSSystems/crosscoap"
)

func startCoapServer(coapPort, httpPort int) error {
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("0.0.0.0:%d", coapPort))
	if err != nil {
		return err
	}

	udpListener, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}
	defer udpListener.Close()

	timeout := time.Second * 10

	p := crosscoap.Proxy{
		Listener:   udpListener,
		BackendURL: fmt.Sprintf("http://127.0.0.1:%d/", httpPort),
		Timeout:    &timeout,
	}

	return p.Serve()
}
