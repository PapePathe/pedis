package commands

import (
	"bytes"
	"fmt"
	"net"
	"pedis/internal/renderer"
	"pedis/internal/storage"
	"strings"

	"go.etcd.io/etcd/raft/v3/raftpb"
)

type IClientRequest interface {
	WriteError(string) error
	WriteString(string) error
	WriteNumber(string) error
	WriteArray([]string) error
	WriteOK() error
	WriteNil() error
	Write([]byte) (int, error)
	Data() [][]byte
	DataRaw() RawRequest
	Store() storage.Storage
	SendClusterConfigChange(raftpb.ConfChange)
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

	for _, i := range sl {
		if len(i) == 2 {
			array = append(array, string(i[1]))
		}
	}

	return array
}

type ClientRequest struct {
	Conn               net.Conn
	data               [][]byte
	dataRaw            RawRequest
	store              storage.Storage
	clusterChangesChan chan<- raftpb.ConfChange
}

func NewClientRequest(c net.Conn, d [][]byte, s storage.Storage, rd RawRequest, cchan chan<- raftpb.ConfChange) IClientRequest {
	return ClientRequest{
		Conn:               c,
		data:               d,
		store:              s,
		dataRaw:            rd,
		clusterChangesChan: cchan,
	}

}

func (c ClientRequest) SendClusterConfigChange(cc raftpb.ConfChange) {
	c.clusterChangesChan <- cc
}

func (c ClientRequest) Data() [][]byte {
	return c.data
}

func (c ClientRequest) DataRaw() RawRequest {
	return c.dataRaw
}

func (c ClientRequest) Store() storage.Storage {
	return c.store
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

func (c ClientRequest) Write(data []byte) (int, error) {
	n, err := c.Conn.Write(data)
	if err != nil {
		return 0, fmt.Errorf("net write error (%v)", err)
	}

	return n, nil
}
