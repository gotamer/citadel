package citadel

import (
	"fmt"
	"strconv"
)

var Floors []Floor

type Floor struct {
	Id    int
	Name  string
	Rooms int
}

func (c *Citadel) Floorer() (ok bool) {
	c.Request("LFLR")
	if c.Code == CODE_LISTING_FOLLOWS {
		res := c.Responce()
		ok = true
		no := len(res)
		if no != 0 {
			for i := 0; i < no; i++ {
				r := res[i]
				id, err := strconv.Atoi(r[0])
				if !Check(err) {
					c.Close()
					return
				}
				no, err := strconv.Atoi(r[2])
				if !Check(err) {
					c.Close()
					return
				}
				Floors = append(Floors, Floor{id, r[1], no})
			}
		}
	}
	return
}
