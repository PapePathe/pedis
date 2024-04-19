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
	body := r.Body()
	data, err := r.Store().HGet(body[0])

	if err != nil {
		_ = r.WriteNil()
		return
	}

	hs := hset{}
	hs.FromBytes(data)

	value, err := hs.Get(body[1])

	if err != nil {
		_ = r.WriteNil()
		return
	}

	_ = r.WriteString(value)
}

func init() {
	RegisterCommand("hget", HGetHandler{})
}
