// Copyright 2018 The PDX Blockchain Hybercloud Authors
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


// PDX sandbox, a secure docker service supporting signed images and fine-grained access control 

package main

import (
	"fmt"
	"log"
	"net/http"
	"net"
	"os"
	"flag"
)

var pdxHome string
var imgHome string
var crtFile string
var imgVerifierScript string
var lockfile string
var datafile string

var  tcpAddr string

func main() {

	flag.Usage = func() {
		fmt.Println("")
		fmt.Println("PDX sandbox,  a secure privileged service for hardcore docker sandboxing")
		fmt.Println("")
		fmt.Println("Please visit https://github.com/PDXbaap/pdx-sandbox for more information")
		fmt.Println("")
		fmt.Println("Use with elevated privileged with PDX iaas-compute installed at $PDX_HOME")
		fmt.Println("")

		flag.PrintDefaults()
	}

	flag.StringVar(&tcpAddr, "addr",  "127.0.0.1:0",        "TCP host:port to listen on")
	flag.StringVar(&pdxHome, "home", os.Getenv("PDX_HOME"), "PDX iaas-compute directory")

	flag.Parse()

	if pdxHome == "" {
		flag.Usage()
		os.Exit(1)
	}

	// set files & directories needed for it to function

	imgHome = pdxHome + "/bin/image"
	crtFile = pdxHome + "/conf/signer.crt"
	imgVerifierScript = pdxHome + "/bin/img-verify.sh"
	lockfile = pdxHome + "/temp/sandbox.lock"
	datafile = pdxHome + "/temp/sandbox.data"

	lock()

	defer unlock()

	listener, err := net.Listen("tcp", tcpAddr)

	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("listening on: %s", listener.Addr().String())

	save(listener.Addr().String())

	http.HandleFunc("/", handler)

	log.Fatalln(http.Serve(listener, nil))
}

