package commands

import (
	"fmt"
)

func HLenHandler(r ClientRequest) {
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
