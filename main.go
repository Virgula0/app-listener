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
		// A critical startup failure (see constants.ErrCriticalStartup) will
		// reproduce identically on every restart: exit with a distinct
		// status instead of log.Fatal's generic 1, so the systemd unit's
		// RestartPreventExitStatus can tell systemd to stop retrying instead
		// of crash-looping forever. Every other error keeps exiting 1, so
		// Restart=on-failure still retries transient failures as before.
		if errors.Is(err, constants.ErrCriticalStartup) {
			os.Exit(constants.CriticalExitCode)
		}
		os.Exit(1)
	}
}
