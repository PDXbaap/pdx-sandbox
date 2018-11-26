package main

import (
	"os"
	"syscall"
	"time"
	"log"
)

var lockF *os.File

func lock() {

	var err error

	for {
		lockF, err = os.OpenFile(lockfile, os.O_CREATE|os.O_WRONLY, os.ModePerm)

		err = syscall.Flock(int(lockF.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)

		if err == syscall.EWOULDBLOCK {
			log.Println("another instance running, waiting")
			time.Sleep(100 * time.Millisecond)
			continue
		}

		break
	}
}

func unlock() {
	syscall.Flock(int(lockF.Fd()), syscall.LOCK_UN)
	lockF.Close()
}

func save(data string) {
	lockF.WriteString(data)
}