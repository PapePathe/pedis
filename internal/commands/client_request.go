package commands

import (
	"fmt"
	"net"
	"pedis/internal/renderer"
	"pedis/internal/storage"

	"github.com/rs/zerolog"
)

type IClientRequest interface {
	WriteError(string) error
	WriteString(string) error
	WriteNumber(string) error
	WriteArray([]string) error
	WriteOK() error
	WriteNil() error
	Write([]byte) (int, error)
}

type ClientRequest struct {
	Conn   net.Conn
	Data   [][]byte
	Store  storage.Storage
	Logger zerolog.Logger
}

func (c ClientRequest) WriteError(s string) error {
	str := fmt.Sprintf("-ERR %s\r\n", s)
	_, err := c.Conn.Write([]byte(str))
	if err != nil {
		return fmt.Errorf("net write error (%v)", err)
	}

	return nil
}

func (c ClientRequest) WriteArray(a []string) error {
	_, err := c.Conn.Write(renderer.RenderArray(a))
	if err != nil {
		return fmt.Errorf("net write error (%v)", err)
	}

	return nil
}

func (c ClientRequest) WriteString(s string) error {
	rdr := renderer.BulkStringRenderer{}
	_, err := c.Conn.Write(rdr.Render(s))
	if err != nil {
		return fmt.Errorf("net write error (%v)", err)
	}

	return nil
}

func (c ClientRequest) WriteNumber(s string) error {
	str := fmt.Sprintf(":%s\r\n", s)
	_, err := c.Conn.Write([]byte(str))
	if err != nil {
		return fmt.Errorf("net write error (%v)", err)
	}

	return nil
}

func (c ClientRequest) WriteOK() error {
	return nil
}

func (c ClientRequest) WriteNil() error {
	_, err := c.Conn.Write([]byte("$-1\r\n"))
	if err != nil {
		return fmt.Errorf("net write error (%v)", err)
	}
	return nil
}

func (c ClientRequest) Write([]byte) (int, error) {
	return 0, nil
}
