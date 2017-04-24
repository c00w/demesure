package lib

import (
	"log"
	"net"
)

func DoListen(addr string) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("Error setting up tcp connection: ", err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print("Error Accepting connection: ", err)
			continue
		}

		b := make([]byte, 1000)
		f := 0
		for {
			n, err := conn.Read(b[f:])

			if err != nil {
				log.Print("Error Reading Data: ", err)
				break
			}

			f += n

			n, err = conn.Write(b[:f])
			if err != nil {
				log.Print("Error Writing Data: ", err)
				break
			}
			f -= n

		}
	}

}

func DoSend(send string) {
	conn, err := net.Dial("tcp", send)
	if err != nil {
		log.Fatal("Error Connecting: ", err)
	}
	b := make([]byte, 1000)
	f := 1000
	for {

		n, err := conn.Write(b[:f])
		if err != nil {
			log.Print("Error Writing Data: ", err)
			break
		}
		f -= n
		n, err = conn.Read(b[f:])

		if err != nil {
			log.Print("Error Reading Data: ", err)
			break
		}

		f += n
	}

}

func DoIt(listen bool, send string) {
	if listen {
		DoListen(send)
	} else {
		DoSend(send)
	}
}
