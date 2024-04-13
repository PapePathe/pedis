package commands

import (
	"fmt"
	"pedis/internal/storage"
)

type AclHandler struct{}

func (ch AclHandler) Authorize(IClientRequest) error {
	return nil
}

func (ch AclHandler) Permissions(IClientRequest) []string {
	return nil
}

func (ch AclHandler) Persistent(IClientRequest) bool {
	return false
}

func (ch AclHandler) Handle(r IClientRequest) {
	data := r.DataRaw().ReadArray()
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

func (aclService) deluser(r IClientRequest) error {
	data := r.DataRaw().ReadArray()
	delCount := 0

	for _, u := range data[1:] {
		err := r.Store().DelUser(u)

		if err == nil {
			delCount++
		}
	}

	r.WriteNumber(fmt.Sprintf("%d", delCount))

	return nil
}

func (aclService) setuser(r IClientRequest) error {
	data := r.DataRaw().ReadArray()
	username := data[1]
	rules := []storage.AclRule{}

	if len(data) >= 2 {
		for _, elem := range data[2:] {
			switch elem {
			case "on":
				rules = append(rules, storage.AclRule{Type: storage.AclActivateUser})
			case "off":
				rules = append(rules, storage.AclRule{Type: storage.AclDisableUser})
			case "nopass":
				rules = append(rules, storage.AclRule{Type: storage.AclDisableUser})
			case "reset":
				rules = append(rules, storage.AclRule{Type: storage.AclResetUser})
			default:
				switch elem[0] {
				case '>':
					rules = append(rules, storage.AclRule{Type: storage.AclSetUserPassword, Value: elem[1 : len(elem)-1]})
				default:
					r.WriteError(fmt.Sprintf("acl rule (%s) not supported", elem))
				}
			}
		}
	}

	_ = r.Store().SetUser(username, rules)
	r.WriteOK()
	return nil
}

func (aclService) users(r IClientRequest) error {
	users := r.Store().Users()
	r.WriteArray(users)

	return nil
}

func init() {
	RegisterCommand("acl", AclHandler{})
}
