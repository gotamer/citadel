// Sync local files such as vCard, vCalender, vNotes, and text
// files with a local or remote Citadel Mail Server
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	//"strconv"
	"strings"
	"time"

	"bitbucket.org/gotamer/cfg"
	"bitbucket.org/gotamer/citadel"
	"bitbucket.org/gotamer/errors"
	//"bitbucket.org/gotamer/tools"
)

const (
	VERSION      = 3
	VCARD_TIME   = "20060102T150405Z0700"
	DEFAULT_PORT = ":504"
)

const LICENSE = `

	The MIT License (MIT)
	=====================

	Copyright © 2013 Dennis T Kaplan <http://www.robotamer.com>

	Permission is hereby granted, free of charge, to any person
	obtaining a copy of this software and associated documentation
	files (the "Software"), to deal in the Software without restriction,
	including without limitation the rights to use, copy, modify, merge,
	publish, distribute, sub-license, and/or sell copies of the Software,
	and to permit persons to whom the Software is furnished to do so,
	subject to the following conditions:

	The above copyright notice and this permission notice shall be included
	in all copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS
	OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NON-INFRINGEMENT. IN NO EVENT SHALL
	THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
	OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
	ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
	OTHER DEALINGS IN THE SOFTWARE.

`

const HELP = `
  ****************************************
               Citadel Sync
  ****************************************

  Sync files with a specified room on a local or remote Citadel Mail Server.

 - Contacts from vCards ".vcf"
 - Notes from vNotes ".vnt"
 - Calendar from vCalendars ".vcs" or ".ics"
 - Task from vCalendars ".vcs" or ".ics"
 - Text from any text based files ".txt"
`

var (
	name     = flag.String("n", "", "Name of the configuration contacts, notes, calender etc.")
	help     = flag.Bool("h", false, "Prints out this help text")
	username = flag.String("u", "", "[Citadel Username] Requiered if not defined in the config file")
	password = flag.String("p", "", "[Citadel Password] Requiered if not defined in the config file")
	purge    = flag.Bool("D", false, "[Delete] will delete all items in the given room WITHOUT WARNING!")
	version  = flag.Bool("v", false, "Version")
	license  = flag.Bool("l", false, "License")
)

var (
	PS       string = string(os.PathSeparator)
	Cfg      *config
	FILE_DB  string
	FILE_CFG string
	FILE_LOG string
	DB       map[string]*database
)

type config struct {
	Version     int8
	Environment e.Env
	LocalDir    string
	Room        string
	Username    string
	Password    string
	Server      string
	Port        string
	Floor       string
	CIT_SSL_CER string
	citRoomType string
}

func init() {
	Cfg = new(config)
	Cfg.Version = VERSION
	Cfg.Environment = e.ENV_PROD
	Cfg.Username = "TaMeR"
	Cfg.Password = "God knows what"
	Cfg.Server = "localhost"
	Cfg.Port = DEFAULT_PORT
	Cfg.Room = "Contacts"
	Cfg.LocalDir = os.TempDir() + PS + "contacts"
	Cfg.Floor = "Not implemented"
	Cfg.CIT_SSL_CER = "Not implemented"

	DB = make(map[string]*database)
}

