package server

import (
	"log"
	"net"
	"pedis/internal/storage"
)

type RedisCommand interface {
	Run([]byte, net.Conn, storage.Storage)
}

type RedisServer struct {
	handlers map[string]RedisCommand
}

func NewRedisServer() *RedisServer {
	return &RedisServer{
		handlers: make(map[string]RedisCommand),
	}
}

func (rs *RedisServer) Start() error {
	listener, err := net.Listen("tcp", "localhost:6379")
	if err != nil {
		return err
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
		}

		go rs.handleConnection(conn)
	}
}

func (rs *RedisServer) AddHandler(firstByte string, c RedisCommand) error {
	rs.handlers[firstByte] = c

	return nil
}

func (rs *RedisServer) handleConnection(conn net.Conn) {
	for {
		b := make([]byte, 1024)

		size, err := conn.Read(b)

		if err != nil || size == 0 {
			continue
		}

		commandId := string(b[0])
		handler, commandNotFound := rs.handlers[commandId]

		if !commandNotFound {
			log.Println(err)
			conn.Close()
			continue
		}

		handler.Run(b[1:size], conn, storage.NewSimpleStorage())
	}
}
