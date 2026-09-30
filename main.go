package main

import (
	"errors"
	"os"
	"runtime"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/cmd"
	"github.com/Virgula0/app-listener/internal/constants"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := cmd.Execute(); err != nil {
		log.Error(err.Error())
		// A critical startup failure (constants.ErrCriticalStartup) reproduces on every restart:
		// exit with a distinct status (not log.Fatal's 1) so the unit's RestartPreventExitStatus
		// stops systemd retrying. Every other error still exits 1, so Restart=on-failure keeps
		// retrying transient failures.
		if errors.Is(err, constants.ErrCriticalStartup) {
			os.Exit(constants.CriticalExitCode)
		}
		os.Exit(1)
	}
}
