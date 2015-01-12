package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/pyk/session"
)

type SMTPserver struct {
	Listener *net.TCPListener
	Stoped   chan bool
	Wg       *sync.WaitGroup
}

func (tcps *SMTPserver) Run() {
	defer tcps.Wg.Done()
	for {
		select {
		case <-tcps.Stoped:
			log.Println("smtpserver: stopping listening on", tcps.Listener.Addr())
			tcps.Listener.Close()
			return
		default:
		}

		// make sure listener.AcceptTCP() doesn't block forever
		// so it can read a stopped channel
		tcps.Listener.SetDeadline(time.Now().Add(1e9))
		conn, err := tcps.Listener.AcceptTCP()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			log.Println(err)
		}

		tcps.Wg.Add(1)
		s := session.New(conn, tcps.Wg, tcps.Stoped)
		go s.Serve()
	}
}

func (tcps *SMTPserver) Stop() {
	close(tcps.Stoped)
	tcps.Wg.Wait()
}

func main() {

	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("smtpserver: listening on %s", listener.Addr())

	server := &SMTPserver{
		Listener: listener,
		Stoped:   make(chan bool),
		Wg:       &sync.WaitGroup{},
	}

	server.Wg.Add(1)
	go server.Run()

	chs := make(chan os.Signal)
	signal.Notify(chs, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-chs)

	server.Stop()
}
