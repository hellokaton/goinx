package main

import (
	"io"
	"log"
	"net"
	"os"
)

type Server struct {
	Name      string   `yaml:"name"`
	Listen    string   `yaml:"listen"`
	Domains   []string `yaml:"domains"`
	Root      string   `yaml:"root"`
	ProxyPass string   `yaml:"proxy_pass"`
	KeyFile   string   `yaml:"key_file"`
	CertFile  string   `yaml:"cert_file"`
}

func (s *Server) Start() {
	log.Printf("[%s] listen %s proxy to %s", s.Name, s.Listen, s.ProxyPass)
	listener, err := net.Listen("tcp", s.Listen)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
			os.Exit(0)
		}

		go s.Serve(conn)
	}

}

func (s *Server) Serve(conn net.Conn) error {
	defer conn.Close()

	// get first 3 bytes of connection as header
	header := make([]byte, 3)
	if _, err := io.ReadAtLeast(conn, header, 3); err != nil {
		return err
	}

	// identify protocol from header
	address := s.ProxyPass
	log.Printf("[INFO] proxy: from=%s to=%s\n", conn.RemoteAddr(), address)

	// connect to remote
	remote, err := net.Dial("tcp", address)
	if err != nil {
		log.Printf("[ERROR] remote: %s\n", err)
		return err
	}
	defer remote.Close()

	// write header we chopped back to remote
	remote.Write(header)

	// proxy between us and remote server
	err = Shovel(conn, remote)
	if err != nil {
		return err
	}

	return nil
}

// proxy between two sockets
func Shovel(local, remote io.ReadWriteCloser) error {
	errch := make(chan error, 1)

	go chanCopy(errch, local, remote)
	go chanCopy(errch, remote, local)

	for i := 0; i < 2; i++ {
		if err := <-errch; err != nil {
			// If this returns early the second func will push into the
			// buffer, and the GC will clean up
			return err
		}
	}
	return nil
}

// copy between pipes, sending errors to channel
func chanCopy(e chan error, dst, src io.ReadWriter) {
	_, err := io.Copy(dst, src)
	e <- err
}
