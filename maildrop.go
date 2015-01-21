package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Maildrop struct {
	path string
}

var (
	ErrMessageNotExist = errors.New("message does not exist")
)

func NewMaildrop(dirname string) Maildrop {
	return Maildrop{path: dirname}
}

func (m Maildrop) MessageCount() int {
	count := 0
	err := filepath.Walk(
		m.path,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			count++
			return nil
		})
	if err != nil {
		fmt.Println(err)
	}
	return count
}

func (m Maildrop) Size() int {
	var size int64
	err := filepath.Walk(
		m.path,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			size += info.Size()
			return nil
		})
	if err != nil {
		fmt.Println(err)
	}

	return int(size / 8) // should be octets.
}

func (m Maildrop) Messages() []*Message {
	messages := make([]*Message, 0, m.MessageCount())
	seq := 1
	err := filepath.Walk(
		m.path,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			messages = append(messages, NewMessage(seq, path))
			seq += 1
			return nil
		})
	if err != nil {
		fmt.Println(err)
	}

	return messages
}

func (m Maildrop) MessageAt(seq int) (*Message, error) {
	for _, message := range m.Messages() {
		if seq == message.Seq() {
			return message, nil
		}
	}

	return nil, ErrMessageNotExist
}
