package commands

func HExistsHandler(r ClientRequest) {
	r.Logger.Debug().Str("key", string(r.Data[4])).Msg("hexists handler")

	data, err := r.Store.HGet(string(r.Data[4]))
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
