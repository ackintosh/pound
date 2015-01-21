package main

import (
	"fmt"
	"net"
	"strconv"
)

type Command struct {
}

func (cmd *Command) User(client net.Conn) {
	// nop.
	client.Write([]byte("+OK\n"))
}

func (cmd *Command) Pass(client net.Conn) {
	// nop.
	client.Write([]byte("+OK\n"))
}

func (cmd *Command) Dele(client net.Conn) {
	// nop.
	client.Write([]byte("+OK\n"))
}

func (cmd *Command) Stat(client net.Conn, maildrop Maildrop) {
	reply := fmt.Sprintf(
		"+OK %d %d",
		maildrop.MessageCount(), maildrop.Size())
	client.Write([]byte(reply + "\n"))
}

func (cmd *Command) List(client net.Conn, maildrop Maildrop, args string) {
	if args == "" {
		reply := fmt.Sprintf(
			"+OK %d messages (%d octets)",
			maildrop.MessageCount(), maildrop.Size()) + "\n"

		for _, m := range maildrop.Messages() {
			reply += fmt.Sprintf("%d %d\n", m.Seq(), m.Size())
		}
		reply += ".\n"
		client.Write([]byte(reply))

		return
	}

	seq, err := strconv.Atoi(args)
	if err != nil {
		fmt.Println(err)
	}

	message, err := maildrop.MessageAt(seq)
	if err != nil {
		fmt.Println(err)
		return
	}
	reply := fmt.Sprintf("+OK %d %d\n", seq, message.Size())
	client.Write([]byte(reply))
}

func (cmd *Command) Retr(client net.Conn, maildrop Maildrop, args string) {
	if args == "" {
		client.Write([]byte("-ERR"))
		return
	}

	seq, err := strconv.Atoi(args)
	if err != nil {
		client.Write([]byte("-ERR"))
		return
	}

	message, err := maildrop.MessageAt(seq)
	if err != nil {
		client.Write([]byte("-ERR"))
		return
	}

	reply := fmt.Sprintf("+OK %d octets\n", message.Size())
	reply += message.Body()
	reply += ".\n"
	client.Write([]byte(reply))
}

func (cmd *Command) Quit(client net.Conn) {
	client.Write([]byte("+OK\n"))
	client.Close()
}
