package commands

type GetHandler struct{}

func (ch GetHandler) Authorize(IClientRequest) error {
	return nil
}

func (ch GetHandler) Permissions(IClientRequest) []string {
	return nil
}

func (ch GetHandler) Persistent(IClientRequest) bool {
	return true
}

func (ch GetHandler) Handle(r IClientRequest) {
  body := r.Body()
	val, err := r.Store().Get(body[0])
	if err != nil {
		r.WriteError(err.Error())
	}

	r.WriteString(val)
}

func init() {
	RegisterCommand("get", GetHandler{})
}
