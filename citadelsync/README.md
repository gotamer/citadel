[CitadelSync]
==============


[![GoDoc](https://godoc.org/bitbucket.org/gotamer/citadel/citadelsync?status.png)](https://godoc.org/bitbucket.org/gotamer/citadel/citadelsync)


************************************************
## Sync local files with a [Citadel] Mail Server
************************************************

[CitadelSync] can sync files with a specified room on a local or remote [Citadel] Mail Server.

You may sync any type of text files but Citadel Sync is most useful for

 - Contacts from vCards `.vcf`
 - Notes from vNotes `.vnt`
 - Calendar `.vcs` or `.ics`
 - Task `.vcs` or `.ics`
 - Text from any text based files `.txt`

____________________________________________________
## Install

### From Executable

This part is compatible with any Linux AMD64 based system
For other systems please see install section below

There is an executable version for Linux AMD64 at [download]

	cd
	mkdir citsync
	cd citsync
	wget https://bitbucket.org/gotamer/citadel/downloads/citadelsync



____________________________________________________

### Install From Go Source

##### First install Go

 - [Linux and FreeBSD](http://golang.org/doc/install#tarball)
 - [Mac OS X](http://golang.org/doc/install#osx)
 - [Windows MSI installer](http://golang.org/doc/install#windows)

##### With go installed run

	go get bitbucket.org/gotamer/citadel
	cd citadelsync
	go install



____________________________________________________

## HowTo use

### Create a config file

	citadelsync -n contacts
	nano contacts.cfg.json

### Edit the config file `contacts.cfg.json`
```
{
	"Version": 3,
	"Environment": 1,
	"LocalDir": "/home/username/PIM/contacts",
	"Room": "Contacts",
	"Username": "TaMeR",
	"Password": "God knows what",
	"Server": "localhost",
	"Port": ":504",
	"Floor": "Not implemented",
	"SSL_KEY": "Not yet implemented",
	"SSL_CER": "Not yet implemented"
}
```

#### Version:
Do not change unless prompted after an upgrade.

#### Environment:

	1. Production
	2. Info mode, prints a lot of info in to the log file
	3. Debug mode, will print to screen, and exit if it finds something not quite right

#### LocalDir:
Point to the local folder containing your vCards, vNotes etc. files.
You should make a folder for each type, since each folder will represent a room in the [Citadel] server.

 - Citadel `Contacts` to the `contacts` folder
 - Citadel `Tasks` to the tasks` folder
 - etc.

#### Room:
The Citadel Room to upload to

#### Username and Password
Your Citadel username or password.
Keep empty "" to specify on the command line.

#### Server:
Your Citadel hostname, such as `example.com`

#### Port:
Your Citadel port. The standard port is 504 if you haven't changed it on the server.
____________________________________________________
### Populate your sync directory

Your sync directory is:

	"LocalDir": "/home/username/PIM/contacts"

In this example we would place our vCard files in to this directory

*Files must be individual files each holding a single entry

____________________________________________________
### Copy to Citadel

Now the setup is done!

Following command will copy all files to the Citadel server.

	citadelsync -n contacts

This will

 1. Check your server connection
 2. Your login information
 3.	The availability of the room
 4. The compatibillity of the room with your file type
 5. Your `LocalDir` sync directory

If everything checks out it will upload all files from your LocalDir to the specified Room on the local or remote Citadel Server.
____________________________________________________
### Notes

 * Files must be individual files each holding a single entry residing in one folder, anywhere on your computer or device.

 * Citadel Sync remembers the state of your files and will only upload modified files

 * Citadel Sync will not delete or modify any files it didn't upload, unless you use the -D flag.

 * A file with the extension `*.db.json` will be created. Please do not modify this file by hand!!!

____________________________________________________

The MIT License (MIT)
=====================

Copyright © 2013 Dennis T Kaplan <http://www.robotamer.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sub-license, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NON-INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.


[Citadel]:(http://www.citadel.org)
[CitadelSync]:(http://bitbucket.org/gotamer/citadel/wiki)
[download]:(https://bitbucket.org/gotamer/citadel/downloads/citadelsync)
[GoDoc]:(https://godoc.org/bitbucket.org/gotamer/citadel/citadelsync)

