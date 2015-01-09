package main

import (
	"net"

	"github.com/pyk/session"
)

func main() {

	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()
	log.Printf("SMTP server listening on %s", listener.Addr())

	for {
		// wait for transmission channel estabilished
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept: %v", err)
			continue
		}

		// define a session
		s := session.New(conn)
		// s.HandleHello(func() {})
		// s.HandleMail(func() {})
		// s.NewExtension()

		// handle every new connected session concurrently
		go s.Serve()
	}
}
