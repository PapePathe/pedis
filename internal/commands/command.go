package commands

import (
	"bytes"
	"log"
	"net"

	"pedis/internal/storage"
)

type RequestHandler struct {
	subcommands map[string]CommandHandler
}

func NewRequestHandler() *RequestHandler {
	return &RequestHandler{}
}

func (s RequestHandler) Run(data []byte, conn net.Conn, store storage.Storage) {
	items := bytes.Split(data, []byte{13, 10})

	log.Println(len(items))

	log.Println("sub command", string(items[2]))
	switch string(items[2]) {
	case "hello":
		HelloHandler(items, store, conn)
	case "get":
		GetHandler(items, store, conn)
	case "set":
		SetHandler(items, store, conn)
	case "client":
		log.Println("going to execute client options command")
	default:
		log.Println("command", string(items[2]), "is not yet supported")
	}
}