func main() {
	flag.Parse()
	switch true {
	case *version:
		fmt.Printf("\n\tCitadel Sync Version %v\n\n", VERSION)
		os.Exit(0)
	case *license:
		fmt.Printf(LICENSE)
		os.Exit(0)
	case *help:
		fmt.Println(HELP)
		flag.PrintDefaults()
		os.Exit(0)
	}
	switch "" {
	case *name:
		fmt.Println(HELP)
		flag.PrintDefaults()
		os.Exit(0)
	}

	FILE_DB = fmt.Sprintf("%s.db.json", *name)
	FILE_CFG = fmt.Sprintf("%s.cfg.json", *name)
	FILE_LOG = fmt.Sprintf("%s.log", *name)

	if err := cfg.Load(FILE_CFG, Cfg); err != nil {
		Cfg.LocalDir = os.Getenv("HOME") + PS + "PIM" + PS + *name
		if err = cfg.Save(FILE_CFG, Cfg); err != nil {
			fmt.Println("cfg.Save error: ", err.Error())
			os.Exit(1)
		} else {
			fmt.Printf("\nPlease edit your config file at:\n\n\t%s\n", FILE_CFG)
			os.Exit(0)
		}
	}

	// Use this if something changes in the config file
	/*
		if Cfg.Version != VERSION {
			fmt.Println("Please check and upgrade your config file version to the new version: ", VERSION)
			os.Exit(0)
		}
	*/
	LogFile()

	if len(*username) != 0 {
		Cfg.Username = *username
	}

	if len(*password) != 0 {
		Cfg.Password = *password
	}

	if *purge {
		CitadelDeleteAll()
		os.Exit(0)
	}

	if !checkRoom() {
		os.Exit(0)
	}
	if Cfg.Server != "" {
		DBLoad()
		FilesInfo()
		CitadelReceive()
		CitadelSend()
		DBSave()
	}
}

func (db *database) Set() {
	DB[db.FileName] = db
}

func (db *database) Del() {
	delete(DB, db.FileName)
}

func (db *database) Modified(fi os.FileInfo) bool {
	if db.FileName != "" {
		if db.FileModTime != fi.ModTime() {
			db.fileModified = true
		}
	} else {
		db.FileModTime = fi.ModTime()
		db.FileName = fi.Name()
		db.fileModified = true
		DB[db.FileName] = db
	}
	return db.fileModified
}

// TODO use linux "file -bip"
func (db *database) mimeType() (ok bool) {
	if db.FileName != "" {
		suffix := Suffix(db.FileName)
		switch suffix {
		case "vcf":
			db.MimeType = "text/vcard"
			ok = true
			if Cfg.citRoomType != citadel.VIEW_ADDRESSBOOK {
				e.Info("WARNING not a vCard Room. Is a: %v", Cfg.citRoomType)
			}
		case "vnt":
			db.MimeType = "text/vnote"
			ok = true
			if Cfg.citRoomType != citadel.VIEW_NOTES {
				e.Info("WARNING not a vNote Room. Is a: %v", Cfg.citRoomType)
			}
		case "vcs", "ics":
			db.MimeType = "text/calendar"
			ok = true
			if Cfg.citRoomType != citadel.VIEW_CALENDAR {
				e.Info("WARNING not a Calender Room. Is a: %v", Cfg.citRoomType)
			}
			if Cfg.citRoomType != citadel.VIEW_TASKS {
				e.Info("WARNING not a Task Room. Is a: %v", Cfg.citRoomType)
			}
		case "txt":
			db.MimeType = "text/text"
			ok = true
			if Cfg.citRoomType != citadel.VIEW_MAILBOX || Cfg.citRoomType != citadel.VIEW_BBS {
				e.Info("WARNING not a text Room. Is a: %v", Cfg.citRoomType)
			}
		default:
			db.MimeType = "text/text"
			e.Info("WARNING Can't determen file type, using text!")
		}
	} else {
		switch Cfg.citRoomType {
		case citadel.VIEW_ADDRESSBOOK:
			db.FileName = db.UID + ".vcf"
			db.MimeType = "text/vcard"
		case citadel.VIEW_NOTES:
			db.FileName = db.UID + ".vnt"
			db.MimeType = "text/vnote"
		case citadel.VIEW_CALENDAR:
			db.FileName = db.UID + ".vcs"
			db.MimeType = "text/calendar"
		case citadel.VIEW_TASKS:
			db.FileName = db.UID + ".vcs"
			db.MimeType = "text/calendar"
		default:
			db.MimeType = "text/text"
			db.FileName = db.UID + ".txt"
		}
		ok = true
	}
	return
}

