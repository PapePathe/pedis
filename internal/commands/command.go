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

	request := ClientRequest{
		Conn:  conn,
		Data:  items,
		store: store,
	}
	switch string(items[2]) {
	case "hello":
		HelloHandler(items, store, conn)
	case "get":
		GetHandler(items, store, conn)
	case "set":
		SetHandler(items, store, conn)
	case "hset":
		go HSetHandler(items[2:], store, conn)
	case "hget":
		go HGetHandler(request)
	case "hlen":
		go HLenHandler(request)
	case "hkeys":
		go HKeysHandler(request)
	case "hvals":
		go HValsHandler(request)
	case "client":
		log.Println("going to execute client options command")
	default:
		log.Println("command", string(items[2]), "is not yet supported")
	}
}
