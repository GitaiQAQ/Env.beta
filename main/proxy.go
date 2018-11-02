// Copyright 2017 Gitai<i@gitai.me> All rights reserved.
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify,
// merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall
// be included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR
// ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF
// CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"strings"
)

type Proxy struct {
	address  net.Addr
	listener net.Listener
	dns ServeDNS
}

func (p *Proxy) init(address string) {
	if address == "" {
		address = "127.0.0.1:0"
	}

	l, err := net.Listen("tcp", address) // listen on localhost with random port
	if err != nil {
		log.Panic(err)
	}

	p.listener = l
	p.address = l.Addr()

	log.Println("Proxy run on address " + p.address.String())

	p.dns.load()

}

func (p *Proxy) loop() {

	for {
		client, err := p.listener.Accept()
		if err != nil {
			log.Panic(err)
		}

		go handleClientRequest(p.dns, client)
	}
}

func handleClientRequest(dns ServeDNS, client net.Conn) {
	if client == nil {
		return
	}
	defer client.Close()

	var b [1024]byte
	n, err := client.Read(b[:])
	if err != nil {
		log.Println(err)
		return
	}
	var method, host, address, port string
	fmt.Sscanf(string(b[:bytes.IndexByte(b[:], '\n')]), "%s%s", &method, &host)
	hostPortURL, err := url.Parse(host)
	if err != nil {
		log.Println(err)
		return
	}

	if hostPortURL.Opaque == "443" {
		address = hostPortURL.Scheme
		port = "443"
	} else { //http访问
		if strings.Index(hostPortURL.Host, ":") == -1 {
			address = hostPortURL.Host
			port = "80"
		} else {
			address = hostPortURL.Host
			port = hostPortURL.Port()
		}
	}

	dnsRes := dns.find(address)

	if address != dnsRes {
		log.Println("Host: " + address + " --> " + dnsRes)
	} else {
		log.Println("Host: " + address)
	}

	server, err := net.Dial("tcp", dnsRes + ":" + port)
	if err != nil {
		log.Println(err)
		return
	}
	if method == "CONNECT" {
		fmt.Fprint(client, "HTTP/1.1 200 Connection established\r\n\r\n")
	} else {
		server.Write(b[:n])
	}

	go io.Copy(server, client)
	io.Copy(client, server)
}

func testa() {
	var proxy = Proxy{}
	proxy.init("")
	proxy.loop()
}
