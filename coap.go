// Copyright (C) 2018 O.S. Systems Sofware LTDA
//
// SPDX-License-Identifier: MIT

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
		BackendURL: fmt.Sprintf("http://:%d/", httpPort),
		Timeout:    &timeout,
	}

	return p.Serve()
}
