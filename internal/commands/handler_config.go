package commands

import (
	"fmt"
	"log"
	"strings"
)

type ConfigHandler struct{}

func (ch ConfigHandler) Authorize(IClientRequest) error {
	return nil
}

func (ch ConfigHandler) Permissions(IClientRequest) []string {
	return nil
}

func (ch ConfigHandler) Persistent(IClientRequest) bool {
	return false
}

func (ch ConfigHandler) HandlePipelined(r IClientRequest) []byte {
	configItems := map[string]string{
		"save":       "3600 1 300 100 60 10000",
		"appendonly": "no",
	}
	params := r.DataRaw().ReadArray()
	switch strings.ToLower(params[0]) {
	case "get":
		v, ok := configItems[params[1]]
		if ok {
			return r.ArrayResponse([]string{params[1], v})
		}

		return r.ErrorResponse("config element not found")
	default:
		return r.ErrorResponse(fmt.Sprintf("subcommand (%s) not supported", string(params[0])))
	}
}

func (ch ConfigHandler) Handle(r IClientRequest) {
	//	configItems := map[string]string{
	//		"save": "",
	//	}
	params := r.DataRaw().ReadArray()
	switch params[0] {
	case "get":
		log.Println("get config item")
		r.WriteArray([]string{"save", "cccc"})
	default:
		r.WriteError(fmt.Sprintf("subcommand (%s) not supported", string(params[0])))
	}
}

func init() {
	RegisterCommand("config", ConfigHandler{})
}
