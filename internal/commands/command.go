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

		// hash related commands
		"hexists": HExistsHandler{},
		"hget":    HGetHandler{},
		"hkeys":   HKeysHandler{},
		"hlen":    HLenHandler{},
		"hset":    HSetHandler{},
		"hvals":   HValsHandler{},
		"config":  ConfigHandler{},
	}
	return &RequestHandler{subcommands}
}

func (s RequestHandler) Run(request ClientRequest) {
	subcommand := strings.ToLower(string(request.Data[2]))

	if h, ok := s.subcommands[subcommand]; ok {
		go h.Handle(request)
	} else {
		switch subcommand {
		case "hello":
			go HelloHandler(request.Data, request.Store, request.Conn)
		case "get":
			go GetHandler(request.Data, request.Store, request.Conn)
		case "set":
			go SetHandler(request.Data, request.Store, request.Conn)
		case "client":
			request.WriteString("OK")
			request.Logger.Debug().Msg("going to execute client options command")
		default:
			request.WriteError(fmt.Sprintf("command not supported %v", string(request.Data[2])))
			request.Logger.Debug().Str("command", string(request.Data[2])).Msg("is not yet supported")
		}
	}
}
