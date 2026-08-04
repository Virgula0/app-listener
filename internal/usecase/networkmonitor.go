package usecase

import (
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/repository"
)

// NetworkMonitorUseCase is the application service for the network monitor
// feature.
type NetworkMonitorUseCase interface {
	Start() error
	Stop()
	Events() <-chan ebpf.NetEvent
	SetEventTypes(types []ebpf.NetEventType)
}

type networkMonitorUseCase struct {
	repo repository.NetworkMonitorRepository
}

// NewNetworkMonitorUseCase wires the network monitor use case to its
// repository implementation.
func NewNetworkMonitorUseCase(repo repository.NetworkMonitorRepository) NetworkMonitorUseCase {
	return &networkMonitorUseCase{repo: repo}
}

func (u *networkMonitorUseCase) Start() error {
	return u.repo.Start()
}

func (u *networkMonitorUseCase) Stop() {
	u.repo.Stop()
}

func (u *networkMonitorUseCase) Events() <-chan ebpf.NetEvent {
	return u.repo.Events()
}

func (u *networkMonitorUseCase) SetEventTypes(types []ebpf.NetEventType) {
	u.repo.SetEventTypes(types)
}
