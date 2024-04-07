package commands

import "log"

func HGetHandler(r ClientRequest) {
	log.Println("hget key", string(r.Data[4]))

	data, err := r.store.HGet(string(r.Data[4]))

	if err != nil {
		_ = r.WriteNil()
		return
	}

	hs := hset{}
	hs.FromBytes(data)

	value, err := hs.Get(string(r.Data[6]))

	if err != nil {
		_ = r.WriteNil()
		return
	}

	_ = r.WriteString(value)
}
