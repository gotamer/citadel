package citadel

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"bitbucket.org/gotamer/errors"
)

const (
	Rx_USERNAME      = "^[A-Za-z0-9_-]{4,}$"
	Rx_USERNAME_TRIM = "_- "
	Rx_PASSWORD      = "^.{5,}$"
	Rx_SPACE_TRIM    = " "
)

func (c *Citadel) Check() (ok bool) {
	if c.Error != nil {
		ok = false
		c.Info()
	}
	return
}

// You can use this to see information about the last command executed
func (c *Citadel) Info() {
	e.Info("\nInfo CODE: %v MESG: %v: %v\n\tResponce RAW: %s\n\n", c.Code, c.Mesg, c.Resp, c.Raw)
}

func (c *Citadel) setError() {
	c.Error = fmt.Errorf("CIT CODE: %v MESG: %v: %v\n\tResponce RAW: %s\n\n", c.Code, c.Mesg, c.Resp, c.Raw)
}

func StrToInt(s string) (i int, ok bool) {
	i, err := strconv.Atoi(s) // int
	if err != nil {
		log.Println("Error StrToInt: ", err.Error())
	} else {
		ok = true
	}
	return
}

func StrToTime(s string) (t time.Time, ok bool) {
	x, err := strconv.ParseInt(s, 10, 0) // int64
	if err != nil {
		log.Println("Error StrToTime: ", err.Error())
	} else {
		t = time.Unix(x, 0)
		ok = true
	}
	return
}

func Validate(t, rx string) bool {
	var validator = regexp.MustCompile(rx)
	return validator.MatchString(t)
}
