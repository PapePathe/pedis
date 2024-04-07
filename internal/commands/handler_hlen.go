package commands

import (
	"fmt"
	"log"
)

func HLenHandler(r ClientRequest) {
	log.Println("hlen key", string(r.Data[4]))

	data, err := r.store.HGet(string(r.Data[4]))

	if err != nil {
		_ = r.WriteString(fmt.Sprintf("%d", 0))
		return
	}

	hs := hset{}
	hs.FromBytes(data)

	_ = r.WriteNumber(fmt.Sprintf("%d", hs.Len()))
}
