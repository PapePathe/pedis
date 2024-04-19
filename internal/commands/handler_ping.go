package commands

import "log"

type PingHandler struct{}

func (ch PingHandler) Authorize(IClientRequest) error {
	return nil
}

func (ch PingHandler) Permissions(IClientRequest) []string {
	return []string{
		"fast",
		"connection",
	}
}

func (ch PingHandler) Persistent(IClientRequest) bool {
	return false
}

func (ch PingHandler) Handle(r IClientRequest) {
	if len(r.Body()) == 0 {
		err := r.WriteString("PONG")
		if err != nil {
			log.Println("error writing to client", err)
		}
		return
	}

	err := r.WriteString(r.Body()[0])
	if err != nil {
		log.Println("error writing to client", err)
	}
}

func init() {
	RegisterCommand("ping", PingHandler{})
}
