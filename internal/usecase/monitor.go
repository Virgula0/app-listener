package usecase

import (
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/repository"
)

// MonitorUseCase is the application service for the file system monitor
// feature. Consumers (the CLI delivery layer) only depend on this interface.
type MonitorUseCase interface {
	Start() error
	Stop()
	Events() <-chan ebpf.FileEvent
}

type monitorUseCase struct {
	repo repository.MonitorRepository
}

// NewMonitorUseCase wires the monitor use case to its repository
// implementation.
func NewMonitorUseCase(repo repository.MonitorRepository) MonitorUseCase {
	return &monitorUseCase{repo: repo}
}

func (u *monitorUseCase) Start() error {
	return u.repo.Start()
}

func (u *monitorUseCase) Stop() {
	u.repo.Stop()
}

func (u *monitorUseCase) Events() <-chan ebpf.FileEvent {
	return u.repo.Events()
}
