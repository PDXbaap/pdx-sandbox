package main

import (
	"os/exec"
	"bytes"
	"syscall"
	"strconv"
	"log"
)

func execute(args []string) (exitcode string, output string) {

	exitcode = "0"

	exe := exec.Command(args[0], args[1:]...)

	var out bytes.Buffer
	exe.Stdout = &out
	exe.Stderr = &out

	err := exe.Run()

	if err != nil {

		log.Println("cmd exec error: ", err.Error(), out.String())

		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				log.Printf("Exit code: %d", status.ExitStatus())
				exitcode = strconv.Itoa(status.ExitStatus())
			}
		}
	}

	output = out.String()

	return
}