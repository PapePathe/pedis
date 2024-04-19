package commands

type HExistsHandler struct{}

func (ch HExistsHandler) Authorize(IClientRequest) error {
	return nil
}

func (ch HExistsHandler) Permissions(IClientRequest) []string {
	return nil
}

func (ch HExistsHandler) Persistent(IClientRequest) bool {
	return false
}

func (ch HExistsHandler) Handle(r IClientRequest) {
	body := r.Body()
	data, err := r.Store().HGet(body[0])
	if err != nil {
		_ = r.WriteNumber("0")
		return
	}

	hs := hset{}
	hs.FromBytes(data)

	_, err = hs.Get(body[1])
	if err != nil {
		_ = r.WriteNumber("0")
		return
	}

	_ = r.WriteNumber("1")
}

func init() {
	RegisterCommand("hexists", HExistsHandler{})
}
