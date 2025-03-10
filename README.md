---
title: Citadel Client
description: Library to access Citadel email and collaboration servers from Go, and a command-line tool that uses this library to sync data with the server.
tags: [email,groupware,vCard,vNotes,vCalendar]
---

Citadel Client
==============
[![GoDoc](https://godoc.org/github.com/gotamer/citadel?status.png)](https://godoc.org/github.com/gotamer/citadel)

If you are looking for the citadelsync application go here:
[citadelsync folder](https://github.com/gotamer/citadel/tree/master/citadelsync)


This is a library to access the [Citadel](http://www.citadel.org "Citadel") email and collaboration servers from Go using the [Citadel](http://www.citadel.org "Citadel") Protocol.

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

*See [GoDoc](https://godoc.org/github.com/gotamer/citadel) for complete list*
________________________________________________________

#### Install

##### First install Go

 - [Linux and FreeBSD](http://golang.org/doc/install#tarball)
 - [Mac OS X](http://golang.org/doc/install#osx)
 - [Windows MSI installer](http://golang.org/doc/install#windows)

##### With go installed run

	go get bitbucket.org/gotamer/citadel

