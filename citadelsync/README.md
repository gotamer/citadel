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

Download the executable for your system from   
https://bitbucket.org/gotamer/citadel/downloads

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

### Create a config file.
Following will create a template for the *contacts* room

	citadelsync -n contacts

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
 3. The availability of the room
 4. The compatibillity of the room with the local file type
 5. Your `LocalDir` sync directory

If everything checks out it will upload all files from your LocalDir to the specified Room on the local or remote Citadel Server.
____________________________________________________
### Notes

 * Files must be individual files each holding a single entry residing in one folder, anywhere on your computer or device.

 * Citadel Sync remembers the state of your files and will only upload modified files

 * Citadel Sync will not delete or modify any files it didn't upload, unless you use the -D flag.

 * A file with the extension `*.db.json` will be created. Please do not modify this file by hand!!! 

____________________________________________________

[Citadel]:(http://www.citadel.org)
[CitadelSync]:(http://bitbucket.org/gotamer/citadel/wiki)
[download]:(https://bitbucket.org/gotamer/citadel/downloads)
[GoDoc]:(https://godoc.org/bitbucket.org/gotamer/citadel/citadelsync)
