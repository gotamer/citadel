package citadel

import (
	"fmt"
	//"strconv"
	//"time"
)

func (c *Citadel) MsgList(room string) (ok bool) {
	if c.Request(fmt.Sprintf("GOTO %s", room)) {
		ok = true
	}
	return
}
