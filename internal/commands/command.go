package commands

type RequestHandler struct {
	subcommands map[string]CommandHandler
}

func NewRequestHandler() *RequestHandler {
	return &RequestHandler{}
}

func DefaultRequestHandler() *RequestHandler {
	subcommands := map[string]CommandHandler{
		"hexists": HExistsHandler,
		"hget":    HGetHandler,
		"hkeys":   HKeysHandler,
		"hlen":    HLenHandler,
		"hset":    HSetHandler,
		"hvals":   HValsHandler,
	}
	return &RequestHandler{subcommands}
}

func (s RequestHandler) Run(request ClientRequest) {
	subcommand := string(request.Data[2])

	if h, ok := s.subcommands[subcommand]; ok {
		go h(request)
	} else {
		switch string(request.Data[2]) {
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
			request.WriteError("command not supported")
			request.Logger.Debug().Str("command", string(request.Data[2])).Msg("is not yet supported")
		}
	}
}
