// Import vCards to the Citadel Mail Server
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"bitbucket.org/gotamer/cfg"
	"bitbucket.org/gotamer/citadel"
	"bitbucket.org/gotamer/errors"
	"bitbucket.org/gotamer/tools"
	"bitbucket.org/llg/vcard"
)

const (
	VERSION      = 1
	VCARD_TIME   = "20060102T150405Z0700"
	DEFAULT_PORT = ":504"
)

var (
	cfg_file    = flag.String("c", "citadelVcard.json", "Config file (*.json)")
	helpflag    = flag.Bool("h", false, "Prints out this help text")
	checkroom   = flag.Bool("r", false, "Check if room exists")
	username    = flag.String("u", "", "Username")
	password    = flag.String("p", "", "Password")
	purge       = flag.Bool("D", false, "Delete all items in the room!")
	version     = flag.Bool("v", false, "version")
	import_file = flag.String("i", "", "Import file (*.vcf)")
)

var help = `
  ***************************************************
       Import vCards to Citadel
  ***************************************************
  A config file is requiered, set it with the -c flag.

  If the specified config file does not exist, one
  will be created with default values.

  -D will delete all items in the given room WITHOUT WARNING

  Optionaly you may specify an import path, to import vcards.
  The import path may be a file, or a folder.

  Username and Password for the Citadel Mail Server may be
  defined in the config file, or optionaly on the command line

  Enviroment:
	0 = Production
	1 = Prints alot of info
	2 = Debug mode, same as 1 but will exit on error

  The -r Flag checks if the room exists one mail server. You
  can use this to verify that you have spelled the room name correctly

  Hint: Don't use the default config file name if you
  are planing to have more then one configuration.
`

var (
	PS          string = string(os.PathSeparator)
	addressBook vcard.AddressBook
	Cfg         *config
)

type config struct {
	Version    int8
	Enviroment e.Env
	PathVcard  string
	Room       string
	Username   string
	Password   string
	Server     string
	Port       string
	LogFile    string
	Floor      string
	SSL_KEY    string
	SSL_CER    string
}

func init() {
	Cfg = new(config)
	Cfg.Version = VERSION
	Cfg.Enviroment = e.ENV_PROD
	Cfg.Username = "TaMeR"
	Cfg.Password = "God knows what"
	Cfg.Server = "localhost"
	Cfg.Port = DEFAULT_PORT
	Cfg.Room = "Contacts"
	Cfg.PathVcard = "/tmp/vcard"
	Cfg.Floor = "Not implemented"
	Cfg.SSL_CER = "Not implemented"
	Cfg.SSL_KEY = "Not implemented"
	Cfg.LogFile = "citadelVcard.log"
}

