package printers

import (
	"github.com/Virgula0/app-listener/internal/constants"
	log "github.com/sirupsen/logrus"
)

func PrintLogo() {
	log.Info("╔══════════════════════════════════╗")
	log.Info("║           " + constants.AppName + "           ║")
	log.Info("║             Version " + constants.Version + "           ║")
	log.Info("║   File System & Network Monitor  ║")
	log.Info("║              (eBPF)              ║")
	log.Info("╚══════════════════════════════════╝")
}
