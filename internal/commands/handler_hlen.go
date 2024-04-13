package commands

import (
	"fmt"
)

type HLenHandler struct{}

func (ch HLenHandler) Authorize(IClientRequest) error {
	return nil
}

func (ch HLenHandler) Permissions(IClientRequest) []string {
	return nil
}

func (ch HLenHandler) Persistent(IClientRequest) bool {
	return false
}

func (ch HLenHandler) Handle(r IClientRequest) {
	data, err := r.Store().HGet(string(r.Data()[4]))

	if err != nil {
		_ = r.WriteString(fmt.Sprintf("%d", 0))
		return
	}

	hs := hset{}
	hs.FromBytes(data)

	_ = r.WriteNumber(fmt.Sprintf("%d", hs.Len()))
}

func init() {
	RegisterCommand("hlen", HLenHandler{})
}
