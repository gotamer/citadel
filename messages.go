package citadel

import (
	"fmt"
	"log"
	"strconv"
	"strings"
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

type Message struct {
	EUID   string
	Mime   string // text/vcard
	UxTime int
	Text   string
}

type Dav struct {
	UID    string
	REV    string
	FN     string
	Object string
}

// Get vcard from current room by MSG ID
func (c *Citadel) GetMessage(MsgID, Mime string) (msg *Message, ok bool) {
	msg = new(Message)
	msg.Mime = Mime
	if c.Request(fmt.Sprintf("MSGP %v", Mime)) { // "MSGP text/vcard"
		var err error
		if c.Request(fmt.Sprintf("MSG4 %v", MsgID)) {
			if c.Code == CODE_LISTING_FOLLOWS {
				res := c.Responce()
				no := len(res)
				for i := 0; i < no; i++ {
					//fmt.Println(res[i][0])
					if ok == false {
						if strings.HasPrefix(res[i][0], "msgn=") {
							msg.EUID = TrimPrefix(res[i][0], "msgn=")
						}
						// Waitting on citadel support on this, what is exti=?
						/*
							if strings.HasPrefix(res[i][0], "exti=") {
								msg.EUID = strings.TrimPrefix(res[i][0], "exti=")
							}
						*/
						if strings.HasPrefix(res[i][0], "time=") {
							ts := TrimPrefix(res[i][0], "time=")
							if msg.UxTime, err = strconv.Atoi(ts); err != nil {
								log.Println(err)
							}
						}
						if res[i][0] == "text" {
							ok = true
						}
					} else {
						if strings.HasPrefix(res[i][0], "Content-type:") {
							msg.Mime = TrimPrefix(res[i][0], "Content-type:")
						}
						msg.Text += res[i][0] + "\n"
					}
				}
			}
		}
	}
	return
}

func TrimPrefix(s, p string) string {
	return strings.TrimSpace(strings.TrimPrefix(s, p))
}

func (dav *Dav) ParseDav(text string) {
	list := strings.Split(text, "\n")
	s := 0
	for _, v := range list {
		u := strings.ToUpper(v)
		if u == "BEGIN:VCARD" || u == "BEGIN:VCALENDAR" || u == "BEGIN:VNOTE" && s == 0 {
			v = u
			s = 1
		}
		if s == 1 {
			if strings.HasPrefix(v, "UID:") {
				dav.UID = strings.TrimPrefix(v, "UID:")
			}
			if u == "END:VCARD" || u == "END:VCALENDAR" || u == "END:VNOTE" {
				v = u
				s = 2
			}
			dav.Object += v + "\n"
		}
	}
}

// Get vcard from current room by MSG ID
func (c *Citadel) GetDav(MsgID string, contentType string) (msg *Message, dav *Dav, ok bool) {
	if contentType == "text/vcard" {
		contentType = "text/vcard|text/x-vcard"
	}
	dav = new(Dav)
	if msg, ok = c.GetMessage(MsgID, contentType); ok {
		dav.ParseDav(msg.Text)
	}
	return
}
