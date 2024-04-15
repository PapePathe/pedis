package commands

type HelloHandler struct{}

func (ch HelloHandler) Authorize(IClientRequest) error {
	return nil
}

func (ch HelloHandler) Permissions(IClientRequest) []string {
	return nil
}

func (ch HelloHandler) Persistent(IClientRequest) bool {
	return false
}

func (ch HelloHandler) HandlePipelined(r IClientRequest) []byte {
	return []byte{}
}

func (ch HelloHandler) Handle(r IClientRequest) {
	data := r.DataRaw().ReadArray()

	user, err := r.Store().GetUser(data[1])
	if err != nil {
		r.WriteError(err.Error())
		return
	}

	if user.AnyPassword {
		r.WriteOK()
		return
	}

	if len(data) == 2 {
		r.WriteError("Password must be supplied")
		return
	}

	if err := user.Authenticate(data[2]); err != nil {
		r.WriteError(err.Error())
		return
	}

	r.WriteOK()
}

func init() {
	RegisterCommand("hello", HelloHandler{})
}
