package commands

import (
	"strconv"
)

type SetHandler struct{}

func (ch SetHandler) Authorize(IClientRequest) error {
	return nil
}

func (ch SetHandler) Permissions(IClientRequest) []string {
	return nil
}

func (ch SetHandler) Persistent(IClientRequest) bool {
	return true
}

func (ch SetHandler) HandlePipelined(r IClientRequest) []byte {
	return []byte{}
}
func (ch SetHandler) Handle(r IClientRequest) {
	value := string(r.Data()[6])
	if len(value) == 0 {
		_, _ = r.Write([]byte("-ERR value is empty\r\n"))
		return
	}

	key := string(r.Data()[4])
	if len(key) == 0 {
		_, _ = r.Write([]byte("-ERR key is empty\r\n"))
		return
	}

	exp := 0
	if len(r.Data()) > 8 {
		var err error
		exp, err = strconv.Atoi(string(r.Data()[10]))
		if err != nil {
			_, _ = r.Write([]byte("-ERR expiration cannot be casted to number\r\n"))
			return
		}
	}

	err := r.Store().Set(key, value, int64(exp))

	if err != nil {
		_, _ = r.Write([]byte("-ERR error\r\n"))
		return
	}

	_, _ = r.Write([]byte("+OK\r\n"))
}

func init() {
	RegisterCommand("set", SetHandler{})
}
