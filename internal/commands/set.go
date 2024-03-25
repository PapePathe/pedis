package commands

import (
	"bytes"
	"log"
	"net"
)

type ArrayHandler struct{}

type HelloResponse struct {
	Server  string
	Version string
	Mode    string
	Proto   uint8
	Id      uint8
	Role    string
}

func (s ArrayHandler) Run(data []byte, conn net.Conn) {
	//	log.Println("raw data", data)
	//	log.Println("string data", string(data))
	log.Println("received request, payload size=", len(data), string(data[0]))

	items := bytes.Split(data, []byte{13, 10})

	log.Println("sub command", string(items[2]))
	switch string(items[2]) {
	case "hello":
		log.Println("Respond to hello command")
		buf := bytes.Buffer{}

		buf.Write([]byte("%"))
		buf.Write([]byte("6"))
		buf.Write([]byte{13, 10})
		buf.Write([]byte("+server"))
		buf.Write([]byte{13, 10})
		buf.Write([]byte("+redis"))
		buf.Write([]byte{13, 10})
		buf.Write([]byte("+version"))
		buf.Write([]byte{13, 10})
		buf.Write([]byte("+6.2.11"))
		buf.Write([]byte{13, 10})
		buf.Write([]byte("+proto"))
		buf.Write([]byte{13, 10})
		buf.Write([]byte(":3"))
		buf.Write([]byte{13, 10})
		buf.Write([]byte("+id"))
		buf.Write([]byte{13, 10})
		buf.Write([]byte(":1"))
		buf.Write([]byte{13, 10})
		buf.Write([]byte("+mode"))
		buf.Write([]byte{13, 10})
		buf.Write([]byte("+standalone"))
		buf.Write([]byte{13, 10})
		buf.Write([]byte("+role"))
		buf.Write([]byte{13, 10})
		buf.Write([]byte("+master"))
		buf.Write([]byte{13, 10})

		n, err := conn.Write(buf.Bytes())

		if err != nil {
			log.Println(err)
		}

		log.Println("Wrote", n, "bytes")
	case "set":
		log.Println("going to execute set command")
	case "client":
		log.Println("going to execute client options command")
	default:
		log.Println(string(items[2]))
	}
}
