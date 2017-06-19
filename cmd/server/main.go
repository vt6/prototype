package main

import (
	"log"
	"net"
	"os"

	"github.com/vt6/prototype/pkg/server"
)

func main() {
	logger := log.New(os.Stderr, "vt6-server: ", log.LstdFlags)
	socket := &server.Socket{
		Handler: func(conn *net.UnixConn) {
			defer conn.Close()
			var buf [1024]byte
			n, err := conn.Read(buf[:])
			log.Printf("read %d bytes: %#v\n", n, string(buf[:n]))
			if err != nil {
				log.Printf("read error was: %s\n", err.Error())
			}
		},
		Logger: logger,
	}

	err := socket.Listen()
	if err != nil {
		logger.Fatal(err)
	}

	socket.Run()
}
