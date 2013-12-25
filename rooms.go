package citadel

import (
	"fmt"
	"time"
)

type room struct {
	Name     string    // Actual name of this room; may include '\' to separate trees
	Flag     int       // QRFlags of this room
	Floor    floor     // The floor number its on
	NoUnread int       // 1   Number of unread messages in this room
	NoTotal  int       // 2   Number of total messages in this room
	Modified time.Time // 15  From Server
	Updated  time.Time // Time this was fatched by us
}

// Enter a specific room by name
func (c *Citadel) Goto(ro string) (ok bool) {
	if c.Request(fmt.Sprintf("GOTO %s", ro)) {
		no := len(c.Resp)
		if no > 14 {
			r := new(room)
			r.Name = c.Resp[0]           // From Citadel to get correct Upper/Lower case
			r.Updated = time.Now().UTC() // email servers should always run in UTC mode
			r.Modified, ok = StrToTime(c.Resp[15])
			ok = true
		}
	}
	return
}

// Retrieve modification time of current room
func (c *Citadel) RoomsStat() (ok bool) {
	if c.Request("STAT") {
		if c.Resp[0] == c.Room.Name {
			c.Room.Updated = time.Now().UTC() // a email servers should always run in UTC mode
			c.Room.Modified, ok = StrToTime(c.Resp[1])
			ok = true
		}
	}
	return
}

// List all accessible Rooms
func (c *Citadel) RoomsAll() bool {
	return c.rooms("LRMS")
}

// List all Public Rooms
func (c *Citadel) RoomsPublic() bool {
	return c.rooms("LPRM")
}

func (c *Citadel) rooms(code string) (ok bool) {
	c.Request(code)
	if c.Code == CODE_LISTING_FOLLOWS {
		res := c.Responce()
		ok = true
		no := len(res)
		if no != 0 {
			for i := 0; i < no; i++ {
				r := res[i]
				flag, _ := StrToInt(r[1])
				floor, _ := StrToInt(r[2])
				Rooms.List = append(Rooms.List, room{Name: r[0], Flag: flag, Floor: Floors[floor]})
			}
		}
	}
	return
}
