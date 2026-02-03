+++
title = 'Tools I Use'
date = '2026-01-02T15:24:53-03:30'
tags = []
draft = false
+++

Here is a list of some of the tools I use day-to-day.


# Software

- [Gentoo Linux](https://www.gentoo.org/) ---
	One of the most solid barebones Linux distros I know of.
	I use it on my desktop and laptop computers.
- [OpenBSD](https://www.openbsd.org/) ---
	Secure server operating system.
	I use it to serve this website.

	Unlike a GNU/Linux distro, OpenBSD is a single homogeneous system.
	Everything makes sense and the documentation is really good.
	Reminds me of [Plan 9](https://plan9.io/plan9/) in that respect.

- [Acme](http://acme.cat-v.org/) ---
	Originally developed for Plan 9, Acme is [Rob Pike's](https://commandcenter.blogspot.com/) user interface for programmers.
	It's what we now call an IDE.
	I'm writing this article in Acme right now.

	Every piece of text in Acme can be a executed or piped into/out-of a script.
	Very powerful.
- [Dwm](https://dwm.suckless.org/) ---
	Nice tiling window manager for X11.
	Goes well with other Suckless accoutrements.
- [Syncthing](https://syncthing.net/) ---
	Syncthing is a sort of distributed filesystem.
	I use it to synchronize files between my laptop, desktop, and phone.


# Hardware

- [Thinkpad T420s](https://www.thinkwiki.org/wiki/Category:T420s) ---
	I've had this laptop for a few years now; no complaints.

	I replaced the hard drive with an SSD and threw some extra DDR3 in there.
	The 10+ year old 4-thread Sandy Bridge i5 is actually fine.
	I'm waiting for someone to make a [serious multithreaded](https://netlib.org/utk/papers/advanced-computers/tera.html) RISC-V CPU, but unfortunately everyone seems to be obsessed with high clock speeds and out-of-order-execution chips that use as much die space and power as possible.
	That [Esperanto ET-SoC-1](https://youtu.be/LmUu-lN7D0k) looked promising, but apparently they went out of business or something?

	Anyway, the T420s has a sturdy magnesium frame, a good keyboard, and a three-button touchpad which is essential for Acme and CAD programs---I don't know how people live without one.

- [USBtin](https://www.fischl.de/usbtin/) ---
	Simple USB-to-CAN interface by Thomas Fischl.
	Works with [SocketCAN](https://docs.kernel.org/networking/can.html).
	Used to test and debug systems that incorporate a CAN bus, like [this](http://git.samanthony.xyz/can-gauge-interface.git/).

- [Espotek Labrador](https://espotek.com/labrador/) ---
	Combined oscilloscope, signal generator, power supply, logic analyzer, and multimeter.
	Obviously a tiny $30 board is not as good as real lab equipment, but it's small and cheap and good enough for now.
