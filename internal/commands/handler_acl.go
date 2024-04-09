package commands

import (
	"fmt"
	"pedis/internal/storage"
)

func AclHandler(r ClientRequest) {
	data := r.DataRaw.ReadArray()
	r.Logger.
		Debug().
		Interface("Data", data).
		Interface("RawData", r.DataRaw.String()).
		Msg("")

	svc := aclService{}

	switch data[0] {
	case "setuser":
		_ = svc.setuser(r)
	case "deluser":
		_ = svc.deluser(r)
	case "users":
		_ = svc.users(r)
	default:
		r.WriteError(fmt.Sprintf("(%s) not implemented by devin", data[0]))
	}
}

type aclService struct{}

func (aclService) deluser(r ClientRequest) error {
	data := r.DataRaw.ReadArray()
	r.Logger.Debug().Interface("usernames", data[1:]).Msg("Going to delete")
	delCount := 0

	for _, u := range data[1:] {
		err := r.Store.DelUser(u)

		if err == nil {
			delCount++
		}
	}

	r.WriteNumber(fmt.Sprintf("%d", delCount))

	return nil
}

func (aclService) setuser(r ClientRequest) error {
	data := r.DataRaw.ReadArray()
	r.Logger.Debug().Msg("Going to create or update existing user")
	username := data[1]

	_ = r.Store.SetUser(username, []storage.AclRule{})
	r.WriteOK()
	return nil
}

func (aclService) users(r ClientRequest) error {
	r.Logger.Debug().Msg("Going to list users")

	users := r.Store.Users()
	r.Logger.Debug().Interface("Users list", users).Msg("")
	r.WriteArray(users)

	return nil
}
