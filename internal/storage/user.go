package storage

import (
	"errors"
)

type AclRuleType string

var (
	AclActivateUser    AclRuleType = "on"
	AclDisableUser     AclRuleType = "off"
	AclResetUser       AclRuleType = "reset"
	AclSetUserPassword AclRuleType = "setuserpassword"
	AclAnyPassword     AclRuleType = "nopass"
	AclClearSelectors  AclRuleType = "clearselectors"
	AclNoCommands      AclRuleType = "nocommands"
)

type AclRule struct {
	Type  AclRuleType
	Value string
}

type User struct {
	Passwords   []string
	Active      bool
	AnyPassword bool
}

func (u User) Authenticate(password string) error {
	for _, p := range u.Passwords {
		if p == password {
			return nil
		}
	}
	return errors.New("Password auth failed")
}

func (u User) Authorize(command string) error {
	return nil
}
