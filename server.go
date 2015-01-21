package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	Host     string
	Port     int
	Maildrop Maildrop
	Conns    chan net.Conn
	Fin      chan bool
}

func NewServer(host string, port int, maildrop Maildrop) Server {
	return Server{Host: host, Port: port, Maildrop: maildrop}
}

func (s *Server) Run() error {
	listener, err := net.Listen("tcp", s.Host+":"+strconv.Itoa(s.Port))
	if err != nil {
		return err
	}

	s.Conns = make(chan net.Conn)
	s.Fin = make(chan bool)

	go s.ClientConns(listener)
loop:
	for {
		select {
		case conns := <-s.Conns:
			go s.HandleConn(conns)
		case <-s.Fin:
			fmt.Print("finished\n")
			break loop
		}
	}
	fmt.Print("Run goroutine finished\n")

	return nil
}

func (s *Server) ClientConns(listener net.Listener) {
	i := 0
	for {
		client, err := listener.Accept()
		if err != nil {
			fmt.Printf("couldn't accept: " + err.Error())
			continue
		}
		i++
		fmt.Printf("%d: %v <-> %v\n", i, client.LocalAddr(), client.RemoteAddr())
		s.Conns <- client
	}

	return
}

func (s *Server) HandleConn(client net.Conn) {
	command := Command{}
loop:
	for {
		cmd, args, err := s.readClientCommand(client)
		if err != nil {
			// When the timer expires,
			// the server should close the TCP connection
			// without sending any response to the client.
			fmt.Println(err)
			client.Close()
			break loop
		}

		switch cmd {
		case "USER":
			command.User(client)
		case "PASS":
			command.Pass(client)
		case "DELE":
			command.Dele(client)
		case "STAT":
			command.Stat(client, s.Maildrop)
		case "LIST":
			command.List(client, s.Maildrop, args)
		case "RETR":
			command.Retr(client, s.Maildrop, args)
		case "QUIT":
			command.Quit(client)
			break loop
		default:
			client.Write([]byte("-ERR unknown command\n"))
		}
	}

	return
}

func (s *Server) Shutdown() {
	fmt.Print("Shutdown called.\n")
	s.Fin <- true

	return
}

func (s *Server) readClientCommand(client net.Conn) (cmd string, args string, err error) {
	b := bufio.NewReader(client)

	// a timer MUST be of at least 10 minutes duration.
	client.SetDeadline(time.Now().Add(10 * time.Minute))
	line, err := b.ReadString('\n')
	// reset the autologout timer.
	client.SetDeadline(time.Time{})

	if err == nil {
		line = strings.TrimRight(line, "\r\n")
		parts := strings.SplitN(line, " ", 2)
		cmd = strings.ToUpper(parts[0])
		if len(parts) == 2 {
			args = parts[1]
		}

	} else {
		if e, ok := err.(net.Error); ok && e.Timeout() {
			fmt.Println("connection timed out.\n")
			err = e
		}
	}

	return cmd, args, err
}
