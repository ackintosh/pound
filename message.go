package main

import (
	"bufio"
	"fmt"
	"os"
)

type Message struct {
	path string
	seq  int
}

func NewMessage(seq int, path string) *Message {
	return &Message{seq: seq, path: path}
}

func (m Message) Path() string {
	return m.path
}

func (m Message) Seq() int {
	return m.seq
}

func (m Message) Size() int {
	var size int64
	fileinfo, err := os.Stat(m.Path())
	if err != nil {
		fmt.Println(err)
	}
	size = fileinfo.Size()

	return int(size / 8) // should be octets.
}

func (m Message) Body() string {
	fp, err := os.Open(m.path)
	if err != nil {
		fmt.Println(err)
	}
	defer fp.Close()

	var body string
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		body += scanner.Text() + "\n"
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	return body
}
