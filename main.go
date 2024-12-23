package main

import (
	"embed"
	"flag"
	"fmt"
)

const (
	hakPrefix = "/.hak"
	ssePath   = "fs/sse"
	domain    = "backloop.dev"
)

var (
	dir  *string
	port *int

	//go:embed frontend/*
	frontend embed.FS
)

func init() {
	//	parse options and args
	dir = flag.String("dir", ".", "what directory to watch")
	port = flag.Int("port", 9443, "what port to listen on")
	flag.Parse()
}

func main() {

	args := flag.Args()
	subcommand := "serve"

	if len(args) > 0 {
		subcommand = args[0]
	}

	switch subcommand {
	case "", "serve":
		err := serve()
		barfOn(err)

	case "init":
		boilerplate()

	default:
		fmt.Println("unsupported subcommand")

	}

}
