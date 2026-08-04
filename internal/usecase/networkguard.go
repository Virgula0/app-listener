package usecase

import (
	"github.com/Virgula0/app-listener/internal/networkguard"
	"github.com/Virgula0/app-listener/internal/repository"
)

// NetworkGuardUseCase is the application service for the network guard
// feature.
type NetworkGuardUseCase interface {
	Start() error
	Stop()
	Events() <-chan networkguard.NetGuardEvent
}

type networkGuardUseCase struct {
	repo repository.NetworkGuardRepository
}

// NewNetworkGuardUseCase wires the network guard use case to its repository
// implementation.
func NewNetworkGuardUseCase(repo repository.NetworkGuardRepository) NetworkGuardUseCase {
	return &networkGuardUseCase{repo: repo}
}

func (u *networkGuardUseCase) Start() error {
	return u.repo.Start()
}

func (u *networkGuardUseCase) Stop() {
	u.repo.Stop()
}

func (u *networkGuardUseCase) Events() <-chan networkguard.NetGuardEvent {
	return u.repo.Events()
}
