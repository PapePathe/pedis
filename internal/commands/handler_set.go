package commands

import (
	"strconv"
)

type SetHandler struct{}

func (ch SetHandler) Authorize(IClientRequest) error {
	return nil
}

func (ch SetHandler) Permissions(IClientRequest) []string {
	return nil
}

func (ch SetHandler) Persistent(IClientRequest) bool {
	return true
}

func (ch SetHandler) Handle(r IClientRequest) {
  body := r.Body()
	key := body[0] 
	if len(key) == 0 {
		_ = r.WriteError("key is empty")
		return
	}

  if len(body) < 2 {
		_ = r.WriteError("value is required")
    return
  }

	value := body[1] 
	if len(value) == 0 {
		_ = r.WriteError("value is empty")
		return
	}

	exp := 0
	if len(body) >= 4 {
		var err error
		exp, err = strconv.Atoi(body[3])
		if err != nil {
			_ = r.WriteError("expiration cannot be casted to number")
			return
		}
	}

	err := r.Store().Set(key, value, int64(exp))

	if err != nil {
		_ = r.WriteError(err.Error())
		return
	}

	_  = r.WriteOK()
}

func init() {
	RegisterCommand("set", SetHandler{})
}
