package commands

type HExistsHandler struct{}

func (ch HExistsHandler) Authorize(ClientRequest) error {
	return nil
}

func (ch HExistsHandler) Permissions() []string {
	return nil
}

func (ch HExistsHandler) Persistent() bool {
	return false
}

func (ch HExistsHandler) Handle(r ClientRequest) {
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

func init() {
	RegisterCommand("hexists", HExistsHandler{})
}
