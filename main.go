package main

import "path/filepath"

var version = "v0.0.1"

func main() {
	current, _ := filepath.Abs(".")
	maildrop := NewMaildrop(current + "/.pound")

	host := "localhost"
	port := 12345
	s := NewServer(host, port, maildrop)

	err := s.Run()
	if err != nil {
		panic(err)
	}
}
