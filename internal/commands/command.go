package commands

import (
	"fmt"
	"strings"
)

type RequestHandler struct {
	subcommands map[string]CommandHandler
}

func NewRequestHandler() *RequestHandler {
	return &RequestHandler{}
}

func DefaultRequestHandler() *RequestHandler {
	subcommands := map[string]CommandHandler{
		// cluster commands
		"cluster": ClusterHandler{},

		// acl commands
		"acl": AclHandler{},

		"del": DelHandler{},
		"get": GetHandler{},
		"set": SetHandler{},

		// hash related commands
		"hexists": HExistsHandler{},
		"hget":    HGetHandler{},
		"hkeys":   HKeysHandler{},
		"hlen":    HLenHandler{},
		"hset":    HSetHandler{},
		"hvals":   HValsHandler{},

		"config": ConfigHandler{},
		"hello":  HelloHandler{},
		"auth":   AuthHandler{},
		"ping":   PingHandler{},
	}
	return &RequestHandler{subcommands}
}

func (s RequestHandler) Run(request ClientRequest) {
	subcommand := strings.ToLower(string(request.Data[2]))

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
			request.Logger.Debug().Msg("going to execute client options command")
		default:
			request.WriteError(fmt.Sprintf("command not supported %v", subcommand))
			request.Logger.Debug().Str("command", subcommand).Msg("is not yet supported")
		}
	}
}
