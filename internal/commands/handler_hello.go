package commands

import (
	"log"
	"net"
	"pedis/internal/response"
	"pedis/internal/storage"
)

func HelloHandler(_ [][]byte, _ storage.Storage, conn net.Conn) {
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
}
