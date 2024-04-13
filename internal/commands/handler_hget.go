package commands

type HGetHandler struct{}

func (ch HGetHandler) Authorize(IClientRequest) error {
	return nil
}

func (ch HGetHandler) Permissions(IClientRequest) []string {
	return nil
}

func (ch HGetHandler) Persistent(IClientRequest) bool {
	return false
}

func (ch HGetHandler) Handle(r IClientRequest) {
	datat := r.DataRaw().ReadArray()
	data, err := r.Store().HGet(datat[0])

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
