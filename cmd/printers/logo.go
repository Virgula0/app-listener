package printers

import log "github.com/sirupsen/logrus"

func PrintLogo() {
	log.Info("╔══════════════════════════════════╗")
	log.Info("║           App-Listener           ║")
	log.Info("║   File System & Network Monitor  ║")
	log.Info("║              (eBPF)              ║")
	log.Info("╚══════════════════════════════════╝")
}
