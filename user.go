package citadel

import (
	"fmt"
)

const (
	USE_LEVEL_DELETED  = "0"
	USE_LEVEL_NEW      = "1"
	USE_LEVEL_PROBLEM  = "2"
	USE_LEVEL_LOCAL    = "3"
	USE_LEVEL_NETWORK  = "4"
	USE_LEVEL_PREFERED = "5"
	USE_LEVEL_ADMIN    = "6"
)

func (c *Citadel) UserCreate(username, password string) (ok bool) {
	if ok = Validate(username, Rx_USERNAME); !ok {
		c.Error = fmt.Errorf("Username not Valid")
		return
	}
	if ok = Validate(password, Rx_PASSWORD); !ok {
		c.Error = fmt.Errorf("Password not Valid")
		return
	}

	cmd := fmt.Sprintf("NEWU %s", username)
	c.Request(cmd)
	if ok = c.code(); !ok {
		return
	}

	ok = c.userSetPassword(password)
	return
}

func (c *Citadel) Login(username, password string) (ok bool) {
	cmd := fmt.Sprintf("USER %s", username)
	c.Request(cmd)
	ok = c.code()
	if ok {
		cmd := fmt.Sprintf("PASS %s", password)
		c.Request(cmd)
		ok = c.code()
	}
	c.FloorsLoader()
	return
}

func (c *Citadel) Logout() {
	c.Request("LOUT")
}

// This command sets a new password for the currently logged in user.
func (c *Citadel) UserSetPassword(password string) (ok bool) {
	if ok = Validate(password, Rx_PASSWORD); !ok {
		c.Error = fmt.Errorf("Password not Valid")
	} else {
		ok = c.userSetPassword(password)
	}
	return
}

func (c *Citadel) userSetPassword(password string) bool {
	cmd := fmt.Sprintf("SETP %s", password)
	c.Request(cmd)
	return c.code()
}

func (c *Citadel) UserCfg(username string) bool {
	cmd := fmt.Sprintf("AGUP %s", username)
	c.Request(cmd)
	return c.code()
}
