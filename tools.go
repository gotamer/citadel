package citadel

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type env uint8

const (
	ENV_PROD env = iota // Environment
	ENV_INFO
	ENV_FAIL
)

var (
	ENVIROMENT env = ENV_FAIL
	LogError   *log.Logger
	LogDebug   *log.Logger
)

func init() {
	Logger(os.Stderr)
}

func Logger(to *os.File) {
	LogError = log.New(to, "", 19)
	LogError.SetPrefix("Error: ")
	LogDebug = log.New(to, "", 19)
	LogDebug.SetPrefix("Info: ")
}

func Check(err error) (ok bool) {
	if err == nil {
		ok = true
	} else {
		LogError.Output(2, err.Error())
		if ENVIROMENT == ENV_FAIL {
			os.Exit(2)
		}
	}
	return
}

func Debug(m interface{}) {
	if ENVIROMENT != ENV_PROD {
		LogDebug.Output(2, fmt.Sprintf("%s", m))
	}
}

func (c *Citadel) setError() {
	c.Error = fmt.Errorf("CIT CODE: %v MESG: %v: %s", c.Code, c.Mesg, c.Resp)
}

func (c *Citadel) Check() (ok bool) {
	if c.Error == nil {
		ok = true
	} else {
		LogError.Output(2, c.Error.Error())
		if ENVIROMENT == ENV_FAIL {
			c.Close()
			os.Exit(1)
		}
	}
	return
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