func (db *database) CitadelEUID() {
	if db.CitEUID == "" {
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "citadel.sync"
		}
		db.CitEUID = fmt.Sprintf("<%v@%v>", db.UID, hostname)
	}
}

func filenameToUID(name string) string {
	var i = strings.LastIndex(name, ".")
	if i > 0 {
		name = name[0:i]
	}
	return name
}

func FilesInfo() {
	fis := readDir(Cfg.LocalDir)
	no := len(fis)
	for i := 0; i < no; i++ {
		if fis[i].IsDir() {
			continue
		}

		uid := filenameToUID(fis[i].Name())
		dbitem, ok := DB[uid]
		if !ok {
			db := new(database)
			db.FileName = fis[i].Name()
			db.CitadelEUID()
			db.mimeType()
			db.CitadelSend()
			os.Remove(Cfg.LocalDir + PS + db.FileName)
			continue
		}
		dbitem.Modified(fis[i])
		dbitem.mimeType()
		dbitem.CitadelEUID()
		DB[uid] = dbitem
	}
}

func readDir(path string) (fis []os.FileInfo) {
	fi, err := os.Stat(path)
	if ok := e.Check(err); !ok {
		err := os.Mkdir(path, 0750)
		if ok := e.Check(err); ok {
			if fi, err = os.Stat(path); err != nil {
				fmt.Println("[ERROR] Something is very wrong, this should never happen ;)")
			}
		} else {
			fmt.Println("[ERROR] Could not make folder")
			os.Exit(0)
		}
	}
	if fi.IsDir() {
		f, err := os.Open(path)
		defer f.Close()
		e.Check(err)
		fis, err = f.Readdir(0) // 0 = All
		e.Check(err)
	}
	return
}

type database struct {
	UID          string
	CitEUID      string
	CitUxTime    int
	FileName     string
	MimeType     string
	FileModTime  time.Time
	fileModified bool
	citModified  bool
}

// Check which one is newer file or citadel
func modified(message *citadel.Message, dav *citadel.Dav) {
	if dbitem, ok := DB[dav.UID]; ok {
		if dbitem.CitUxTime > message.UxTime {
			fmt.Println("CitUx is bigger: ", dav.UID)
			dbitem.citModified = true
			if dbitem.fileModified {
				fmt.Println("fileModified is true ", dav.UID)
				if dbitem.FileModTime.Unix() > int64(message.UxTime) {
					fmt.Println("FileModTime >  ", dav.UID)
					dbitem.citModified = false
				} else {
					dbitem.fileModified = false
				}
			}
		}
	} else {
		dbitem = new(database)
		dbitem.UID = dav.UID
		dbitem.CitUxTime = message.UxTime
		dbitem.mimeType()
		dbitem.citModified = true
		dbitem.fileModified = false
		dbitem.CitadelEUID()
		DB[dav.UID] = dbitem
	}
}

func CitadelReceive() {
	c := citadel.New(Cfg.Server + Cfg.Port)
	defer c.Close()
	c.Login(Cfg.Username, Cfg.Password)
	c.Goto(Cfg.Room)
	db := new(database)
	db.mimeType()
	if list, ok := c.MsgListAll(); ok {
		for _, v := range list {
			if message, dav, ok := c.GetDav(v, db.MimeType); ok {
				modified(message, dav)
				if dbitem, ok := DB[dav.UID]; ok {
					if dbitem.citModified {
						ioutil.WriteFile(Cfg.LocalDir+PS+dbitem.FileName, []byte(dav.Object), 0640)
						dbitem.setDbFileModTime()
					}
				} else {
					// This should never happen
					err := e.New("Item should be in the database!")
					e.Check(err)
				}
			}
		}
	}
}

func (db *database) setDbFileModTime() {
	fi, err := os.Stat(Cfg.LocalDir + PS + db.FileName)
	e.Check(err)
	DB[db.UID].FileModTime = fi.ModTime()
}

