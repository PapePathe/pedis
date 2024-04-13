package commands

type ConfigHandler struct{}

func (ch ConfigHandler) Authorize(IClientRequest) error {
	return nil
}

func (ch ConfigHandler) Permissions(IClientRequest) []string {
	return nil
}

func (ch ConfigHandler) Persistent(IClientRequest) bool {
	return false
}

func (ch ConfigHandler) Handle(r IClientRequest) {
	r.WriteError("not yet implemented")
}

func init() {
	RegisterCommand("config", ConfigHandler{})
}
