package commands

import "fmt"
import "log"

// ACL categories: @keyspace, @write, @slow

type DelHandler struct{}

func (ch DelHandler) Authorize(IClientRequest) error {
	return nil
}

func (ch DelHandler) Permissions(IClientRequest) []string {
	return nil
}

func (ch DelHandler) Persistent(IClientRequest) bool {
	return true
}

func (ch DelHandler) Handle(r IClientRequest) {
	delCount := 0

	for _, key := range r.Body() {
		err := r.Store().Del(key)

		if err == nil {
      delCount++
		}
	}

  log.Println("del-count", delCount)

	r.WriteNumber(fmt.Sprintf("%d", delCount))
}

func init() {
	RegisterCommand("del", DelHandler{})
}
