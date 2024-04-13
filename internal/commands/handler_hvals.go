package commands

type HValsHandler struct{}

func (ch HValsHandler) Authorize(ClientRequest) error {
	return nil
}

func (ch HValsHandler) Permissions() []string {
	return nil
}

func (ch HValsHandler) Persistent() bool {
	return false
}

func (ch HValsHandler) Handle(r ClientRequest) {
	r.Logger.Debug().Str("key", string(r.Data[4])).Msg("hvals handler")

	data, err := r.Store.HGet(string(r.Data[4]))

	if err != nil {
		_ = r.WriteArray([]string{})
		return
	}

	hs := hset{}
	hs.FromBytes(data)

	_ = r.WriteArray(hs.Values())
}

func init() {
	RegisterCommand("hvals", HValsHandler{})
}
