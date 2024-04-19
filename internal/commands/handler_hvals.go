package commands

type HValsHandler struct{}

func (ch HValsHandler) Authorize(IClientRequest) error {
	return nil
}

func (ch HValsHandler) Permissions(IClientRequest) []string {
	return nil
}

func (ch HValsHandler) Persistent(IClientRequest) bool {
	return false
}

func (ch HValsHandler) Handle(r IClientRequest) {
	body := r.Body()
	data, err := r.Store().HGet(body[0])

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
