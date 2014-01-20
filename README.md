Go Citadel Client
=================
[![GoDoc](https://godoc.org/bitbucket.org/gotamer/citadel?status.png)](https://godoc.org/bitbucket.org/gotamer/citadel)

This is a library to access [Citadel] email and collaboration servers from Go using the [Citadel] Protocol.
________________________________________________________

#### Features impemented
 - Users
	+ Create
	+ Login
	+ Logout
	+ Change Password

 - Floors
	+ List all floors *with id, name and count of rooms*

 - Rooms
	+ List all rooms
	+ List public rooms

 - Room
	+ Goto room
	+ Stat room

 - Messages
    + Read Message
	+ Parse vCard, vNotes and vCalendar files

*See [GoDoc] for complete list*
________________________________________________________

#### Install

##### First install Go

 - [Linux and FreeBSD](http://golang.org/doc/install#tarball)
 - [Mac OS X](http://golang.org/doc/install#osx)
 - [Windows MSI installer](http://golang.org/doc/install#windows)

##### With go installed run

	go get bitbucket.org/gotamer/citadel


________________________________________________________

#### The MIT License (MIT)

Copyright © 2013 Dennis T Kaplan <http://www.robotamer.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.


[Citadel]:(http://www.citadel.org "Citadel")
[GoDoc]:(https://godoc.org/bitbucket.org/gotamer/citadel)
