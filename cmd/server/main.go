package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/vt6/prototype/pkg/server"
)

var logger = log.New(os.Stderr, "vt6-server: ", log.LstdFlags)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "usage: %s COMMAND [ARG]...\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  e.g. %s bash\n", os.Args[0])
		os.Exit(255)
	}

	//create the server socket
	socket := &server.Socket{
		Handler: handler,
		Logger:  logger,
	}
	err := socket.Listen()
	if err != nil {
		logger.Fatal(err)
	}
	go socket.Run()

	//shutdown socket on SIGINT or SIGTERM
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go shutdownOnSignal(c, socket)

	//execute command line
	os.Setenv("VT6", socket.Path())
	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	cmd.Stdin = os.Stdin   //TODO: wrap os.Stdin to catch VT100 input events, and convert them into VT6 input events if someone is listening on them
	cmd.Stdout = os.Stdout //TODO: wrap os.Stdout to catch VT100 escape sequences
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		logger.Fatal(err)
	}

	socket.Close()
}

func handler(conn *net.UnixConn) {
	defer conn.Close()
	var buf [1024]byte
	n, err := conn.Read(buf[:])
	logger.Printf("read %d bytes: %#v\n", n, string(buf[:n]))
	if err != nil {
		logger.Printf("read error was: %s\n", err.Error())
	}
}

func shutdownOnSignal(c <-chan os.Signal, socket *server.Socket) {
	if <-c != nil {
		socket.Close()
		os.Exit(0)
	}
}
