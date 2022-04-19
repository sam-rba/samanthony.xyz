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
