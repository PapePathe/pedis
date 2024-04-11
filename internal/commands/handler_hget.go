package commands

func HGetHandler(r ClientRequest) {
	r.Logger.Debug().Interface("command", r.DataRaw.ReadArray()).Msg("hget handler")

	datat := r.DataRaw.ReadArray()

	if len(datat) <= 3 {
		_ = r.WriteError("incomplete command args")
		return
	}

	data, err := r.Store.HGet(datat[0])

	if err != nil {
		_ = r.WriteNil()
		return
	}

	hs := hset{}
	hs.FromBytes(data)

	value, err := hs.Get(datat[1])

	if err != nil {
		_ = r.WriteNil()
		return
	}

	_ = r.WriteString(value)
}
