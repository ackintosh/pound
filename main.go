package main

import (
	"os"
	"path/filepath"

	"github.com/jessevdk/go-flags"
)

var version = "v0.0.1"
var opts Options

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	current, _ := filepath.Abs(".")
	maildrop := NewMaildrop(current + "/.pound")

	host := "localhost"
	var port int
	if opts.Port != 0 {
		port = opts.Port
	} else {
		port = 12345
	}
	s := NewServer(host, port, maildrop)

	err = s.Run()
	if err != nil {
		panic(err)
	}
}
