package commands

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"pedis/internal/renderer"
	"pedis/internal/storage"
	"strings"

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

type RawRequest []byte

func (r RawRequest) String() string {
	return strings.ReplaceAll(string(r), "\\", "/")
}

func SliceAsChunks(slice [][]byte, chunkSize int) [][][]byte {
	var chunks [][][]byte

	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}
func (r RawRequest) ReadArray() []string {
	items := bytes.Split(r[2:], []byte{13, 10})
	sl := SliceAsChunks(items[3:], 2)
	array := []string{}

	log.Println(sl)

	for _, i := range sl {
		if len(i) == 2 {
			array = append(array, string(i[1]))
		}
	}

	return array
}

type ClientRequest struct {
	Conn    net.Conn
	Data    [][]byte
	DataRaw RawRequest
	Store   storage.Storage
	Logger  zerolog.Logger
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
	_, err := c.Conn.Write([]byte("+OK\r\n"))
	if err != nil {
		return fmt.Errorf("net write error (%v)", err)
	}
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
