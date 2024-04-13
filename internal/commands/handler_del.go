package commands

import "fmt"

// ACL categories: @keyspace, @write, @slow

type DelHandler struct{}

func (ch DelHandler) Authorize(ClientRequest) error {
	return nil
}

func (ch DelHandler) Permissions() []string {
	return nil
}

func (ch DelHandler) Persistent() bool {
	return false
}

func (ch DelHandler) Handle(r ClientRequest) {
	r.Logger.Debug().Interface("Data", r.DataRaw.ReadArray()).Str("RawData", r.DataRaw.String()).Msg("del handler")

	delCount := 0

	for _, key := range r.DataRaw.ReadArray() {
		err := r.Store.Del(key)

		if err != nil {
			r.Logger.Error().Str("Key", key).Err(err).Msg("del handler")
			continue
		}

		delCount++
	}

	r.WriteNumber(fmt.Sprintf("%d", delCount))
}

func init() {
	RegisterCommand("del", DelHandler{})
}
