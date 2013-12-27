package citadel

import (
	"fmt"
	//"strconv"
	//"time"
)

type msglist string

const (
	LIST_ALL    msglist = "all"
	LIST_OLD    msglist = "old"
	LIST_NEW    msglist = "new"
	LIST_FIRST  msglist = "first"
	LIST_LAST   msglist = "last"
	LIST_GT     msglist = "gt"
	LIST_LT     msglist = "lt"
	LIST_SEARCH msglist = "search"
)

type MsgList struct {
	Cmd       msglist
	FirstLast int
	GtLt      int
	Find      string
}

// Get list of all messages in current room
func (c *Citadel) MsgListAll() (list []string, ok bool) {
	m := new(MsgList)
	m.Cmd = LIST_ALL
	list, ok = c.mlist(m)
	return
}

// Get messages list for current room
// This command must be passed a single parameter: LIST_ALL, LIST_OLD, or LIST_NEW
func (c *Citadel) mlist(m *MsgList) (list []string, ok bool) {
	if c.Request(fmt.Sprintf("MSGS %s", m.Cmd)) {
		if c.Code == CODE_LISTING_FOLLOWS {
			res := c.Responce()
			no := len(res)
			for i := 0; i < no; i++ {
				list = append(list, res[i][0])
			}
			ok = true
		}
	}
	return
}

// Del a list of messages from current room
func (c *Citadel) MsgsDel(list []string) (ok bool) {
	var cmd string
	no := len(list)
	if no != 0 {
		cmd = fmt.Sprintf("DELE %s", list[0])
		if no > 1 {
			for i := 1; i < no; i++ {
				cmd = fmt.Sprintf("%s,%s", cmd, list[i])
			}
		}
		ok = c.Request(cmd)
	}
	return
}
