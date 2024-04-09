package commands

func HKeysHandler(r ClientRequest) {
	r.Logger.Debug().Str("key", string(r.Data[4])).Msg("hkeys handler")

	data, err := r.Store.HGet(string(r.Data[4]))

	if err != nil {
		_ = r.WriteArray([]string{})
		return
	}

	hs := hset{}
	hs.FromBytes(data)

	_ = r.WriteArray(hs.Keys())
}
