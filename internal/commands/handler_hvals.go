package commands

import (
	"log"
)

func HValsHandler(r ClientRequest) {
	log.Println("hkeys handler")

	data, err := r.store.HGet(string(r.Data[4]))

	if err != nil {
		_ = r.WriteArray([]string{})
		return
	}

	hs := hset{}
	hs.FromBytes(data)

	_ = r.WriteArray(hs.Values())
}
