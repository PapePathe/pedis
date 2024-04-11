package commands

type ConfigHandler struct{}

func (ch ConfigHandler) Authorize(ClientRequest) error {
	return nil
}

func (ch ConfigHandler) Permissions() []string {
	return nil
}

func (ch ConfigHandler) Persistent() bool {
	return false
}

func (ch ConfigHandler) Handle(r ClientRequest) {
	r.WriteError("not yet implemented")
}
