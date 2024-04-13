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

	user, err := r.Store.GetUser(data[0])
	if err != nil {
		r.WriteError(err.Error())
		return
	}

	if user.AnyPassword {
		r.WriteOK()
		return
	}

	if len(data) == 1 {
		r.WriteError("Password must be supplied")
		return
	}

	if err := user.Authenticate(data[1]); err != nil {
		r.WriteError(err.Error())
		return
	}

	r.WriteOK()
}

func init() {
	RegisterCommand("auth", AuthHandler{})
}
