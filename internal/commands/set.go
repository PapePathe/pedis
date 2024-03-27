package commands

import (
	"bytes"
	"log"
	"net"

	"pedis/internal/response"
	"pedis/internal/storage"
)

type ObjectHandler struct{}

func (s ObjectHandler) Run(data []byte, conn net.Conn, store storage.Storage) {
	//	log.Println("raw data", data)
	//	log.Println("string data", string(data))
	log.Println("received request, payload size=", len(data), string(data[0]))

	items := bytes.Split(data, []byte{13, 10})

	log.Println("sub command", string(items[2]))
	switch string(items[2]) {
	case "hello":
		log.Println("Respond to hello command")
		hr := response.HelloResponse{
			Server:  "redis",
			Version: "6.2.1",
			Mode:    "standalone",
			Proto:   3,
			Role:    "master",
		}

		_, err := conn.Write(hr.Render())

		if err != nil {
			log.Println(err)
		}
	case "set":
		log.Println("going to execute set command")
		log.Println(string(items[0]), string(items[1]), string(items[2]), string(items[3]), string(items[4]), string(items[5]), string(items[6]))
		err := store.Set(string(items[4]), string(items[6]))

		if err != nil {
			conn.Write([]byte("-ERR error\r\n"))
		}

		conn.Write([]byte("+OK\r\n"))
	case "client":
		log.Println("going to execute client options command")
	default:
		log.Println(string(items[2]))
	}
}
