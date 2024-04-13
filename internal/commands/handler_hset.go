package commands

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"strings"
)

type HSetHandler struct{}

func (ch HSetHandler) Authorize(ClientRequest) error {
	return nil
}

func (ch HSetHandler) Permissions() []string {
	return nil
}

func (ch HSetHandler) Persistent() bool {
	return true
}

func (ch HSetHandler) Handle(r ClientRequest) {
	r.Logger.Info().Str("hset key", string(r.Data[4])).Msg("hset handler")
	hs := chunkSlice(r.Data[5:], 4)

	data, err := hs.ToBytes()

	if err != nil {
		r.WriteError(err.Error())
		return
	}

	_, err = r.Store.HSet(string(r.Data[4]), data, 0)

	if err != nil {
		r.WriteError(err.Error())
		return
	}

	_ = hs.FromBytes(data)
	_, _ = r.Write([]byte(fmt.Sprintf(":%d\r\n", hs.Len())))
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

func (hs hset) Values() []string {
	keys := []string{}

	for _, i := range hs {
		keys = append(keys, i.Value())
	}

	return keys
}

func (hs hset) ToBytes() ([]byte, error) {
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

func init() {
	RegisterCommand("hset", HSetHandler{})
}
