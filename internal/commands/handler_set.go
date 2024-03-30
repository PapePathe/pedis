package commands

import (
	"log"
	"net"
	"pedis/internal/storage"
)

func SetHandler(items [][]byte, store storage.Storage, conn net.Conn) {
	value := string(items[6])
	if len(value) == 0 {
		_, _ = conn.Write([]byte("-ERR value is empty\r\n"))
		return
	}

	key := string(items[4])

	if len(key) == 0 {
		_, _ = conn.Write([]byte("-ERR key is empty\r\n"))
		return
	}

	log.Println("going to execute set command")
	err := store.Set(key, value)

	if err != nil {
		_, _ = conn.Write([]byte("-ERR error\r\n"))
		return
	}

	_, _ = conn.Write([]byte("+OK\r\n"))
}
