package main

import (
	"runtime"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/cmd"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := cmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}
