package commands

import (
	"fmt"
)

type HLenHandler struct{}

func (ch HLenHandler) Authorize(ClientRequest) error {
	return nil
}

func (ch HLenHandler) Permissions() []string {
	return nil
}

func (ch HLenHandler) Persistent() bool {
	return false
}

func (ch HLenHandler) Handle(r ClientRequest) {
	r.Logger.Debug().Str("key", string(r.Data[4])).Msg("hlen handler")

	data, err := r.Store.HGet(string(r.Data[4]))

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
