package commands

func HGetHandler(r ClientRequest) {
	r.Logger.Debug().Str("hset key", string(r.Data[4])).Msg("hget handler")

	data, err := r.Store.HGet(string(r.Data[4]))

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
