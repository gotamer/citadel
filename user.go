package citadel

import (
	"fmt"
	"regexp"
)

const (
	Rx_USERNAME      = "^[A-Za-z0-9_-]{4,}$"
	Rx_USERNAME_TRIM = "_- "
	Rx_PASSWORD      = "^.{5,}$"
	Rx_SPACE_TRIM    = " "
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
	c.request(cmd)
	if ok = c.code(); !ok {
		return
	}

	ok = c.userSetPassword(password)
	return
}

func (c *Citadel) Login(username, password string) (ok bool) {
	cmd := fmt.Sprintf("USER %s", username)
	c.request(cmd)
	ok = c.code()
	if ok {
		cmd := fmt.Sprintf("PASS %s", password)
		c.request(cmd)
		ok = c.code()
	}
	return
}

func (c *Citadel) Logout() {
	c.request("LOUT")
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
	c.request(cmd)
	return c.code()
}

func Validate(t, rx string) bool {
	var validator = regexp.MustCompile(rx)
	return validator.MatchString(t)
}
