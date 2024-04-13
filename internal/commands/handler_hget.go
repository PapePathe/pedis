package commands

type HGetHandler struct{}

func (ch HGetHandler) Authorize(ClientRequest) error {
	return nil
}

func (ch HGetHandler) Permissions() []string {
	return nil
}

func (ch HGetHandler) Persistent() bool {
	return false
}

func (ch HGetHandler) Handle(r ClientRequest) {
	r.Logger.Debug().Interface("command", r.DataRaw.ReadArray()).Msg("hget handler")

	datat := r.DataRaw.ReadArray()
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

func init() {
	RegisterCommand("hget", HGetHandler{})
}
