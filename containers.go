package main

import (
	"os"
	"bufio"
	"fmt"
	"log"
	"os/exec"
)

var startedContainers  = make(map[string]string)

func init() {
	loadStartedContainers()
}

func loadStartedContainers() {

	file, err := os.OpenFile(datafile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		text := scanner.Text()
		log.Println("found started container: " + text)
		startedContainers[text] = text
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}

func saveStartedContainers(name string) {

	removeDeadContainers()

	startedContainers[name] = name

	os.Remove(datafile)

	file, err := os.Create(datafile)
	if err != nil {
		log.Println("cannot write data file: ", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, v := range startedContainers {
		log.Println("save started container:" + v)
		fmt.Fprintln(writer, v)
	}

	writer.Flush()
}

func removeDeadContainers() {

	// remove dead container ids off the record

	for _, v := range startedContainers {

		//docker inspect -f '{{.State.Running}}' v

		out, err := exec.Command("docker", "-f", "{{.State.Running}}", v).CombinedOutput()

		if err != nil {
			delete(startedContainers, v)
			log.Println("inspect error:", err)
			log.Println("removeDeadContainers dead container: ", v)
			continue
		}

		if string(out) != "true" {
			delete(startedContainers, v)
			log.Println("removeDeadContainers dead container: ", v)
		}
	}

	log.Println("removeDeadContainers done ")
}