func LogFile() {
	e.ENVIROMENT = Cfg.Enviroment
	if Cfg.LogFile != "" {
		f, err := os.OpenFile(Cfg.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
		e.Check(err)
		e.Logger(f)
	}
}

func main() {
	flag.Parse()
	if *version {
		fmt.Printf("\n\tCitadel Import Version %v\n\n", VERSION)
		os.Exit(0)
	}
	if *helpflag {
		fmt.Println(help)
		flag.PrintDefaults()
		fmt.Println("\n")
		os.Exit(0)
	}
	if err := cfg.Load(*cfg_file, Cfg); err != nil {
		if err = cfg.Save(*cfg_file, Cfg); err != nil {
			fmt.Println("cfg.Save error: ", err.Error())
			os.Exit(1)
		} else {
			fmt.Printf("\nPlease edit your config file at:\n\n\t%s\n", *cfg_file)
			os.Exit(0)
		}
	}

	if Cfg.Version != VERSION {
		fmt.Println("Please check and upgrade your config file version to the new version: ", VERSION)
		os.Exit(0)
	}

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

	if *checkroom {
		checkRoom()
		os.Exit(0)
	}

	VcardsLoad()
	if *import_file != "" {
		VcardsImport(*import_file)
		VcardsCheck()
	}
	VcardsSave()
	CitadelSend()
}

func CitadelSend() {

	c := citadel.New(Cfg.Server + Cfg.Port)
	defer c.Close()
	c.Login(Cfg.Username, Cfg.Password)
	c.Goto(Cfg.Room)
	c.Info()
	var err error

	fis := readDir(Cfg.PathVcard)
	var no = len(fis)
	for i := 0; i < no; i++ {
		fi := fis[i]
		card := VcardLoad(fi)
		if card == nil {
			continue
		}
		contact := VcardString(card)

		c.Request("ENT0 1|||4")
		if c.Code == citadel.CODE_SEND_LISTING {
			c.Error = c.Conn.PrintfLine("Content-type: text/x-vcard; charset=UTF-8")
			c.Check()
			err = c.Conn.PrintfLine("%s", "\n")
			e.Check(err)
			err = c.Conn.PrintfLine("%s", contact)
			e.Check(err)
			err = c.Conn.PrintfLine("%s", citadel.DE)
			e.Check(err)
			e.Info("%s", contact)
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
		ok = c.MsgsDel(list)
		if ok {
			list, _ := c.MsgListAll()
			fmt.Printf("list: %s\n\n", list)
		}
	}
}

func VcardString(address *vcard.AddressBook) string {
	var b = new(bytes.Buffer)
	writer := vcard.NewDirectoryInfoWriter(b)
	address.WriteTo(writer)
	return b.String()
}

func VcardLoad(fi os.FileInfo) (contact *vcard.AddressBook) {
	if fi.IsDir() {
		e.LogError.Printf("Path must be a vcatd file, this is a folder: %s\n", fi.Name())
		return
	}
	if fi.Size() < 200 {
		e.LogError.Printf("File is too small for a vcard: %s\n", fi.Name())
		return
	}
	if !strings.HasSuffix(fi.Name(), ".vcf") {
		e.LogError.Printf("File must be a vcard ending with .vcf Got: %s\n", fi.Name())
		return
	}
	file := fmt.Sprintf("%s%s%s", Cfg.PathVcard, PS, fi.Name())
	f, err := os.Open(file)
	if err != nil {
		e.LogError.Printf("Can't read file %s\n", file)
		return
	}
	reader := vcard.NewDirectoryInfoReader(f)
	contact = new(vcard.AddressBook)
	contact.ReadFrom(reader)
	f.Close()
	return contact
}

func VcardsLoad() {
	fis := readDir(Cfg.PathVcard)
	no := len(fis)
	for i := 0; i < no; i++ {
		if fis[i].IsDir() {
			continue
		}
		if fis[i].Size() < 200 {
			continue
		}
		file := fmt.Sprintf("%s%s%s", Cfg.PathVcard, PS, fis[i].Name())
		if !strings.HasSuffix(file, ".vcf") {
			continue
		}
		f, err := os.Open(file)
		if err != nil {
			e.LogError.Printf("Can't read file %s\n", file)
			continue
		}
		reader := vcard.NewDirectoryInfoReader(f)
		addressBook.ReadFrom(reader)
		f.Close()
	}
}

func readDir(path string) (fis []os.FileInfo) {
	fi, err := os.Stat(path)
	e.Check(err)
	if fi.IsDir() {
		f, err := os.Open(path)
		defer f.Close()
		e.Check(err)
		fis, err = f.Readdir(0) // 0 = All
		e.Check(err)
	}
	return
}

func VcardsImport(path string) {
	fi, err := os.Stat(path)
	e.Check(err)
	if fi.IsDir() {
		f, err := os.Open(path)
		if err != nil {
			e.LogError.Printf("Can't read file %s\n", path)
			return
		}
		fil, err := f.Readdir(0) // 0 = All
		f.Close()
		no := len(fil)
		for i := 0; i < no; i++ {
			if fil[i].IsDir() {
				continue
			}
			if fil[i].Size() < 200 {
				continue
			}
			file := fmt.Sprintf("%s%s%s", path, PS, fil[i].Name())
			if !strings.HasSuffix(file, ".vcf") {
				continue
			}
			f, err := os.Open(file)
			if err != nil {
				e.LogError.Printf("Can't read file %s\n", file)
				continue
			}
			reader := vcard.NewDirectoryInfoReader(f)
			addressBook.ReadFrom(reader)
			f.Close()
		}
	} else {
		if fi.Size() < 200 {
			e.Info("File too small %s", fi.Size())
			return
		}
		if !strings.HasSuffix(path, ".vcf") {
			e.Info("Wrong file type")
			return
		}
		f, err := os.Open(path)
		e.Check(err)
		defer f.Close()
		reader := vcard.NewDirectoryInfoReader(f)
		addressBook.ReadFrom(reader)
	}

}

func VcardsSave() {
	for _, c := range addressBook.Contacts {
		var address = new(vcard.AddressBook)
		address.Contacts = append(address.Contacts, c)
		uid := address.Contacts[0].UID
		path := fmt.Sprintf("%s%s%s.vcf", Cfg.PathVcard, PS, uid)
		file, err := os.Create(path)
		e.Check(err)
		defer file.Close()
		bufoutput := bufio.NewWriter(file)
		var output io.Writer
		output = bufoutput
		defer bufoutput.Flush()
		writer := vcard.NewDirectoryInfoWriter(output)
		address.WriteTo(writer)
	}
}

func VcardsCheck() {
	for i := 0; i < len(addressBook.Contacts); i++ {
		contact := &(addressBook.Contacts[i])
		var modified bool

		origianal := contact

		if len(contact.REV) != 16 {
			e.Info("Modified: REV should be 16 but is %s", contact.REV)
			modified = true
		}

		if len(contact.UID) < 2 {
			e.Info("Modified: UID!")
			contact.UID = tools.Uid16()
			modified = true
		}

		if modified {
			e.Info("Org: %s", origianal)
			contact.REV = time.Now().UTC().Format(VCARD_TIME)
			e.Info("New: %s", contact)
		}
	}
}

func fixFamilyName(contact *vcard.VCard) {
	if len(contact.FamilyNames) == 0 {
		if len(contact.FormattedName) > 1 {
			fn := strings.Split(contact.FormattedName, " ")
			no := len(fn)
			contact.FamilyNames = []string{fn[no-1]}
		}
	}
}

func fixGivenName(contact *vcard.VCard) {
	if len(contact.GivenNames) == 0 {
		if len(contact.FormattedName) > 1 {
			fn := strings.Split(contact.FormattedName, " ")
			contact.GivenNames = []string{fn[0]}
		}
	}
}

func checkRoom() {
	c := citadel.New(Cfg.Server + Cfg.Port)
	defer c.Close()
	c.Login(Cfg.Username, Cfg.Password)
	c.Goto(Cfg.Room)
	list, _ := c.MsgListAll()
	fmt.Printf("\nList id of contacts in room:\n%s\n\n", list)
}
