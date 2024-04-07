package commands

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net"
	"pedis/internal/storage"
	"strings"
)

// HSetHandler
func HSetHandler(items [][]byte, store storage.Storage, conn net.Conn) {
	hs := chunkSlice(items[3:], 4)

	data, err := hs.ToBytes()

	if err != nil {
		_, _ = conn.Write([]byte("-ERR future error message\r\n"))
		return
	}

	log.Println("hset key", string(items[2]))
	_, err = store.HSet(string(items[2]), data, 0)

	if err != nil {
		_, _ = conn.Write([]byte("-ERR future error message\r\n"))
		return
	}

	err = hs.FromBytes(data)
	log.Println(err)
	_, _ = conn.Write([]byte(fmt.Sprintf(":%d\r\n", hs.Len())))
}

type hasharray [][]byte

func (ha hasharray) Key() string {
	return string(ha[1])
}

func (ha hasharray) Value() string {
	if len(ha) < 4 {
		return ""
	}
	return string(ha[3])
}

func (ha hasharray) String() string {
	sb := strings.Builder{}

	sb.WriteString("hasharray[")
	for idx, item := range ha {
		str := fmt.Sprintf("(i=%d v=%s),", idx, string(item))
		sb.WriteString(str)
	}
	sb.WriteString("]")

	return sb.String()
}

type hset []hasharray

func (hs hset) Len() int {
	return len(hs)
}

func (hs hset) Get(k string) (string, error) {
	for _, i := range hs {
		if i.Key() == k {
			return i.Value(), nil
		}
	}

	return "", errors.New("key not found in set")
}

func (hs hset) Keys() []string {
	keys := []string{}

	for _, i := range hs {
		keys = append(keys, i.Key())
	}

	return keys
}

func (hs hset) ToBytes() ([]byte, error) {
	log.Println(hs)
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(hs); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (hs *hset) FromBytes(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	if err := dec.Decode(&hs); err != nil {
		return err
	}

	return nil
}

func chunkSlice(slice [][]byte, chunkSize int) hset {
	var chunks hset

	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}
