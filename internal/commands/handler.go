package commands

import (
	"net"
	"pedis/internal/storage"
)

type CommandHandler func([][]byte, storage.Storage, net.Conn)
