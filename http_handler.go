package main

import (
	"io/ioutil"
	"strings"
	"fmt"
	"net/http"
	"log"
)

func handler(w http.ResponseWriter, r *http.Request) {

	var cmd string

	if r.Method == http.MethodGet {
		cmd = r.URL.Query().Get("cmd")
	} else if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			log.Println("invalid request")
			http.Error(w, "invalid request", http.StatusUnauthorized)
			return
		}
		cmd = string(body)
	} else {
		log.Println("unsupported http method: ", r.Method)
		http.Error(w, "unsupported http method: " + r.Method, http.StatusUnauthorized)
		return
	}

	log.Println("received cmd: " + cmd)

	args := strings.Fields(cmd)

	ok, reason, name := authorized(args);

	if !ok {
		log.Println("unauthorized:", reason)
		http.Error(w, "unauthorized: " + reason, http.StatusUnauthorized)
		return
	}

	if args[1] == "run" {

		if name == "" {
			log.Println("missing container name, noop")
			http.Error(w, "missing container name", http.StatusBadRequest)
			return
		}

		saveStartedContainers(name)
		log.Println("starting container: ", name)

	} else if args[1] == "load" {

		args[3] = imgHome + "/" + args[3]

	} else {
		// noop
	}

	exitcode, output := execute(args)

	w.Header().Set("DOCKER_EXIT_CODE", exitcode)

	if exitcode != "0" {
		http.Error(w, output, http.StatusBadRequest)
	} else {
		log.Println("started container:", name)
		fmt.Println(w, output)
	}

	return
}
