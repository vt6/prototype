package server

import (
	"errors"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

//Socket creates and handles the socket that a VT6 server publishes when in
//normal mode.
type Socket struct {
	//The Handler function is called once (in a separate goroutine) for each
	//accepted connection. The Handler should defer a call to Close() on the
	//connection it is given, to make sure that it is properly cleaned up.
	Handler    func(*net.UnixConn)
	Logger     *log.Logger
	socketPath string
	listener   *net.UnixListener
}

//Listen starts listening on the VT6 server socket.
func (s *Socket) Listen() (err error) {
	s.socketPath, err = chooseSocketPath()
	if err != nil {
		return errors.New("cannot create VT6 server socket: " + err.Error())
	}

	s.listener, err = net.ListenUnix("unix",
		&net.UnixAddr{s.socketPath, "unix"},
	)
	if err != nil {
		return errors.New("cannot listen on " + s.socketPath + ": " + err.Error())
	}

	return nil
}

//Run accepts connections on this socket forever, then cleans up the socket
//before exiting.
func (s *Socket) Run() error {
	defer s.listener.Close()
	defer os.Remove(s.socketPath)

	for {
		conn, err := s.listener.AcceptUnix()
		if err != nil {
			s.Logger.Printf("cannot accept: %s\n", err.Error())
			continue
		}

		go s.Handler(conn)
	}
}

func chooseSocketPath() (string, error) {
	pid := strconv.Itoa(os.Getpid())

	//prefer $XDG_RUNTIME_DIR which is usually on a tmpfs and already chmod'd so
	//that only the current user can access it
	//FIXME: This could be a problem when we use sudo(1), because the user
	//account under which we execute the target program could not have access to
	//the VT6 socket anymore. Maybe we need multiplexed mode there as well?
	str := os.Getenv("XDG_RUNTIME_DIR")
	if str != "" {
		dir := filepath.Join(str, "vt6")
		path := filepath.Join(dir, pid)
		return path, os.MkdirAll(dir, 0700)
	}

	//put a socket into the tempdir (without extra directory to ensure that
	//everything is cleaned up when this process exits)
	return filepath.Join(os.TempDir(), pid), nil
}
