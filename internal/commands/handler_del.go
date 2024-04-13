package commands

import "fmt"

// ACL categories: @keyspace, @write, @slow

type DelHandler struct{}

func (ch DelHandler) Authorize(IClientRequest) error {
	return nil
}

func (ch DelHandler) Permissions(IClientRequest) []string {
	return nil
}

func (ch DelHandler) Persistent(IClientRequest) bool {
	return false
}

func (ch DelHandler) Handle(r IClientRequest) {
	delCount := 0

	for _, key := range r.DataRaw().ReadArray() {
		err := r.Store().Del(key)

		if err != nil {
			continue
		}

		delCount++
	}

	r.WriteNumber(fmt.Sprintf("%d", delCount))
}

func init() {
	RegisterCommand("del", DelHandler{})
}
