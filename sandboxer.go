// Copyright 2017 The PDX Blockchain Hybercloud Authors
// This file is part of the PDX chainmux implementation.
//
// The PDX Blcockchain Hypercloud is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PDX Blockchain Hypercloud is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the software. If not, see <http://www.gnu.org/licenses/>.


// PDX sandboxer, a setgid docker helper for PDX smart-contract sandboxing.

package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"
	"log"
	"os/exec"
	"strings"
	"time"
)



func main() {

	flag.Usage = func() {
		fmt.Println("")
		fmt.Println("PDX sandboxer, a hardened setgid docker helper for PDX smart-contract sandboxing")
		fmt.Println("")
		fmt.Println("Usage [after privilege elevation]: ./sandboxer docker [OPTIONS] COMMAND [ARG...]")
		fmt.Println("")
		fmt.Println("For example,")
		fmt.Println("")
		fmt.Println("	sudo ./sandboxer docker run -it -v $PDX_HOME/dapps:/dapps/ pdx-dapp-omni /bin/sh")
		fmt.Println("")
                fmt.Println("Please visit https://github.com/PDXbaap/pdx-sandboxer to get the latest version.")
		fmt.Println("")
	}

	if len(os.Args) < 2 || os.Args[1] == "-h" {
		flag.Usage();
		return
	}

	binary, err := exec.LookPath(os.Args[1])
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	log.Println(os.Args)

	var lockfile *os.File

	for {
		lockfile, err = os.OpenFile(os.Getenv("PDX_HOME")+"/temp/sandboxer.lock", os.O_CREATE, os.ModePerm)

		err = syscall.Flock(int(lockfile.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)

		if err == syscall.EWOULDBLOCK {
			log.Println("another instance running, waiting")
			time.Sleep(100 * time.Millisecond)
			continue
		}

		break
	}

	defer func() {
		syscall.Flock(int(lockfile.Fd()), syscall.LOCK_UN)
		lockfile.Close()
	}()

	log.Println("now only myself is running")

	time.Sleep(10000 * time.Millisecond)

	housekeeping()

	if !accessControl(os.Args[1:]) {
		log.Fatal("unauthorized priviledged access, exiting ...")
	}

	if err := syscall.Exec(binary, os.Args[1:], os.Environ()); err != nil {
		log.Fatal(err)
		panic(err)
	}

	fmt.Println("goodbye")
}


func accessControl(args []string) bool {

	// IMPORTANT: sandboxer whitelist rules
	//
	// 1) Only allow docker run/stop
	//
	// 2) docker run [OPTIONS] IMAGE [COMMAND] [ARG...]
	//
	//		a) unpriviledged only
	//
	//		b)
	//
	//		 ) record the container id
	//
	// 3) docker stop [OPTIONS] CONTAINER [CONTAINER...]
	//
	// 		only containers started by sandboxer
	//
	// 4) docker stats [OPTIONS] [CONTAINER...]
	//
	//		only containers started by sandboxer

	if args[0] != "docker" {
		log.Println("not a docker binary")
		return false
	}

	if args[1] != "run" && args[1] != "kill" && args[1] != "stop" {
		log.Println("not docker run/kill/stop")
		return false
	}

	var cmd = strings.Join(args," ");

	log.Println(cmd);

	return false
}

func housekeeping() {
	// remove dead container ids off the record
}