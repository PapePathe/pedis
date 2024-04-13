package commands

import (
	"fmt"
	"strings"
	"sync"
)

var (
	defaultCommands     map[string]CommandHandler
	defaultCommandsLock sync.Mutex
)

func RegisterCommand(cmd string, h CommandHandler) {
	defaultCommandsLock.Lock()
	defer defaultCommandsLock.Unlock()

	if defaultCommands == nil {
		defaultCommands = make(map[string]CommandHandler)
	}

	defaultCommands[cmd] = h
}

type RequestHandler struct {
	subcommands map[string]CommandHandler
}

func NewRequestHandler() *RequestHandler {
	return &RequestHandler{}
}

func DefaultRequestHandler() *RequestHandler {
	return &RequestHandler{defaultCommands}
}

func (s RequestHandler) Run(request IClientRequest) {
	subcommand := strings.ToLower(string(request.Data()[2]))

	if h, ok := s.subcommands[subcommand]; ok {
		if err := h.Authorize(request); err != nil {
			request.WriteError("not authorized to access command")
			return
		}
		go h.Handle(request)
	} else {
		switch subcommand {
		case "client":
			request.WriteString("OK")
		default:
			request.WriteError(fmt.Sprintf("command not supported %v", subcommand))
		}
	}
}
