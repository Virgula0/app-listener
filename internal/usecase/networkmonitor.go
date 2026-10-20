package usecase

import (
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/repository"
)

// NetworkMonitorUseCase is the application service for the network monitor feature.
type NetworkMonitorUseCase interface {
	StreamUseCase[ebpf.NetEvent]
	SetEventTypes(types []ebpf.NetEventType)
}

type networkMonitorUseCase struct {
	streamUseCase[ebpf.NetEvent]
	repo repository.NetworkMonitorRepository
}

// NewNetworkMonitorUseCase wires the network monitor use case to its repository implementation.
func NewNetworkMonitorUseCase(repo repository.NetworkMonitorRepository) NetworkMonitorUseCase {
	return &networkMonitorUseCase{streamUseCase: streamUseCase[ebpf.NetEvent]{repo: repo}, repo: repo}
}

func (u *networkMonitorUseCase) SetEventTypes(types []ebpf.NetEventType) {
	u.repo.SetEventTypes(types)
}
