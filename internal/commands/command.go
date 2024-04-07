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

func DefaultRequestHandler() *RequestHandler {
	subcommands := map[string]CommandHandler{
		"hexists": HExistsHandler,
		"hget":    HGetHandler,
		"hkeys":   HKeysHandler,
		"hlen":    HLenHandler,
		//"hset":    HSetHandler,
		"hvals": HValsHandler,
	}
	return &RequestHandler{subcommands}
}

func (s RequestHandler) Run(data []byte, conn net.Conn, store storage.Storage) {
	items := bytes.Split(data, []byte{13, 10})

	subcommand := string(items[2])

	request := ClientRequest{
		Conn:  conn,
		Data:  items,
		store: store,
	}

	if h, ok := s.subcommands[subcommand]; ok {
		go h(request)
	} else {
		switch string(items[2]) {
		case "hello":
			go HelloHandler(items, store, conn)
		case "get":
			go GetHandler(items, store, conn)
		case "set":
			go SetHandler(items, store, conn)
		case "hset":
			go HSetHandler(items[2:], store, conn)
		case "client":
			log.Println("going to execute client options command")
			request.WriteString("OK")
		default:
			log.Println("command", string(items[2]), "is not yet supported")
			request.WriteError("command not supported")
		}
	}
}
