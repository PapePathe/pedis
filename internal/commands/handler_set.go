package commands

import (
	"net"
	"pedis/internal/storage"
	"strconv"
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

	exp, err := strconv.Atoi(string(items[10]))
	if err != nil {
		_, _ = conn.Write([]byte("-ERR expiration cannot be casted to number\r\n"))
		return
	}

	err = store.Set(key, value, int64(exp))

	if err != nil {
		_, _ = conn.Write([]byte("-ERR error\r\n"))
		return
	}

	_, _ = conn.Write([]byte("+OK\r\n"))
}
