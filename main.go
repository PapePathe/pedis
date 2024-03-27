package main

import (
	"log"
	"pedis/internal/commands"
	"pedis/internal/server"
)

func main() {
	rs := server.NewRedisServer()
	rs.AddHandler("*", commands.ObjectHandler{})

	if err := rs.Start(); err != nil {
		log.Fatal(err)
	}
}
