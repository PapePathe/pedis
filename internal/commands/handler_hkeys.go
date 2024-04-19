package commands

type HKeysHandler struct{}

func (ch HKeysHandler) Authorize(IClientRequest) error {
	return nil
}

func (ch HKeysHandler) Permissions(IClientRequest) []string {
	return nil
}

func (ch HKeysHandler) Persistent(IClientRequest) bool {
	return false
}

func (ch HKeysHandler) Handle(r IClientRequest) {
	body := r.Body()
	data, err := r.Store().HGet(body[0])

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
