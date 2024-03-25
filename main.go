package main

import (
	"fmt"
	"log"
	"net"
	"pedis/internal/commands"
)

var data map[string]string

type RedisCommand interface {
	Run([]byte, net.Conn)
}

func main() {
	data = make(map[string]string)

	handlers := make(map[string]RedisCommand)
	handlers["*"] = commands.ArrayHandler{}

	listener, err := net.Listen("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on port 6379...")
	for {
		conn, err := listener.Accept()
		fmt.Println("received connection")

		if err != nil {
			fmt.Println("Error accepting connection:", err)
		}

		go func(conn net.Conn) {
			for {
				b := make([]byte, 1024)

				size, err := conn.Read(b)

				if err != nil || size == 0 {
					// write redis error message to the connection and close it
					conn.Close()
					continue
				}

				commandId := string(b[0])
				log.Println("command id", commandId)

				handler, commandNotFound := handlers[commandId]

				if !commandNotFound {
					log.Println(err)
					conn.Close()
					continue
				}

				handler.Run(b[1:size], conn)
			}
		}(conn)
	}
}
