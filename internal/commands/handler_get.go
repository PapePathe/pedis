package commands

import (
	"pedis/internal/renderer"
)

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
	val, err := r.Store().Get(string(r.Data()[4]))
	if err != nil {
		r.Write([]byte("-ERR key not found\r\n"))
	}

	rr := renderer.BulkStringRenderer{}
	r.Write(rr.Render(val))
}

func init() {
	RegisterCommand("get", GetHandler{})
}
