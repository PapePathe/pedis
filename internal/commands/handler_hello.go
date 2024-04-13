package commands

import (
	"log"
	"pedis/internal/response"
)

type HelloHandler struct{}

func (ch HelloHandler) Authorize(ClientRequest) error {
	return nil
}

func (ch HelloHandler) Permissions() []string {
	return nil
}

func (ch HelloHandler) Persistent() bool {
	return false
}

func (ch HelloHandler) Handle(r ClientRequest) {
	hr := response.HelloResponse{
		Server:  "redis",
		Version: "6.2.1",
		Mode:    "standalone",
		Proto:   3,
		Role:    "master",
	}

	_, err := r.Write(hr.Render())

	if err != nil {
		log.Println(err)
	}
}
