// Package bitbucket.org/gotamer/citadel is a Citadel Client library to access Citadel email and collaboration servers from Go using the Citadel Protocol.
// http://www.citadel.org
package citadel

import (
	"fmt"
	net "net/textproto"
	"os"
	"strconv"
	"strings"

	"bitbucket.org/gotamer/errors"
)

const (
	VERSION = 0.1

	DS = "|"
	DE = "000"

	VIEW_BBS         = "0" // Bulletin board view
	VIEW_MAILBOX     = "1" // Mailbox summary
	VIEW_ADDRESSBOOK = "2" // Address book view
	VIEW_CALENDAR    = "3" // Calendar view
	VIEW_TASKS       = "4" // Tasks view
	VIEW_NOTES       = "5" // Notes view
	VIEW_WIKI        = "6" // Wiki view

	CODE_DONE            = 0
	CODE_LISTING_FOLLOWS = 1 // The requested operation is progressing and is now delivering text. The client *must* now read lines of text until it receives the termination sequence (“000” on a line by itself).
	CODE_OK              = 2 // The requested operation succeeded.
	CODE_MORE_DATA       = 3 // The requested operation succeeded so far, but another command is required to complete it.
	CODE_SEND_LISTING    = 4 // The requested operation is progressing and is now expecting text. The client *must* now transmit zero or more lines of text followed by the termination sequence (“000” on a line by itself).
	CODE_ERROR           = 5 // The requested operation failed. The second and third digits of the error code and/or the error message following it describes why.
	CODE_BINARY_FOLLOWS  = 6 // After this line please read n bytes. (n follows after a blank)
	CODE_SEND_BINARY     = 7 // you now may send us n bytes binary data. (n follows after a blank)
	CODE_START_CHAT_MODE = 8 // ok, we are in chat mode now. every line you send will be broadcasted.
	CODE_ASYNC_MSG       = 9 // there is a page waiting for you, please fetch it.

	MESG_OK                      = 0
	MESG_ASYNC_GEXP              = 02
	MESG_INTERNAL_ERROR          = 10
	MESG_TOO_BIG                 = 11
	MESG_ILLEGAL_VALUE           = 12
	MESG_NOT_LOGGED_IN           = 20
	MESG_CMD_NOT_SUPPORTED       = 30
	MESG_SERVER_SHUTTING_DOWN    = 31
	MESG_PASSWORD_REQUIRED       = 40
	MESG_ALREADY_LOGGED_IN       = 41
	MESG_USERNAME_REQUIRED       = 42
	MESG_HIGHER_ACCESS_REQUIRED  = 50
	MESG_MAX_SESSIONS_EXCEEDED   = 51
	MESG_RESOURCE_BUSY           = 52
	MESG_RESOURCE_NOT_OPEN       = 53
	MESG_NOT_HERE                = 60
	MESG_INVALID_FLOOR_OPERATION = 61
	MESG_NO_SUCH_USER            = 70
	MESG_FILE_NOT_FOUND          = 71
	MESG_ROOM_NOT_FOUND          = 72
	MESG_NO_SUCH_SYSTEM          = 73
	MESG_ALREADY_EXISTS          = 74
	MESG_MESSAGE_NOT_FOUND       = 75
)

type Citadel struct {
	Conn  *net.Conn
	Room  room  // Current room data
	Floor floor // Current floor data
	Code  int   // Citadel reponce CODE_XXXX
	Mesg  int   // Citadel responce MESG_XXXX
	Resp  []string
	Raw   string // Raw responce from citadel
	Error error
}

func New(addr string) (c *Citadel) {
	c = new(Citadel)
	c.Open(addr)
	return
}

func (c *Citadel) Open(addr string) {
	var err error
	c.Conn, err = net.Dial("tcp", addr)
	e.Check(err)
	_, c.Error = c.Conn.ReadLine()
	c.Check()
	c.Iden()
}

func (c *Citadel) Close() {
	c.Request("QUIT")
	err := c.Conn.Close()
	e.Check(err)
}

func (c *Citadel) Iden() {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}
	cmd := fmt.Sprintf("IDEN 1|1|%v|GoTamer|%v", VERSION, hostname)
	c.Request(cmd)
}

