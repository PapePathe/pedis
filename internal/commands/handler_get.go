package commands

import (
	"log"
	"net"
	"pedis/internal/renderer"
	"pedis/internal/storage"
)

func GetHandler(items [][]byte, store storage.Storage, conn net.Conn) {
	log.Println("going to execute get command")
	log.Println(string(items[0]), string(items[1]), string(items[2]), string(items[3]), string(items[4]))
	val, err := store.Get(string(items[4]))
	if err != nil {
		conn.Write([]byte("-ERR  key not found\r\n"))
	}
	log.Println("value from store get", val)
	r := renderer.BulkStringRenderer{}
	conn.Write(r.Render(val))
}
