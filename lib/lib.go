package lib

import (
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"time"
)

// Translates a byte count into a human readable string
func PrintRateHuman(count int) string {
	prefixes := []string{
		"B",
		"KB",
		"MB",
		"GB",
		"TB",
		"PB"}
	prefix := 0
	for {
		if count/1024 == 0 {
			break
		}
		prefix += 1
		count = count / 1024
	}

	return fmt.Sprint(count, prefixes[prefix], " ", 8*count, prefixes[prefix], "its")

}

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

		b := make([]byte, 10000)
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

	count := 0
	b := make([]byte, 10000)
	f := 1000
	last := time.Now()
	for i := 0; ; i++ {

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
		count += n
		if i%1000 == 0 {
			test := time.Now()
			if test.Sub(last) > time.Second {
				log.Print("Bounced ", PrintRateHuman(count))
				count = 0
				last = test
			}
		}
	}

}

func DoIt(listen bool, send string) {
	go http.ListenAndServe(":8080", nil)
	if listen {
		DoListen(send)
	} else {
		DoSend(send)
	}
}
