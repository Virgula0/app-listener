package usecase

import (
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/repository"
)

// MonitorUseCase is the application service for the file system monitor feature.
type MonitorUseCase = StreamUseCase[ebpf.FileEvent]

// NewMonitorUseCase wires the monitor use case to its repository implementation.
func NewMonitorUseCase(repo repository.MonitorRepository) MonitorUseCase {
	return &streamUseCase[ebpf.FileEvent]{repo: repo}
}
