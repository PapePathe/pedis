package commands

type PingHandler struct{}

func (ch PingHandler) Authorize(ClientRequest) error {
	return nil
}

func (ch PingHandler) Permissions() []string {
	return []string{}
}

func (ch PingHandler) Persistent() bool {
	return false
}

func (ch PingHandler) Handle(r ClientRequest) {
	data := r.DataRaw.ReadArray()
	r.Logger.Info().Str("Cmd", r.DataRaw.String()).Interface("ping params", data).Msg("")

	if len(data) == 0 {
		r.WriteString("PONG")
		return
	}

	r.WriteString(data[0])
}

func init() {
	RegisterCommand("ping", PingHandler{})
}
