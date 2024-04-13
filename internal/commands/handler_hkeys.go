package commands

type HKeysHandler struct{}

func (ch HKeysHandler) Authorize(ClientRequest) error {
	return nil
}

func (ch HKeysHandler) Permissions() []string {
	return nil
}

func (ch HKeysHandler) Persistent() bool {
	return false
}

func (ch HKeysHandler) Handle(r ClientRequest) {
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

func init() {
	RegisterCommand("hkeys", HKeysHandler{})
}
