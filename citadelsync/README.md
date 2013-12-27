[Citadel] Sync
============

***************************************************
## Import vCards to [Citadel]
***************************************************

A config file is requiered, set it with the -c flag.

If the specified config file does not exist, one
will be created with default values.

 > -D will delete all items in the given room WITHOUT WARNING

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

**Hint**: Don't use the default config file name if you
are planing to have more then one configuration.


	-v=false: version
	-h=false: Prints out this help text
	-c="citadelVcard.json": Config file (*.json)
	-u="": Username
	-p="": Password
	-r=false: Check if room exists
	-D=false: Delete all items in the room!
	-i="": Import file (*.vcf)

## Install

### Executable

 > There is an exectable version for Linux at:
https://bitbucket.org/gotamer/citadel/downloads/citadelsync

### From Source

 > Install go then run

	go get bitbucket.org/gotamer/citadel


________________________________________________________

#### The MIT License (MIT)

Copyright © 2013 Dennis T Kaplan <http://www.robotamer.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.


[Citadel]:(http://www.citadel.org "Citadel")
[GoDoc]:(https://godoc.org/bitbucket.org/gotamer/citadel)
