package usecase

import (
	"github.com/Virgula0/app-listener/internal/networkguard"
	"github.com/Virgula0/app-listener/internal/repository"
)

// NetworkGuardUseCase is the application service for the network guard feature.
type NetworkGuardUseCase = StreamUseCase[networkguard.NetGuardEvent]

// NewNetworkGuardUseCase wires the network guard use case to its repository implementation.
func NewNetworkGuardUseCase(repo repository.NetworkGuardRepository) NetworkGuardUseCase {
	return &streamUseCase[networkguard.NetGuardEvent]{repo: repo}
}
