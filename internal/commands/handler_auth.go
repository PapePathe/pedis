package commands

type AuthHandler struct{}

func (ch AuthHandler) Authorize(ClientRequest) error {
	return nil
}

func (ch AuthHandler) Permissions() []string {
	return nil
}

func (ch AuthHandler) Persistent() bool {
	return false
}

func (ch AuthHandler) Handle(r ClientRequest) {
	data := r.DataRaw.ReadArray()
	r.Logger.Info().Interface("auth params", data).Msg("")
	r.WriteError("not yet deployed")
}