func CitadelSend() {
	c := citadel.New(Cfg.Server + Cfg.Port)
	defer c.Close()
	c.Login(Cfg.Username, Cfg.Password)
	c.Goto(Cfg.Room)
	c.Info()

	for _, dbitem := range DB {
		if dbitem.fileModified == false {
			e.Info("File not modified; not sending! %v", dbitem.UID)
			fmt.Printf("Not Sending file %v\n", dbitem.UID)
			continue
		}
		e.Info("Sending file %v", dbitem.UID)
		fmt.Printf("Sending file %v\n", dbitem.UID)
		bytes, err := ioutil.ReadFile(Cfg.LocalDir + PS + dbitem.FileName)
		ok := e.Check(err)
		if ok {
			cmd := fmt.Sprintf("ENT0 1|||4||||||%s", dbitem.CitEUID)
			c.Request(cmd)
			if c.Code == citadel.CODE_SEND_LISTING {
				c.Error = c.Conn.PrintfLine("Content-type: " + dbitem.MimeType)
				c.Check()
				err = c.Conn.PrintfLine("%s", "\n")
				e.Check(err)
				err = c.Conn.PrintfLine("%s", bytes)
				e.Check(err)
				err = c.Conn.PrintfLine("%s", citadel.DE)
				e.Check(err)
				e.Info("%s", bytes)
				dbitem.setDbFileModTime()
			}
		}
	}
}

func (db *database) CitadelSend() {
	c := citadel.New(Cfg.Server + Cfg.Port)
	defer c.Close()
	c.Login(Cfg.Username, Cfg.Password)
	c.Goto(Cfg.Room)
	c.Info()

	bytes, err := ioutil.ReadFile(Cfg.LocalDir + PS + db.FileName)
	ok := e.Check(err)
	if ok {
		cmd := fmt.Sprintf("ENT0 1|||4||||||%s", db.CitEUID)
		c.Request(cmd)
		if c.Code == citadel.CODE_SEND_LISTING {
			c.Error = c.Conn.PrintfLine("Content-type: " + db.MimeType)
			c.Check()
			err = c.Conn.PrintfLine("%s", "\n")
			e.Check(err)
			err = c.Conn.PrintfLine("%s", bytes)
			e.Check(err)
			err = c.Conn.PrintfLine("%s", citadel.DE)
			e.Check(err)
			e.Info("%s", bytes)
		}
	}
}

func CitadelDeleteAll() {
	c := citadel.New(Cfg.Server + Cfg.Port)
	defer c.Close()
	c.Login(Cfg.Username, Cfg.Password)
	c.Goto(Cfg.Room)
	c.Info()
	list, ok := c.MsgListAll()
	if ok {
		c.MsgsDel(list)
	}
}

func checkRoom() (ok bool) {
	c := citadel.New(Cfg.Server + Cfg.Port)
	defer c.Close()
	c.Login(Cfg.Username, Cfg.Password)
	if ok = c.Goto(Cfg.Room); ok {
		Cfg.citRoomType = c.Resp[12]
	} else {
		e.Info("Room not availabe")
		c.Info()
	}
	return
}

func DBLoad() {
	if err := cfg.Load(FILE_DB, &DB); err != nil {
		e.Info("Will create new database at: %v", FILE_DB)
	}
}

func DBSave() {
	if err := cfg.Save(FILE_DB, DB); err != nil {
		e.Info("Could not create database at: ", FILE_DB)
		fmt.Println("Could not create database at: %v", FILE_DB)
	}
}

func LogFile() {
	e.ENVIRONMENT = Cfg.Environment
	if e.ENVIRONMENT != e.ENV_FAIL {
		f, err := os.OpenFile(FILE_LOG, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
		e.Check(err)
		e.Logger(f)
	}
}

func Suffix(filename string) (suffix string) {
	no := strings.LastIndex(filename, ".")
	if no > 0 {
		suffix = filename[no+1:]
	}
	return
}