// Makes a request and returns the 1st responce
func (c *Citadel) Request(cmd string) (ok bool) {
	var err error
	e.Info(cmd)
	err = c.Conn.PrintfLine("%s", cmd)
	e.Check(err)
	c.Raw, c.Error = c.Conn.ReadLine()
	c.Check()
	c.Raw = strings.Trim(c.Raw, " |")
	c.Code, err = strconv.Atoi(c.Raw[0:1])
	e.Check(err)
	if c.Code != 0 && len(c.Raw) > 2 {
		c.Mesg, err = strconv.Atoi(c.Raw[1:3])
		e.Check(err)
	}
	if len(c.Raw) > 4 {
		c.Resp = strings.Split(c.Raw[4:], DS)
		e.Info("Resp: %v", c.Resp)
	}
	ok = true
	return
}

// This is for multi line responces.
// Used when receiving CODE_LISTING_FOLLOWS
func (c *Citadel) Responce() (r [][]string) {
	var err error
	var text string
	for {
		text, err = c.Conn.ReadLine()
		if e.Check(err) {
			if text == DE {
				break
			}
			r = append(r, strings.Split(text, "|"))
		}
	}
	return
}

func (c *Citadel) code() (ok bool) {
	switch c.Code {
	case CODE_LISTING_FOLLOWS:
		c.Error = nil
		ok = true

	case CODE_OK:
		c.Error = nil
		ok = true

	case CODE_MORE_DATA:
		c.Error = nil
		ok = true

	case CODE_SEND_LISTING:
		c.Error = nil
		ok = true

	case CODE_ERROR:
		c.setError()

	case CODE_BINARY_FOLLOWS:
		c.Error = nil
		ok = true

	case CODE_SEND_BINARY:
		c.Error = nil
		ok = true

	case CODE_START_CHAT_MODE:
		c.Error = nil
		ok = true

	case CODE_ASYNC_MSG:
		c.Error = nil
		ok = true
	default:
		c.setError()
	}
	c.Check()
	return
}

/*

	switch c.Mesg {
	case MESG_ALREADY_EXISTS:
		log.Println("Mesg: ", c.Mesg)

	case MESG_ALREADY_LOGGED_IN:
		log.Println("Mesg: ", c.Mesg)

	case MESG_ASYNC_GEXP:
		log.Println("Mesg: ", c.Mesg)

	case MESG_CMD_NOT_SUPPORTED:
		log.Println("Mesg: ", c.Mesg)

	case MESG_FILE_NOT_FOUND:
		log.Println("Mesg: ", c.Mesg)

	case MESG_HIGHER_ACCESS_REQUIRED:
		log.Println("Mesg: ", c.Mesg)

	case MESG_ILLEGAL_VALUE:
		log.Println("Mesg: ", c.Mesg)

	case MESG_INTERNAL_ERROR:
		log.Println("Mesg: ", c.Mesg)

	case MESG_INVALID_FLOOR_OPERATION:
		log.Println("Mesg: ", c.Mesg)

	case MESG_MAX_SESSIONS_EXCEEDED:
		log.Println("Mesg: ", c.Mesg)

	case MESG_MESSAGE_NOT_FOUND:
		log.Println("Mesg: ", c.Mesg)

	case MESG_NOT_HERE:
		log.Println("Mesg: ", c.Mesg)

	case MESG_NOT_LOGGED_IN:
		log.Println("Mesg: ", c.Mesg)

	case MESG_NO_SUCH_SYSTEM:
		log.Println("Mesg: ", c.Mesg)

	case MESG_NO_SUCH_USER:
		log.Println("Mesg: ", c.Mesg)

	case MESG_OK:
		log.Println("Mesg: ", c.Mesg)

	case MESG_PASSWORD_REQUIRED:
		log.Println("Mesg: ", c.Mesg)

	case MESG_RESOURCE_BUSY:
		log.Println("Mesg: ", c.Mesg)

	case MESG_RESOURCE_NOT_OPEN:
		log.Println("Mesg: ", c.Mesg)

	case MESG_ROOM_NOT_FOUND:
		log.Println("Mesg: ", c.Mesg)

	case MESG_SERVER_SHUTTING_DOWN:
		log.Println("Mesg: ", c.Mesg)

	case MESG_TOO_BIG:
		log.Println("Mesg: ", c.Mesg)

	case MESG_USERNAME_REQUIRED:
		log.Println("Mesg: ", c.Mesg)
	default:
		log.Println("Unknown Message: ")
		c.Error = fmt.Errorf("CIT CODE: %v CIT MESG: %v %s", c.Code, c.Mesg, c.Resp)
	}

*/
