package commands

type PingHandler struct{}

func (ch PingHandler) Authorize(IClientRequest) error {
	return nil
}

func (ch PingHandler) Permissions(IClientRequest) []string {
	return []string{}
}

func (ch PingHandler) Persistent(IClientRequest) bool {
	return false
}

func (ch PingHandler) Handle(r IClientRequest) {
	data := r.DataRaw().ReadArray()

	if len(data) == 0 {
		r.WriteString("PONG")
		return
	}

	r.WriteString(data[0])
}

func init() {
	RegisterCommand("ping", PingHandler{})
}
