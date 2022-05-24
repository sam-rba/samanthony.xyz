/*
Copyright 2022 Sam Anthony

This file is part of samanthony.xyz.

samanthony.xyz is free software: you can redistribute it and/or modify it under
the terms of the GNU General Public License as published by the Free Software
Foundation, either version 3 of the License, or (at your option) any later
version.

samanthony.xyz is distributed in the hope that it will be useful, but WITHOUT
ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with
samanthony.xyz. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/sys/unix"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func uidOf(user string) (int, error) {
	passwdFile, err := os.Open("/etc/passwd")
	if err != nil {
		return -1, err
	}
	defer passwdFile.Close()

	scanner := bufio.NewScanner(passwdFile)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()

		parsed := strings.Split(line, ":")

		name := parsed[0]

		if name == user {
			uid, err := strconv.Atoi(parsed[2])
			if err != nil {
				return -1, err
			}
			return uid, nil
		}
	}
	return -1, errors.New(fmt.Sprintf("user '%s' not in /etc/passwd", user))
}

func gidOf(group string) (int, error) {
	groupFile, err := os.Open("/etc/group")
	if err != nil {
		return -1, err
	}
	defer groupFile.Close()

	scanner := bufio.NewScanner(groupFile)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()

		parsed := strings.Split(line, ":")

		name := parsed[0]

		if name == group {
			gid, err := strconv.Atoi(parsed[2])
			if err != nil {
				return -1, err
			}
			return gid, nil
		}
	}
	return -1, errors.New(fmt.Sprintf("group '%s' not in /etc/group", group))
}

func dropPerms(uid, gid int) error {
	if runtime.GOOS != "linux" {
		if err := unix.Setgid(gid); err != nil {
			return errors.New(fmt.Sprintf("setgid(%d): %v", gid, err))
		}
		if err := unix.Setuid(uid); err != nil {
			return errors.New(fmt.Sprintf("setuid(%d): %v", uid, err))
		}
		return nil
	} else {
		// setuid/setgid has supposedly been fully supported on Linux
		// since go 1.16 but I can't seem to get it to work properly.
		log.Print("setgid not supported on Linux, skipping.")
		log.Print("setuid not supported on Linux, skipping.")
		return nil
	}
}
