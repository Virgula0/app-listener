package printers

import log "github.com/sirupsen/logrus"

func PrintLogo() {
	log.Info("╔══════════════════════════════╗")
	log.Info("║        app-listener          ║")
	log.Info("║    File System Monitor       ║")
	log.Info("║         (eBPF)               ║")
	log.Info("╚══════════════════════════════╝")
}
