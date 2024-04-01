package commands

import (
	"net"
	"pedis/internal/renderer"
	"pedis/internal/storage"
)

func GetHandler(items [][]byte, store storage.Storage, conn net.Conn) {
	val, err := store.Get(string(items[4]))
	if err != nil {
		conn.Write([]byte("-ERR key not found\r\n"))
	}

	r := renderer.BulkStringRenderer{}
	conn.Write(r.Render(val))
}
