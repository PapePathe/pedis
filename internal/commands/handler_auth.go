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

	if len(data) < 2 {
		r.WriteError("password is required")
	}

	_, err := r.Store.GetUser(data[0])
	if err != nil {
		r.WriteError(err.Error())
	}

	r.WriteOK()
}
