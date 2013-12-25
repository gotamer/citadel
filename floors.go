package citadel

var Floors []floor

type floor struct {
	Id    int
	Name  string
	Rooms int
}

func (c *Citadel) FloorsLoader() (ok bool) {
	c.Request("LFLR")
	if c.Code == CODE_LISTING_FOLLOWS {
		res := c.Responce()
		ok = true
		no := len(res)
		if no != 0 {
			for i := 0; i < no; i++ {
				r := res[i]
				id, _ := StrToInt(r[0])
				count, _ := StrToInt(r[2])
				Floors = append(Floors, floor{id, r[1], count})
			}
		}
	}
	return
}
