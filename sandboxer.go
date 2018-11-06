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
	"bufio"
)

var lockfile = os.Getenv("PDX_HOME")+"/temp/sandboxer.lock"
var datafile = os.Getenv("PDX_HOME")+"/temp/sandboxer.data"

var startedContainers  = make(map[string]string)

func main() {

	flag.Usage = func() {
		fmt.Println("")
		fmt.Println("PDX sandboxer, a hardened setgid docker helper for PDX smart-contract sandboxing")
		fmt.Println("")
		fmt.Println("Usage [after privilege elevation via sudo or setgid]:")
		fmt.Println("")
		fmt.Println("	./sandboxer docker run [OPTIONS] IMAGE [COMMAND] [ARG...]")
		fmt.Println("	./sandboxer docker stop [OPTIONS] CONTAINER [CONTAINER...]")
		fmt.Println("	./sandboxer docker stats [OPTIONS] [CONTAINER...]")
		fmt.Println("")
		fmt.Println("Note: Option name and argument (if present) MUST be ONE token")
		fmt.Println("")
		fmt.Println("For example,")
		fmt.Println("")
		fmt.Println("	./sandboxer docker run -it -v=$PDX_HOME/dapps:/dapps/ pdx-dapp-omni /bin/sh")
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

	var lockF *os.File

	for {
		lockF, err = os.OpenFile(lockfile, os.O_CREATE, os.ModePerm)

		err = syscall.Flock(int(lockF.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)

		if err == syscall.EWOULDBLOCK {
			log.Println("another instance running, waiting")
			time.Sleep(100 * time.Millisecond)
			continue
		}

		break
	}

	defer func() {
		syscall.Flock(int(lockF.Fd()), syscall.LOCK_UN)
		lockF.Close()
	}()

	log.Println("now only myself is running")

	if !accessControl(os.Args[1:]) {
		log.Fatal("unauthorized priviledged access, exiting ...")
	}

	loadStarted()

	housekeeping()

	// get name of container to be created

	var name string = ""

	for i, v := range os.Args[1:] {
		if v == "--name" {
			name = os.Args[i + 2]
			break
		}
	}

	if name == "" {
		log.Println("missing container name, existing ...")
		os.Exit(-1)
	}

	startedContainers[name] = name

	saveStarted()

	log.Println("starting container: ", name)

	if err := syscall.Exec(binary, os.Args[1:], os.Environ()); err != nil {
		log.Fatal(err)
		panic(err)
	}

	fmt.Println("goodbye")
}


func accessControl(args []string) bool {

	////////////////////////////////////////////////////////
	//
	// IMPORTANT: sandboxer whitelist rules
	//
	// 1) Only allow docker run/stop/stats
	//
	// 2) docker run [OPTIONS] IMAGE [COMMAND] [ARG...]
	//
	//		unprivileged no-harm options only
	//
	// 3) docker stop [OPTIONS] CONTAINER [CONTAINER...]
	//
	// 		only containers started by sandboxer
	//
	// 4) docker stats [OPTIONS] [CONTAINER...]
	//
	//		only containers started by sandboxer
	//
    //////////////////////////////////////////////////////

	// only do docker, nothing else

	if args[0] != "docker" {
		log.Println("not a docker binary")
		return false
	}

	// only stats what we have started

	if args[1] == "stats" {

		for _, v := range args[2:] {

			if strings.HasPrefix(v, "-") {
				continue
			}

			if _, ok := startedContainers[v]; !ok {
				return false
			}
		}

		return true
	}

	// only stop what we have started

	if args[1] == "stop" {

		for _, v := range args[2:] {

			if strings.HasPrefix(v, "-") {
				continue
			}

			if _, ok := startedContainers[v]; !ok {
				return false
			}
		}

		return true
	}

	if args[1] != "run" {
		log.Println("not docker run/stop/stats")
		return false
	}

	// check docker run [OPTIONS] IMAGE [COMMAND] [ARG...]

	for i := 2; i < len(args) ; i++ {

		v := args[i]

		if strings.HasPrefix(v,"--privileged") {
			if !strings.Contains(v,"=false") {
				return false
			}
		}

		if strings.HasPrefix(v, "--cap-add") {
			return false
		}

		if strings.HasPrefix(v,"--device") {
			return false
		}

		if strings.HasPrefix(v, "--group-add") {
			return false
		}

		if strings.HasPrefix(v,"--ipc") {
			if strings.Contains(v, "host") || strings.Contains(v, "shareable") ||
				strings.Contains(v, "container:") {
				return false
			}
		}

		if strings.HasPrefix(v, "--security-opt")  {
			if !strings.Contains(v, "no-new-privileges") {
				return false
			}
		}

		if strings.HasPrefix(v, "-v") || strings.HasPrefix(v, "--volume") {
			if !strings.Contains(v,"ro") {
				return false;
			}
		}

		if !strings.HasPrefix(v, "-") { //docker image now
			if strings.HasPrefix(v, "pdxbaap/pdx-sandbox") || strings.HasPrefix(v, "pdx-sandbox") ||
				strings.HasPrefix(v, "pdxbaap/pdx-dapp-omni") || strings.HasPrefix(v, "pdx-dapp-omni")  ||
				strings.HasPrefix(v, "pdxbaap/pdx-dapp-exec") || strings.HasPrefix(v, "pdx-dapp-exec") ||
				strings.HasPrefix(v, "pdxbaap/pdx-dapp-java") || strings.HasPrefix(v, "pdx-dapp-java") {
				return true
			}
		}
	}

	return false
}

func housekeeping() {
	// remove dead container ids off the record
}

func loadStarted() {

	file, err := os.OpenFile(datafile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Println("cannot read data file: ", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		text := scanner.Text()
		startedContainers[text] = text
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}

func saveStarted() {

	os.Remove(datafile)

	file, err := os.Create(datafile)
	if err != nil {
		log.Println("cannot write data file: ", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, v := range startedContainers {
		fmt.Fprintln(writer, v)
	}

	writer.Flush()
}