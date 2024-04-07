package commands

import "log"

func HExistsHandler(r ClientRequest) {
	log.Println("hexists key", string(r.Data[4]))

	data, err := r.store.HGet(string(r.Data[4]))
	if err != nil {
		_ = r.WriteNumber("0")
		return
	}

	hs := hset{}
	hs.FromBytes(data)

	_, err = hs.Get(string(r.Data[6]))
	if err != nil {
		_ = r.WriteNumber("0")
		return
	}

	_ = r.WriteNumber("1")
}
