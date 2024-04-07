package commands

import (
	"net"
	"pedis/internal/storage"
)

type CommandHandler func([][]byte, storage.Storage, net.Conn)

type CommandHandlerV2 func(ClientRequest)
